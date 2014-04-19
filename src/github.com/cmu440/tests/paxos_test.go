package main

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/cmu440/paxos"
	"github.com/cmu440/rpc/paxosrpc"
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
	pc        proxycounter.ProxyCounter
	paxos     Paxos
)

var LOGE = log.New(os.Stderr, "", log.Lshortfile|log.Lmicroseconds)

func initNodeServer(masterServerHostPort string, NodePort int) error {
	nodeServerHostPort := net.JoinHostPort("localhost", strconv.Itoa(NodePort))
	proxyCounter, err := proxycounter.NewProxyCounter(masterServerHostPort, nodeServerHostPort)
	if err != nil {
		LOGE.Println("Failed to setup test:", err)
		return err
	}
	pc = proxyCounter
	rpc.RegisterName("NodeServer", paxosrpc.Wrap(pc))

	// Create and start the TribServer.
	numNodes := 3
	paxos, err := paxos.NewPaxos("", numNodes, masterServerHostPort, 0)

	for i := 1; i < numNodes; i++ {
		nodeServerHostPort := net.JoinHostPort("localhost", strconv.Itoa(NodePort+i))
		_, _ := paxos.NewPaxos(masterServerHostPort, numNodes, nodeServerHostPort, i)
	}

	if err != nil {
		LOGE.Println("Failed to create Paxos Node:", err)
		return err
	}
	return nil
}

func main() {
	tests := []testFunc{
	//{"testCreateUserValid", testPaxosBasic1}
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
