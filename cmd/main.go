package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"

	utils "luscheduler/utils"
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
		utils.RunFileOnce(scenario)
		os.Exit(0)
	}

	file, _ := os.Open(*config)
	configuration := Config{}
	target, _ := ioutil.ReadAll(file)

	err := yaml.Unmarshal(target, &configuration)
	if err != nil {
		log.Println("[ERROR] Error while parsing configuration: ", err)
	}

	utils.RunFileOnce(configuration.InitScript)
}
