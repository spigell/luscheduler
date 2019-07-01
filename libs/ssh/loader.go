package ssh

import (
	lua "github.com/yuin/gopher-lua"
)

func Preload(L *lua.LState) {
	L.PreloadModule(`ssh`, Loader)
}

// Loader is the module loader function.
func Loader(L *lua.LState) int {
	ssh := L.NewTypeMetatable(`ssh`)
	L.SetGlobal(`ssh`, ssh)
	L.SetField(ssh, `__index`, L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		`execute`: Execute,
		//"copy":    config.dslScpCopy,
	}))
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)
	return 1

}

var api = map[string]lua.LGFunction{
	`auth`: Auth,
}
