package dsl

import (

	lua "github.com/yuin/gopher-lua"
	"log"

)


func (d *dslConfig) dslNewSchedule(L *lua.LState) int {
	schedule := L.CheckString(1)
	scenario := L.CheckString(2)
        log.Printf("[INFO] add new schedule: `%s` with scenario `%s`\n", schedule, scenario)
        
        cron := d.Cron

        cron.AddFunc(schedule, func (){ (run(scenario)) })

        return 0
}


func run (s string) {
	state := lua.NewState()
        script, config := Prepare(s)
        Register(config, state)
	if err := state.DoFile(script); err != nil {
		log.Printf("[ERROR]: ", err)
	}
}

