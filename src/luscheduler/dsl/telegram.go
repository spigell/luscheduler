package dsl


import (

        "net/http"
        "net/url"
        "log"
        "io/ioutil"
        "fmt"

        lua "github.com/yuin/gopher-lua"


)


func (d *dslState) TelegramSendMessage (L *lua.LState) int {

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
        log.Printf("[DEBUG] telegram response: ", string(body))

        if err != nil {
                log.Println(err)
                return 1
        }

        return 0

}