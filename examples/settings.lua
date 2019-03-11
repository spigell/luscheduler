
function sendTelegram(message)
	response, err = telegram.sendmessage("CHATID", "TOKEN", message)
	print(response)
end
