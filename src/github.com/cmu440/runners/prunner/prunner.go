package main

import (
	"flag"
	"github.com/cmu440/backend"
	"github.com/cmu440/paxos"
	"log"
	"os"
	"strconv"
)

const defaultMasterPort = 9009

var (
	port           = flag.Int("port", defaultMasterPort, "port number to listen on")
	masterHostPort = flag.String("master", "", "master storage server host port (if non-empty then this storage server is a slave)")
	numNodes       = flag.Int("N", 1, "the number of nodes in the ring (including the master)")
	nodeID         = flag.Uint64("id", 0, "a 64-bit unsigned node ID to use for sequencing")
	masterID       = flag.Uint64("masterID", 0, "a 64-bit unsigned node ID to use for sequencing")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	f, _ := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)
}

func main() {
	flag.Parse()
	if *masterHostPort == "" && *port == 0 {
		// If masterHostPort string is empty, then this node is the master.
		*port = defaultMasterPort
	}
	// Create and start the paxos server
	_, err := paxos.NewPaxos(*masterHostPort, *numNodes, "localhost:"+strconv.Itoa(*port), *nodeID, *masterID, backend.NewStub())
	if err != nil {
		log.Fatalln("Failed to create Paxos node:", err)
	}
	log.Println("Started", *nodeID)
	// Run the node forever.
	select {}
}
