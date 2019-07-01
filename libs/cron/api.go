package cron

import (
	"log"

	"github.com/robfig/cron"


	lua "github.com/yuin/gopher-lua"
)


func Run(s string) {
	state := lua.NewState()
	if err := state.DoFile(s); err != nil {
		log.Println("[ERROR] Error executing scenario: ", err)
	}
}

func NewSchedule(L *lua.LState) int {
	cron :=  checkCron(L)
	schedule := L.CheckString(2)
	scenario := L.CheckString(3)
	log.Printf("[INFO] add new schedule: `%s` with scenario `%s`\n", schedule, scenario)

	cron.AddFunc(schedule, func() { (Run(scenario)) })
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
