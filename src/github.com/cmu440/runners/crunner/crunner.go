package main

import (
	"flag"
	"github.com/cmu440/client"
	"log"
	"os"
)

var (
	clientHostPort  = flag.String("client", "localhost:8080", "Default port to listen on from client")
	jmasterHostPort = flag.String("master", "localhost:9009", "Default port for the master paxos server")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	f, _ := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)
}

func main() {
	flag.Parse()
	_, err := client.NewJeopardyClient(*jmasterHostPort, *clientHostPort)
	if err != nil {
		log.Fatalln("Failed to create Jeopardy! Client:", err)
	}
	log.Println("Started Jeopardy! client")
	select {}
}
