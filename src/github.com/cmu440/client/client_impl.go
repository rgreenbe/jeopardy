package client

import (
	"github.com/cmu440/rpc/paxosrpc"
	"log"
	"net"
	"net/rpc"
)

type jeopardyClient struct {
	master *rpc.Client
}

func (j *jeopardyClient) handleClients(l *net.TCPListener) {
	for {
		conn, err := l.AcceptTCP()
		if err != nil {
			log.Println(err)
		}
		go j.handleReads(conn)
	}
}

func (j *jeopardyClient) handleReads(conn *net.TCPConn) {
	data := make([]byte, 4096)
	for {
		n, err := conn.Read(data)
		if err != nil {
			log.Println(err)
		} else {
			args := &paxosrpc.ProposeArgs{data[:n]}
			j.master.Call("Paxos.Propose", args, new(paxosrpc.ProposeReply))
		}
	}
}

func NewJeopardyClient(serverHostPort, clientHostPort string) (jeopardyClient, error) {
	j := jeopardyClient{nil}
	master, err := rpc.DialHTTP("tcp", serverHostPort)
	j.master = master
	if err != nil {
		log.Println(err)
	}
	addr, _ := net.ResolveTCPAddr("tcp", clientHostPort)
	client, err := net.ListenTCP("tcp", addr)
	go j.handleClients(client)
	if err != nil {
		log.Println(err)
	}
	return jeopardyClient{master}, nil
}
