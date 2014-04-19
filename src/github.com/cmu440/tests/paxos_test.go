package main

import (
	"flag"
	"fmt"
	"github.com/cmu440/paxos"
	"github.com/cmu440/rpc/paxosrpc"
	"log"
	"net"
	"os"
	"regexp"
	"strconv"
)

type testFunc struct {
	name string
	f    func()
}

var (
	port      = flag.Int("port", 9010, "TribServer port number")
	testRegex = flag.String("t", "", "test to run")
	passCount int
	failCount int
	master    paxos.Paxos
)

var LOGE = log.New(os.Stderr, "", log.Lshortfile|log.Lmicroseconds)

func initNodeServer(masterServerHostPort string, NodePort int) error {
	// Create and start the TribServer.
	numNodes := 3
	var err error
	master, err = paxos.NewPaxos("", numNodes, masterServerHostPort, 0)
	for i := 1; i < numNodes; i++ {
		nodeServerHostPort := net.JoinHostPort("localhost", strconv.Itoa(NodePort+i))
		_, _ = paxos.NewPaxos(masterServerHostPort, numNodes, nodeServerHostPort, uint32(i))
	}

	if err != nil {
		LOGE.Println("Failed to create Paxos Node:", err)
		return err
	}
	return nil
}

func testPaxosBasic1() {
	fmt.Println("FUCK DA POLICE")
	master.Propose(&paxosrpc.ProposeArgs{struct{}{}}, new(paxosrpc.ProposeReply))
}

func main() {
	tests := []testFunc{
		{"testCreateUserValid", testPaxosBasic1},
	}
	flag.Parse()
	if flag.NArg() < 1 {
		LOGE.Fatal("Usage: tribtest <storage master host:port>")
	}

	if err := initNodeServer(flag.Arg(0), *port); err != nil {
		LOGE.Fatalln("Failed to setup TribServer:", err)
	}

	// Run tests.
	for _, t := range tests {
		if b, err := regexp.MatchString(*testRegex, t.name); b && err == nil {
			fmt.Printf("Running %s:\n", t.name)
			t.f()
		}
	}

	fmt.Printf("Passed (%d/%d) tests\n", passCount, passCount+failCount)
}
