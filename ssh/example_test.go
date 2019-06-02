package ssh

import (
	"log"

	"github.com/vadv/gopher-lua-libs/time"
	lua "github.com/yuin/gopher-lua"
)

func Example_package() {
	state := lua.NewState()
	Preload(state)
	time.Preload(state)
	source := `
session, err = ssh.auth{host = 'localhost', user = 'spigell', key = '/home/spigell/.ssh/keys/spigell.key'}
command = session:execute{command = "whoami && hostname -f"}
`
	if err := state.DoString(source); err != nil {
		log.Fatal(err.Error())
	}
	// Output:
	// true
}
