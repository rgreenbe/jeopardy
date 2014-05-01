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

/* The backend stub simply gets the messages that are sent through Paxos
* and echos them back to the provided host:port
*/
func (s *stub) RecvCommit(commitMessage []byte, master bool) error {
	commit := string(commitMessage)
	items := strings.Split(commit, ",")
	//log.Println(commit)
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
	conn.Close()
	return nil
}
