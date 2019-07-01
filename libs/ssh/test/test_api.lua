local ssh = require("ssh")
client, err = ssh.client{host = 'localhost', user = 'spigell', key = '/home/spigell/.ssh/keys/spigell.key'}

command = client:execute{command = "whoami && hostname -f && df -h"}
command2, err2 = client:execute{command = "broken_command"}

print(command.stdout)
print(command2.stderr)
