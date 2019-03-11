SCHEDULE = '@every 5s'


local http = require("http")
local client = http.client({insecure_ssl=true})


local request = http.request("GET", "http://127.0.0.1:8081/get")
local result, err = client:do_request(request)

if err then error(err) end

if not(result.code == 200) then error("code") end

print(result.body)
