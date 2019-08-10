package cron

import (
	"github.com/robfig/cron"

	plugin "github.com/vadv/gopher-lua-libs/plugin"
	lua "github.com/yuin/gopher-lua"
)

func NewSchedule(L *lua.LState) int {
	cron := checkCron(L)
	schedule := L.CheckString(2)
	scenario := L.CheckUserData(3)

	state := lua.NewState()
	defer state.Close()
	state.Push(scenario)
	cron.AddFunc(schedule, func() { plugin.Run(state) })
	return 1
}

func New(L *lua.LState) int {
	currentCron := cron.New()
	currentCron.Start()

	ud := L.NewUserData()
	ud.Value = currentCron
	L.SetMetatable(ud, L.GetTypeMetatable(`cron`))
	L.Push(ud)
	return 1
}

func checkCron(L *lua.LState) *cron.Cron {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*cron.Cron); ok {
		return v
	}
	L.ArgError(1, "This is not a Cron")
	return nil
}
