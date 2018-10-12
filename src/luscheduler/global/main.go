package global

import (
        "strconv"
        "time"
        "io"
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


func MergeSettings( storage, filename, settings string ) string {
        t := strconv.FormatInt(time.Now().Unix(), 10)
        target := storage + "/" + t + ".lua"

        final, err := os.OpenFile(target, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
        if err != nil {
                log.Fatal("[ERROR] failed to open file for writing: ", err)
        }

        first, err := os.Open(settings)
        if err != nil {
                log.Println("[ERROR] failed to open settings file:", err)
        }
        defer first.Close()
 
        second, err := os.Open(filename)
        if err != nil {
                log.Fatal("[ERROR] failed to open scenario file: ", err)
        }
        defer second.Close()

        if _, err := io.Copy(final, first); err != nil {
                log.Panic(err)
        }

        if _, err := io.Copy(final, second); err != nil {
                log.Panic(err)
        }

        if err := final.Close(); err != nil {
                log.Panic(err)
        }

        return target

} 