package dsl

import (
	"bytes"
	"io/ioutil"
	"log"

	"golang.org/x/crypto/ssh"

	lua "github.com/yuin/gopher-lua"
)

type dslSshConfig struct {
	user   string
	host   string
	port   string
	key    string
	config *ssh.ClientConfig
}

func (s *dslSshConfig) buildConfig() error {

	key, err := ioutil.ReadFile(s.key)
	if err != nil {
		return err
	}

	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return err
	}

	sshConfig := &ssh.ClientConfig{
		User: s.user,
		Auth: []ssh.AuthMethod{ssh.PublicKeys(signer)},
	}

	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	s.config = sshConfig

	return nil
}

func (s *dslSshConfig) makeSession() (*ssh.Session, error) {

	err := s.buildConfig()
	if err != nil {
		return nil, err
	}

	host := s.host + ":" + s.port

	client, err := ssh.Dial("tcp", host, s.config)
	if err != nil {
		return nil, err
	}

	session, err := client.NewSession()
	if err != nil {
		return nil, err
	}
	return session, nil
}

func (d *dslState) dslSshExecute(L *lua.LState) int {
	session := checkSshConn(L)
	args := L.CheckTable(2)
	command := args.RawGetString("command").String()

	var output bytes.Buffer
	session.Stdout = &output

	if err := session.Run(command); err != nil {
		log.Printf("[ERROR] Ssh execution failed: ", err)
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}
	session.Close()

	result := L.NewTable()
	L.SetField(result, "output", lua.LString(output.String()))
	L.Push(result)

	return 1

}
func (d *dslState) dslSshAuth(L *lua.LState) int {
	args := L.CheckTable(1)
	host := args.RawGetString("host").String()
	user := args.RawGetString("user").String()
	port := args.RawGetString("port").String()
	if port == "nil" {
		port = "22"
	}

	key := args.RawGetString("key").String()

	conn := dslSshConfig{host: host,
		user: user,
		port: port,
		key:  key,
	}

	session, err := conn.makeSession()

	if err != nil {
		log.Printf("[ERROR] Ssh auth/session failed: ", err)
		L.Push(lua.LNil)
		L.Push(lua.LString(err.Error()))
		return 2
	}

	ud := L.NewUserData()
	ud.Value = session
	L.SetMetatable(ud, L.GetTypeMetatable("ssh"))
	L.Push(ud)
	log.Printf("[INFO] New ssh connection to `%s:%s`\n", host, port)
	return 1
}

func checkSshConn(L *lua.LState) *ssh.Session {
	ud := L.CheckUserData(1)
	if v, ok := ud.Value.(*ssh.Session); ok {
		return v
	}
	L.ArgError(1, "This is not a ssh connection")
	return nil
}
