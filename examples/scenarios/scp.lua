SCHEDULE = '@every 1500s'

session, err = ssh.auth{host = 'localhost', user = 'spigell', key = '/home/spigell/.ssh/keys/spigell.key'}

_, err = session:copy{source = '/etc/hosts', destination = '/tmp/hosts'}

if err then
  print(err)
end

