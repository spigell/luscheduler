package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	lua "github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v2"

        native "luscheduler"
        cron "luscheduler/libs/cron"
)

var (
	config       = flag.String("config", "/etc/luscheduler.yml", "path to config file")
	version      = flag.Bool("version", false, "show version")
	exec         = flag.String("execute", "", "execute file(for testing purposes)")
	BuildVersion = "None"
)

type Config struct {
	InitScript string
}

func main() {

	flag.Parse()
	if *version {
		fmt.Printf("%s\n", BuildVersion)
		os.Exit(0)
	}

	if *exec != `` {
		scenario := *exec
		cron.Run(scenario)
		os.Exit(0)
	}

	file, _ := os.Open(*config)
	configuration := Config{}
	target, _ := ioutil.ReadAll(file)

	err := yaml.Unmarshal(target, &configuration)
	if err != nil {
		log.Println("[ERROR] Error while parsing configuration: ", err)
	}

	state := lua.NewState()
	defer state.Close()
	native.Preload(state)
	if err := state.DoFile(configuration.InitScript); err != nil {
		log.Printf("[FATAL] Main file: %s\n", err.Error())
	}
}
