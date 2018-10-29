package dsl

import (
	"log"

	"github.com/robfig/cron"

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

	telegram := L.NewTypeMetatable("telegram")
	L.SetGlobal("telegram", telegram)
	L.SetField(telegram, "sendmessage", L.NewFunction(config.TelegramSendMessage))

	filepath := L.NewTypeMetatable("filepath")
	L.SetGlobal("filepath", filepath)
	L.SetField(filepath, "base", L.NewFunction(config.dslFilepathBasename))
	L.SetField(filepath, "dir", L.NewFunction(config.dslFilepathDir))
	L.SetField(filepath, "ext", L.NewFunction(config.dslFilepathExt))
	L.SetField(filepath, "glob", L.NewFunction(config.dslFilepathGlob))

	zabbix := L.NewTypeMetatable("zabbix")
	L.SetGlobal("zabbix", zabbix)
	L.SetField(zabbix, "login", L.NewFunction(config.dslZabbixLogin))
	L.SetField(zabbix, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"alarms": config.dslZabbixGetTriggers,
		"logout": config.dslZabbixLogout,
	}))

	http := L.NewTypeMetatable("http")
	L.SetGlobal("http", http)
	L.SetField(http, "request", L.NewFunction(config.dslHttpRequest))

	json := L.NewTypeMetatable("json")
	L.SetGlobal("json", json)
	L.SetField(json, "decode", L.NewFunction(config.dslJsonDecode))
	L.SetField(json, "encode", L.NewFunction(config.dslJsonEncode))

	ssh := L.NewTypeMetatable("ssh")
	L.SetGlobal("ssh", ssh)
	L.SetField(ssh, "auth", L.NewFunction(config.dslSshAuth))
	L.SetField(ssh, "__index", L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		"execute": config.dslSshExecute,
		"copy":    config.dslScpCopy,
	}))

}
