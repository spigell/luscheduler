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
        config := dsl.Prepare()
        dsl.Register(config, state)
        if err := state.DoFile(conf.InitScript); err != nil {
                log.Printf("[FATAL] Main file: %s\n", err.Error())
	}

}

