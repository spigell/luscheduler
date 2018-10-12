package main

import (

	"log"

	"luscheduler/global"
	"luscheduler/dsl"
	
	lua "github.com/yuin/gopher-lua"

)




func main() {

	conf := global.ReadConfiguration()

        state := lua.NewState()
        script, config := dsl.Prepare(conf.InitScript)
        dsl.Register(config, state)
        if err := state.DoFile(script); err != nil {
                log.Printf("[FATAL] Main file: %s\n", err.Error())
	}

}

