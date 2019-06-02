local ssh = require("ssh")
session, err = ssh.auth{host = 'localhost', user = 'spigell', key = '/home/spigell/.ssh/keys/spigell.key'}


command = session:execute{command = "whoami && hostname -f"}

session2 = ssh.auth{host = 'localhost', user = 'spigell', key = '/home/spigell/.ssh/keys/spigell.key'}
command2, err2 = session2:execute{command = "broken command"}

print(command.output)
print(err2)
