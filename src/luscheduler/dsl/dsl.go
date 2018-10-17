package dsl

import (

        "github.com/robfig/cron"

        lua "github.com/yuin/gopher-lua"

        "luscheduler/global"

)

var(
        CurrentCron = cron.New()
        Conf = global.ReadConfiguration()
)     

type dslConfig struct{
        Telegram        Telegram
	Cron           *cron.Cron
}

type Telegram struct {
        Token   string
        ChatId  string
}



func Prepare() *dslConfig {
        CurrentCron.Start()

        Tg := Telegram{ Token: Conf.Telegram.Token, ChatId: Conf.Telegram.ChatId }

        return &dslConfig{ Cron: CurrentCron, Telegram: Tg }
}


func Register(config *dslConfig, L *lua.LState) {

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

}
