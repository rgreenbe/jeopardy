package backend

import (
	"log"
	"net"
	"strings"
	"time"
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
	cAddr, err := net.ResolveUDPAddr("udp", ":0")
	sAddr, err := net.ResolveUDPAddr("udp", hostPort)
	if err != nil {
		log.Fatalln(err)
	}
	conn, err := net.DialUDP("udp", cAddr, sAddr)
	if err != nil {
		log.Fatalln(err)
	}
	_, err = conn.Write([]byte(message))
	if err != nil {
		log.Println(err)
	}
	time.Sleep(100 * time.Millisecond)
	conn.Close()
	return nil
}
