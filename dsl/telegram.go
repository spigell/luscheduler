package dsl

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"

	lua "github.com/yuin/gopher-lua"
)

func (d *dslState) TelegramSendMessage(L *lua.LState) int {

	var client http.Client
	chatId := L.CheckString(1)
	token := L.CheckString(2)
	message := L.CheckString(3)

	response, err := client.PostForm(
		fmt.Sprintf("https://api.telegram.org/bot%s/sendMessage", token),
		url.Values{"chat_id": {chatId}, "text": {message}})

	defer response.Body.Close()

	if err != nil {
		log.Printf("[ERROR] telegram send message failed: ", err)
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	body, err := ioutil.ReadAll(response.Body)

	if err != nil {
		log.Printf("[ERROR] telegram read telegram response: ", err)
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	//	log.Printf("[DEBUG] telegram response: ", string(body))
	L.Push(lua.LString(string(body)))

	return 1

}
