local cron = require("cron")
local time = require("time")
local inspect = require("inspect")
local filepath = require('filepath')

scheduler = cron.new({verbose = 'true'})

job = scheduler:add_file('@every 3s', filepath.dir(debug.getinfo(1).source) .. '/hello.lua')
error_plugin = scheduler:add_file('@every 3s', filepath.dir(debug.getinfo(1).source) .. '/error.lua')

list = scheduler:list()
print(inspect(list))

time.sleep(4)
if job:last_error() then
  error("job must be without error. got - " .. tostring(job:last_error()))
end
