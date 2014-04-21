package backend

import (
	"log"
	"net"
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
	hostPort := items[0]
	message := items[1]
	log.Println("Backend got: ", message, "from: ", hostPort)
	conn, err := net.Dial("tcp", hostPort)
	if err != nil {
		log.Println(err)
	}
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Println(err)
	}
	conn.Close()
	return nil
}
