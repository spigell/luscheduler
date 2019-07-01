package cron

import (
	"log"

	lua "github.com/yuin/gopher-lua"
	time "github.com/vadv/gopher-lua-libs/time"
)

func Example_package() {
	state := lua.NewState()
	Preload(state)
	time.Preload(state)
	source := `
local time = require("time")
local cron = require("cron")
scheduler = cron.new()
scheduler:new('@every 1s', './test/hello.lua')
time.sleep(1)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// GO + lua = Cool
	// GO + lua = Cool
}
