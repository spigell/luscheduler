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
	Cron *cron.Cron
}

type Telegram struct {
        Token   string
        ChatId  string
}



func Prepare( file string ) (string, *dslConfig) {
        CurrentCron.Start()

        script := global.MergeSettings(Conf.Storage, file, Conf.Settings)

        Tg := Telegram{ Token: Conf.Telegram.Token, ChatId: Conf.Telegram.ChatId }

        return script, &dslConfig{ Cron: CurrentCron, Telegram: Tg }
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


}
