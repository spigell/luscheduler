SCHEDULE = '@every 5s'

package.path = filepath.dir(debug.getinfo(1).source)..'/../?.lua;'.. package.path
settings = require "settings"

local params = {
  type = "GET", 
  url = "http://127.0.0.1:8081/get",
  useragent = "luscheduler", 
  headers = "token: 1234; Authorization: Bearer 132"
}

local result, err = http.request(params)

if not err then

  print(result.code)
  print(result.body)

else

  print(err)
end

--

params.type = "POST"
params.url = "http://127.0.0.1:8081/post"

local result, err = http.request(params)

if not err then

  print(result.code)
  print(result.body)

else

  print(err)
end
