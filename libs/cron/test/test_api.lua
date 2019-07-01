local cron = require("cron")
local time = require("time")


scheduler = cron.new()
scheduler:add('@every 1s', './test/hello.lua')
--scheduler:new('@every 2s', '../../libs/ssh/test/test_api.lua')
time.sleep(100)
