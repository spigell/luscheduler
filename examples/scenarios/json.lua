SCHEDULE = '@every 20s'

package.path = filepath.dir(debug.getinfo(1).source)..'/../?.lua;'.. package.path
settings = require "settings"

request, err = http.request{
  type = "POST", 
  url = "http://127.0.0.1:8081/post",
  body = "{ testo: testo }",
}

resp = json.decode(request.body)

print(resp.status)

request, err = http.request{
  type = "POST", 
  url = "http://127.0.0.1:8081/postlong",
  body = "{ testo: testo }",
}

resp = json.decode(request.body)

-- json:  { "items": [{"status": "success post", "description": "you are posted smth!"}, {"status": "cool", "description": "test" }]}
for _, item in pairs(resp) do

  for _, i in pairs(item) do

    print("id: "..i.id)
    print("status: "..i.status)
    print("description: "..i.description)

  end

end






