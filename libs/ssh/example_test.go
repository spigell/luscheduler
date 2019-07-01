package ssh

import (
	"log"

	lua "github.com/yuin/gopher-lua"
)

func Example_package() {
	state := lua.NewState()
	Preload(state)
	source := `
local ssh = require("ssh")
session, err = ssh.auth{host = 'localhost', user = 'spigell', key = '/home/spigell/.ssh/keys/spigell.key'}
command = session:execute{command = "echo true"}
print(command.output)
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// true
}
