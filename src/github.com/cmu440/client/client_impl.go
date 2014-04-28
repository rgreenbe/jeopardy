package client

import (
	"bufio"
	"fmt"
	"net"
)

type jeopardyClient struct {
	client  *rpc.Client
	readCH  chan string
	writeCH chan string
}

func (*jeopardyClient) Start(port int) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}

		go j.handleConnection(conn)
	}
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

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}

		go handleConnection(conn)
	}
}
func (j *jeopardyClient) read(conn net.Conn) {
	reader := bufio.NewReader(conn)
	for {

		buf, _ := reader.ReadBytes('\n')
		fmt.Println("read")
		j.readCH <- string(buf)

	}

}
func (j *jeopardyClient) write(conn net.Conn) {
	for msg := range j.writeCH {
		conn.Write([]byte(msg))
	}
}

func NewJeopardyClient(serverHost string, serverPort int, clientPort int) (jeopardyClient, error) {
	/*cli, err := rpc.DialHTTP("tcp", net.JoinHostPort(serverHost, strconv.Itoa(serverPort)))
	if err != nil {
		return nil, err
	}*/
	ln, err := net.Listen("tcp", ":"+strconv.Itoa(clientPort))
	if err != nil {
		// handle error
	}
	return &jeopardyClient{client: nil, readCH: make(chan string), writeCH: make(chan string)}, nil
}

/*func (jc *jeopardyClient) SendMessage(userMessage string) {
	args := paxosrpc.Propose{[]byte(userMessage)}

}*/
