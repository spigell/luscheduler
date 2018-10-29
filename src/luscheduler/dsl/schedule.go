package dsl

import (

	lua "github.com/yuin/gopher-lua"
	"log"

)


func (d *dslState) dslNewSchedule(L *lua.LState) int {
	schedule := L.CheckString(1)
	scenario := L.CheckString(2)
        log.Printf("[INFO] add new schedule: `%s` with scenario `%s`\n", schedule, scenario)
        
        cron := d.Cron

        cron.AddFunc(schedule, func (){ (Run(scenario)) })

        return 0
}



