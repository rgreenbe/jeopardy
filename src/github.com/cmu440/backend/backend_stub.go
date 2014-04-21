package backend

import (
	"log"
	//"net"
	"strings"
)

type stub struct {
}

func NewStub() Backend {
	return &stub{}
}

func (s *stub) RecvCommit(commitMessage []byte) error {
	commit := string(commitMessage)
	items := strings.Split(commit, ",")
	//hostPort := items[0]
	message := items[1]
	log.Println("Backend got", message)
	//conn, _ := net.Dial("tcp", hostPort)
	//conn.Write([]byte(message))
	//conn.Close()
	return nil
}
