package client

/*
import (
	"bufio"
	"fmt"
	"net"
)

func handleConnection(conn net.Conn, in chan string, out chan string) {
	go read(conn, in)
	go write(conn, out)
	for {
		select {
		case msg := <-in:
			fmt.Println("Message ", msg)
			out <- string(msg)

		}
	}

}

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		// handle error
	}
	read := make(chan string)
	write := make(chan string)
	for {
		conn, err := ln.Accept()
		if err != nil {
			// handle error
			continue
		}

		go handleConnection(conn, read, write)
	}
}
func read(conn net.Conn, in chan string) {
	reader := bufio.NewReader(conn)
	for {

		buf, _ := reader.ReadBytes('\n')
		fmt.Println("read")
		in <- string(buf)

	}

}
func write(conn net.Conn, out chan string) {
	for msg := range out {
		conn.Write([]byte(msg))
	}
}*/
