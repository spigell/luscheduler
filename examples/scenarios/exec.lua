SCHEDULE = '@every 1s'

local cmd = require("cmd")

local command = "echo 1adsfasfasfdsaf"

local result, err = cmd.exec(command)

print(result.stdout)