package dsl

import (
	"log"

	"github.com/robfig/cron"

	libs "github.com/vadv/gopher-lua-libs"
	lua "github.com/yuin/gopher-lua"
)

var (
	CurrentCron = cron.New()
)

type dslState struct {
	Cron *cron.Cron
}

func Run(s string) {
	state := lua.NewState()
	config := Prepare()
	libs.Preload(state)
	Register(config, state)
	if err := state.DoFile(s); err != nil {
		log.Printf("[ERROR] Error executing scenario: ", err)
	}
}

func Prepare() *dslState {
	CurrentCron.Start()

	return &dslState{Cron: CurrentCron}
}

func Register(config *dslState, L *lua.LState) {

	schedule := L.NewTypeMetatable("schedule")
	L.SetGlobal("schedule", schedule)
	L.SetField(schedule, "new", L.NewFunction(config.dslNewSchedule))

	zabbix := L.NewTypeMetatable("zabbix")
	L.SetGlobal("zabbix", zabbix)
	L.SetField(zabbix, "login", L.NewFunction(config.dslZabbixLogin))
	L.SetField(zabbix, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"alarms": config.dslZabbixGetTriggers,
		"logout": config.dslZabbixLogout,
	}))

	ssh := L.NewTypeMetatable("ssh")
	L.SetGlobal("ssh", ssh)
	L.SetField(ssh, "auth", L.NewFunction(config.dslSshAuth))
	L.SetField(ssh, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"execute": config.dslSshExecute,
		"copy":    config.dslScpCopy,
	}))

}
