local ssh = require("ssh")
local inspect = require('inspect')

client, err = ssh.client{host = 'localhost', user = 'spigell', key = '/home/spigell/.ssh/keys/spigell.key'}

stdout = client:execute{command = "whoami && hostname -f && df -h && sleep 2"}
