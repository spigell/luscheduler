package cron

import (
	"log"

	time "github.com/vadv/gopher-lua-libs/time"
	lua "github.com/yuin/gopher-lua"
)

func Example_package() {
	state := lua.NewState()
	Preload(state)
	time.Preload(state)
	source := `
local time = require("time")
local cron = require("cron")
scheduler = cron.new()
scheduler:add('@every 1s', './test/hello.lua')
time.sleep(1)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// GO + lua = Cool
	// GO + lua = Cool
}
