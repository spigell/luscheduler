package cron

import (
	"testing"

	filepath "github.com/vadv/gopher-lua-libs/filepath"
	inspect "github.com/vadv/gopher-lua-libs/inspect"
	plugin "github.com/vadv/gopher-lua-libs/plugin"
	time "github.com/vadv/gopher-lua-libs/time"
	lua "github.com/yuin/gopher-lua"

	ssh "luscheduler/libs/ssh"
)

func TestApi(t *testing.T) {

	state := lua.NewState()
	Preload(state)
	time.Preload(state)
	ssh.Preload(state)
	inspect.Preload(state)

	plugin.Preload(state)
	if err := state.DoFile("./test/test_api.lua"); err != nil {
		t.Fatalf("execute test: %s\n", err.Error())
	}
}
