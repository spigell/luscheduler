package libs

import (
	ExternalLibs "github.com/vadv/gopher-lua-libs"
	ssh "luscheduler/libs/ssh"
	cron "luscheduler/libs/cron"

	lua "github.com/yuin/gopher-lua"
)

// Preload preload all gopher libs and luscheduler packages
func Preload(L *lua.LState) {
	ExternalLibs.Preload(L)
	cron.Preload(L)
	ssh.Preload(L)
}
