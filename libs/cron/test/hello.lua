local filepath = require("filepath")
print("GO + lua = Cool")

package.path = filepath.dir(debug.getinfo(1).source) .. '/?.lua;' .. package.path
print(package.path)
local external = require("test_api")

local ssh = external.ssh

client, err = ssh.client{host = 'localhost', user = 'spigell', key = '/home/spigell/.ssh/keys/spigell.key'}

command = client:execute{command = "whoami && hostname -f && df -h"}


print(command.stdout)