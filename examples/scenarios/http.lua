SCHEDULE = '@every 5s'

package.path = filepath.dir(debug.getinfo(1).source)..'/../?.lua;'.. package.path
settings = require "settings"

request, err = http.request{
  type = "GET", 
  url = "http://127.0.0.1:8081", 
  useragent = "luscheduler", 
  headers = "token: 1234; Authorization: Bearer 132"
}
print(request.code)
print(request.body)

request, err = http.request{
  type = "POST", 
  url = "http://127.0.0.1:8081", 
  useragent = "luscheduler", 
  user = "admin",
  password = "admin",
  body = "{ testo: testo }",
}

print(request.body)
print(request.code)





