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
		`add_file`: AddJobWithDoFile,
		`list`:     ListJobs,
	}))
	t := L.NewTable()
	L.SetFuncs(t, api)
	L.Push(t)

	job := L.NewTypeMetatable(`job`)
	L.SetGlobal(`job`, job)
	L.SetField(job, `__index`, L.SetFuncs(L.NewTable(), map[string]lua.LGFunction{
		`is_running`: IsJobRunning,
		`last_error`: Error,
	}))
	j := L.NewTable()
	L.SetFuncs(j, api)
	L.Push(j)
	return 1

}

var api = map[string]lua.LGFunction{
	`new`: New,
}
