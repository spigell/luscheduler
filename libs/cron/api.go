package cron

import (
	cron "github.com/robfig/cron/v3"

	externalLibs "github.com/vadv/gopher-lua-libs"
	lua "github.com/yuin/gopher-lua"
	"io/ioutil"
	"log"
	ssh "luscheduler/libs/ssh"
	"os"
	"sync"
)

type luaCronJob struct {
	sync.Mutex
	name      string
	id        int
	schedule  string
	state     *lua.LState
	running   bool
	error     error
	cronEntry cron.Entry
}

type luaScheduler struct {
	entries   []*luaCronJob
	scheduler *cron.Cron
	logger    *cron.Logger
}

func AddJobWithDoFile(L *lua.LState) int {
	c := checkCron(L, 1)
	schedule := L.CheckString(2)
	filepath := L.CheckString(3)

	s := &luaCronJob{name: filepath}
	s.Lock()
	state := newScheduleState()
	s.state = state
	s.schedule = schedule
	s.error = nil
	s.running = false

	job := cron.FuncJob(func() { go s.execute() })

	id, _ := c.scheduler.AddJob(schedule, job)
	s.id = int(id)
	s.cronEntry = c.scheduler.Entry(id)

	s.Unlock()

	c.entries = append(c.entries, s)
	ud := L.NewUserData()
	ud.Value = s
	L.SetMetatable(ud, L.GetTypeMetatable(`job`))
	L.Push(ud)

	return 1
}

func New(L *lua.LState) int {
	var params *lua.LTable

	if L.GetTop() == 0 {
		params = L.NewTable()
		params.RawSetString(`verbose`, lua.LString(`false`))
	} else {
		params = L.CheckTable(1)
	}

	verbose := params.RawGetString(`verbose`).String()
	var DefaultLogger cron.Logger
	if verbose == `true` {
		DefaultLogger = cron.VerbosePrintfLogger(log.New(os.Stdout, "cron: ", log.LstdFlags))
	} else {
		DefaultLogger = cron.PrintfLogger(log.New(ioutil.Discard, "", 0))
	}

	logger := cron.WithLogger(DefaultLogger)
	skipIfStillRunning := cron.WithChain(cron.SkipIfStillRunning(DefaultLogger))

	currentCron := cron.New(logger, skipIfStillRunning)
	currentCron.Start()

	scheduler := &luaScheduler{scheduler: currentCron}

	ud := L.NewUserData()
	ud.Value = scheduler
	L.SetMetatable(ud, L.GetTypeMetatable(`cron`))
	L.Push(ud)
	return 1
}

func IsJobRunning(L *lua.LState) int {
	s := checkluaCronJob(L, 1)
	L.Push(lua.LBool(s.getRunning()))
	return 1
}

func Error(L *lua.LState) int {
	s := checkluaCronJob(L, 1)
	err := s.getError()
	if err == nil {
		return 0
	}
	L.Push(lua.LString(err.Error()))
	return 1
}

func ListJobs(L *lua.LState) int {
	s := checkCron(L, 1)
	entries := s.entries
	list := L.NewTable()
	for i, v := range entries {

		var last_error string
		err := v.getError()

		if err == nil {
			last_error = `none`
		} else {
			last_error = err.Error()
		}

		v.Lock()
		t := L.NewTable()
		t.RawSetString(`name`, lua.LString(v.name))
		t.RawSetString(`schedule`, lua.LString(v.schedule))
		t.RawSetString(`id`, lua.LNumber(v.id))
		t.RawSetString(`next`, lua.LString(v.cronEntry.Next.String()))
		t.RawSetString(`last_error`, lua.LString(last_error))
		t.RawSetString(`running`, lua.LBool(v.running))
		v.Unlock()

		list.Insert(i, t)
	}

	L.Push(list)
	return 1
}

func checkCron(L *lua.LState, i int) *luaScheduler {
	ud := L.CheckUserData(i)
	if v, ok := ud.Value.(*luaScheduler); ok {
		return v
	}
	L.ArgError(1, "This is not a Cron instance")
	return nil
}

func checkluaCronJob(L *lua.LState, i int) *luaCronJob {
	ud := L.CheckUserData(i)
	if v, ok := ud.Value.(*luaCronJob); ok {
		return v
	}
	L.ArgError(1, "This is not a job")
	return nil
}

func newScheduleState() *lua.LState {
	state := lua.NewState()

	externalLibs.Preload(state)
	ssh.Preload(state)

	return state
}

func (s *luaCronJob) execute() {
	s.setRunning(true)
	s.setError(s.state.DoFile(s.name))
	s.setRunning(false)
}

func (s *luaCronJob) getRunning() bool {
	s.Lock()
	defer s.Unlock()
	running := s.running
	return running
}

func (s *luaCronJob) setRunning(value bool) {
	s.Lock()
	defer s.Unlock()
	s.running = value
}

func (s *luaCronJob) setError(err error) {
	s.Lock()
	defer s.Unlock()
	s.error = err
}

func (s *luaCronJob) getError() error {
	s.Lock()
	defer s.Unlock()
	err := s.error
	return err
}
