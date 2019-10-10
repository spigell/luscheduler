package libs

import (
	ExternalLibs "github.com/vadv/gopher-lua-libs"
	cron "luscheduler/libs/cron"
	ssh "luscheduler/libs/ssh"

	lua "github.com/yuin/gopher-lua"
)

// Preload preload all gopher libs and luscheduler packages
func Preload(L *lua.LState) {
	ExternalLibs.Preload(L)
	cron.Preload(L)
	ssh.Preload(L)
}
