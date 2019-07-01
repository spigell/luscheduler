package cron

import (
	lua "github.com/yuin/gopher-lua"
)

func Preload(L *lua.LState) {
	L.PreloadModule(`cron`, Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	cron := L.NewTypeMetatable(`cron`)
	L.SetGlobal(`cron`, cron)
	L.SetField(cron, `__index`, L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		`new`: NewSchedule,
	}))
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1

}

var api = map[string]lua.LGFunction{
	`new`: New,
}
