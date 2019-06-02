package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	lua "github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v2"

	oldDsl "luscheduler/dsl"
        native "luscheduler"
)

var (
	config       = flag.String("config", "/etc/luscheduler.yml", "path to config file")
	version      = flag.Bool("version", false, "show version")
	exec         = flag.String("execute", "", "execute file(for testing purposes)")
	BuildVersion = "None"
)

type Config struct {
	Storage    string
	Settings   string
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
		oldDsl.Run(scenario)
		os.Exit(0)
	}

	file, _ := os.Open(*config)
	configuration := Config{}
	target, _ := ioutil.ReadAll(file)

	err := yaml.Unmarshal(target, &configuration)
	if err != nil {
		log.Printf("[ERROR] Error while parsing configuration: ", err)
	}

	state := lua.NewState()
	defer state.Close()
	config := oldDsl.Prepare()
	native.Preload(state)
	oldDsl.Register(config, state)
	if err := state.DoFile(configuration.InitScript); err != nil {
		log.Printf("[FATAL] Main file: %s\n", err.Error())
	}

}
