package cron

import (
	"log"

	inspect "github.com/vadv/gopher-lua-libs/inspect"
	time "github.com/vadv/gopher-lua-libs/time"
	lua "github.com/yuin/gopher-lua"
)

func Example_package() {
	state := lua.NewState()
	Preload(state)
	time.Preload(state)
	inspect.Preload(state)
	source := `
local time = require("time")
local cron = require("cron")
local inspect = require('inspect')
scheduler = cron.new()
scheduler:add_file('@every 1s', './test/hello.lua')
time.sleep(1.5)
list = scheduler:list()

print(list[0]['name'])
print(list[0]['id'])
print(list[0]['running'])
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// GO + lua = Cool
	//./test/hello.lua
	//1
	//true

}
