package global

import (
        "flag"
        "os"
        "log"
        "io/ioutil"


        "gopkg.in/yaml.v2"

)

var (
        config = flag.String("config", "/etc/luscheduler.yml", "path to config file")
)


type Config struct {
        Storage                 string 
        Settings                string
        InitScript              string
        Telegram                Telegram
}

type Telegram struct {

        ChatId          string
        Token        string
}


func ReadConfiguration() Config {
        flag.Parse()

        file, _ := os.Open(*config)
        configuration := Config{}
        target, _ := ioutil.ReadAll(file)

        err := yaml.Unmarshal(target, &configuration)
        if err != nil {
                log.Printf("[ERROR] Error while parsing configuration: ", err)
        }

        return configuration
}


