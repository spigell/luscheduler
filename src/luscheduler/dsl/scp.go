package dsl

import (
	"log"

	"github.com/tmc/scp"

	lua "github.com/yuin/gopher-lua"
)

func (d *dslState) dslScpCopy(L *lua.LState) int {
	session := checkSshConn(L)
	args := L.CheckTable(2)
	source := args.RawGetString("source").String()
	dest := args.RawGetString("destination").String()

	err := scp.CopyPath(source, dest, session)
	if err != nil {
		log.Printf("[ERROR] Scp failed: ", err)
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	session.Close()

	L.Push(lua.LString("Copied!"))

	return 1
}
