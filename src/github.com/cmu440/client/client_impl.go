package client

import (
	"bufio"
	"fmt"
	"net"
)

type jeopardyClient struct {
	master *rpc.Client
}

func (j *jeopardyClient) handleConnection(conn net.Conn) {
	go read(conn, in)
	go write(conn, out)
	for {
		select {
		case msg := <-j.readCH:
			fmt.Println("Message ", msg)
			j.writeCH <- string(msg)

		}
	}

}

func (j *jeopardyClient) handleClients(l *net.Listener) {
	for {
		conn, err := l.Accept()
		if err != nil {
			log.Println(err)
		}
		go j.handleReads(conn)
	}
}

func (j *jeopardyClient) handleReads(conn *net.Client) {
	data := make([]byte, 4096)
	for {
		n, err := conn.Read(data)
		args := &paxosrpc.ProposeArgs{data[:n]}
		master.Call("Paxos.Propose", args, new(paxosrpc.ProposeReply))
	}
}

func NewJeopardyClient(serverHost string, serverPort int, clientPort int) (jeopardyClient, error) {
	j := jeopardyClient{nil}
	master, err := rpc.DialHTTP("tcp", net.JoinHostPort(serverHost, strconv.Itoa(serverPort)))
	j.master = master
	if err != nil {
		log.Println(err)
	}
	client, err := net.Listen("tcp", ":"+strconv.Itoa(clientPort))
	go j.handleClients(client)
	if err != nil {
		log.Println(err)
	}
	return jeopardyClient{master}, nil
}
