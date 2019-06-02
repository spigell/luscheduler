package libs

import (

	externalLibs "github.com/vadv/gopher-lua-libs"
	ssh "luscheduler/libs/ssh"

	lua "github.com/yuin/gopher-lua"
)

// Preload preload all gopher libs and luscheduler packages
func Preload(L *lua.LState) {
	externalLibs.Preload(L)
	ssh.Preload(L)
}