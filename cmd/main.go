package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	lua "github.com/yuin/gopher-lua"
	"gopkg.in/yaml.v2"

	libs "luscheduler"
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
		run(scenario)
		os.Exit(0)
	}

	file, _ := os.Open(*config)
	configuration := Config{}
	target, _ := ioutil.ReadAll(file)

	err := yaml.Unmarshal(target, &configuration)
	if err != nil {
		log.Println("[ERROR] Error while parsing configuration: ", err)
	}

	run(configuration.InitScript)

}

func run(file string) {
	state := lua.NewState()
	defer state.Close()
	libs.Preload(state)
	if err := state.DoFile(file); err != nil {
		log.Printf("[ERROR] file: %s\n", err.Error())
	}
}
