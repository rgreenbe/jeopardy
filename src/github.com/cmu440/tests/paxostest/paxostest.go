package main

import (
	"flag"
	"github.com/cmu440/rpc/paxosrpc"
	"log"
	"net/rpc"
	"os"
	"regexp"
	"strconv"
	"time"
)

type testFunc struct {
	name string
	f    func()
}

var (
	passCount      int
	failCount      int
	master         *rpc.Client
	testRegex      = flag.String("t", "", "test to run")
	masterHostPort = flag.String("master", "", "The host:port of the master server")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	f, _ := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)
}

func testPaxosBasic1() {
	message := []byte("localhost:10004,0")
	err := master.Call("Paxos.Propose", &paxosrpc.ProposeArgs{message}, new(paxosrpc.ProposeReply))
	if err != nil {
		log.Println(err)
	}
	time.Sleep(1 * time.Second)
}

func testPaxosBasic2() {
	for i := 0; i < 5; i++ {
		message := []byte("localhost:10004," + strconv.Itoa(i))
		err := master.Call("Paxos.Propose", &paxosrpc.ProposeArgs{message}, new(paxosrpc.ProposeReply))
		if err != nil {
			log.Println(err)
		}
	}
	time.Sleep(1 * time.Second)
}

func main() {
	flag.Parse()
	var err error
	master, err = rpc.DialHTTP("tcp", *masterHostPort)
	if err != nil {
		log.Fatalln("Failed to connect to the master server")
	}

	tests := []testFunc{
		{"testPaxosBasic1", testPaxosBasic1},
		{"testPaxosBasic2", testPaxosBasic2},
	}

	// Run tests.
	for _, t := range tests {
		if b, err := regexp.MatchString(*testRegex, t.name); b && err == nil {
			log.Printf("Running %s:\n", t.name)
			t.f()
		}
	}

	log.Printf("Passed (%d/%d) tests\n", passCount, passCount+failCount)
}
