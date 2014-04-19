package paxostest

import (
	"flag"
	"github.com/cmu440/paxos"
	"github.com/cmu440/rpc/paxosrpc"
	"log"
	"os"
	"regexp"
	"time"
)

type testFunc struct {
	name string
	f    func()
}

var (
	passCount int
	failCount int
	master    paxos.Paxos
	testRegex = flag.String("t", "", "test to run")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	f, _ := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)
}

var LOGE = log.New(os.Stderr, "", log.Lshortfile|log.Lmicroseconds)

func testPaxosBasic1() {
	log.Println("FUCK DA POLICE")
	master.Propose(&paxosrpc.ProposeArgs{struct{}{}}, new(paxosrpc.ProposeReply))
	time.Sleep(5 * time.Second)
}

func StartTests(master paxos.Paxos) {
	tests := []testFunc{
		{"testPaxosBasic1", testPaxosBasic1},
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
