local cron = require("cron")
local time = require("time")


scheduler = cron.new()
scheduler:new('@every 1s', './test/hello.lua')
time.sleep(1)
