package dsl


import (

        "net/http"
        "net/url"
        "log"
        "io/ioutil"
        "fmt"

	"gopkg.in/telegram-bot-api.v4"

        lua "github.com/yuin/gopher-lua"


)


// Do not use it
func TelegramStartBot (token string) {

	bot, err := tgbotapi.NewBotAPI(token)
        if err != nil {
                log.Panic(err)
        }

        bot.Debug = true

        log.Printf("Authorized on account %s", bot.Self.UserName)

        u := tgbotapi.NewUpdate(0)
        u.Timeout = 45

}


func (d *dslConfig) TelegramSendMessage (L *lua.LState) int {

	var client http.Client
        chatId := L.CheckString(1)
        token := L.CheckString(2)
        message := L.CheckString(3)

        response, err := client.PostForm(
                fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token),
                url.Values{"chat_id": {chatId}, "text": {message}})

        defer response.Body.Close()

        if err != nil {
                log.Println(err)
                return 1
        }

        body, err := ioutil.ReadAll(response.Body)
        log.Printf("[DEBUG] telegram responce: ", string(body))

        if err != nil {
                log.Println(err)
                return 1
        }

        return 0

}