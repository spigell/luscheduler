SCHEDULE = '@every 10s'

p, err = zabbix.login("http://127.0.0.1:8081/api_jsonrpc.php", "admin", "zabbix")
alarms = p:alarms{pattern = ".*", severity = "5", duration = "100"}

for _, alarm in pairs(alarms) do
	sendTelegram("Описание: "..alarm["description"].."\nХост: "..alarm["host"].."\nПриоритет: "..alarm["priority"].."\nТекущее значение: "..alarm["lastvalue"])
end

p:logout()

