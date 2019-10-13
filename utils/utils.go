package utils
import (
	lua "github.com/yuin/gopher-lua"

	"log"
	libs "luscheduler"
)

func RunFileOnce(file string) {
	state := lua.NewState()
	defer state.Close()
	libs.Preload(state)
	if err := state.DoFile(file); err != nil {
		log.Printf("[ERROR] file: %s\n", err.Error())
	}
}
