local cron = require("cron")
local time = require("time")
local plugin = require("plugin")
ssh = require("ssh")


scheduler = cron.new()


local plugin_body = [[
    local time = require("time")
    local chef = require("chef")
--    local cron = require("cron")

    local i = 1
    print(i)
    time.sleep(2)
]]

local string_plugin = plugin.do_string(plugin_body)
local file_plugin = plugin.do_file('./test/hello.lua')

scheduler:add('@every 3s', string_plugin)
scheduler:add('@every 10s', file_plugin)

time.sleep(11)
if file_plugin:is_running() then
  print("HI")
else
  print("NO")
  print(file_plugin:error())
end


time.sleep(100)
