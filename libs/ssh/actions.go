package ssh

import (
	"bytes"
	
	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
	lua "github.com/yuin/gopher-lua"
)

func Execute(L *lua.LState) int {
	client := checkClient(L)
	args := L.CheckTable(2)
	command := args.RawGetString("command").String()

	session, err := makeSession(client)
	if err != nil {
		sshError(L, err)
	}
	defer session.Close()

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	session.Stdout = &stdout
	session.Stderr = &stderr

	if err := session.Run(command); err != nil {
		sshError(L, err)
	}

	result := L.NewTable()
	L.SetField(result, `stdout`, lua.LString(stdout.String()))
	L.SetField(result, `stderr`, lua.LString(stderr.String()))
	L.Push(result)

	return 1
}


func makeSession(client *ssh.Client) (*ssh.Session, error) {
	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func Copy(L *lua.LState) int {
	client := checkClient(L)
	args := L.CheckTable(2)
	source := args.RawGetString("source").String()
	dest := args.RawGetString("destination").String()

	session, err := makeSession(client)
	defer session.Close()

	if err != nil {
		sshError(L, err)
	}

	if err := scp.CopyPath(source, dest, session); err != nil {
		sshError(L, err)
	}

	L.Push(lua.LString("Copied!"))

	return 1
}
