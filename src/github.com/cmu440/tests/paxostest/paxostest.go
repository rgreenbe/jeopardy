package main

import (
	"flag"
	"github.com/cmu440/rpc/paxosrpc"
	"log"
	"math/rand"
	"net"
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

type tester struct {
	messages   chan string
	hostPort   string
	connection *net.UDPConn
	count      int
}

var (
	passCount      int
	failCount      int
	master         *rpc.Client
	myHostPort     string
	conn           *net.Listener
	t              tester
	testRegex      = flag.String("t", "", "test to run")
	masterHostPort = flag.String("master", "", "The host:port of the master server")
	numNodes       = flag.Int("nodes", 3, "The number of nodes in the paxos ring")
)

func init() {
	log.SetFlags(log.Lshortfile | log.Lmicroseconds)
	f, _ := os.OpenFile("log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	log.SetOutput(f)
}

func (t *tester) sendAndListen(n int) {
	for i := 0; i < n; i++ {
		message := []byte(myHostPort + "," + strconv.Itoa(t.count))
		t.count++
		err := master.Call("Paxos.Propose", &paxosrpc.ProposeArgs{&message}, new(paxosrpc.ProposeReply))
		if err != nil {
			log.Println(err)
		}
	}
	for i := 0; i < n*(*numNodes); i++ {
		_ = <-t.messages
	}
	passCount++
	log.Println("Passed!")
}

func (t *tester) send(n int, node *rpc.Client) {
	for i := 0; i < n; i++ {
		message := []byte(myHostPort + "," + strconv.Itoa(i))
		err := node.Call("Paxos.Propose", &paxosrpc.ProposeArgs{&message}, new(paxosrpc.ProposeReply))
		if err != nil {
			log.Println(err)
		}
	}
}

func (t *tester) acceptConnections() {
	addr, err := net.ResolveUDPAddr("udp", t.hostPort)
	if err != nil {
		log.Println(err)
	}
	t.connection, err = net.ListenUDP("udp", addr)
	if err != nil {
		log.Println(err)
	}
	buf := make([]byte, 64)
	for {
		n, err := t.connection.Read(buf)
		if err != nil {
			log.Println(err)
			return
		}
		t.messages <- string(buf[0:n])
	}
}

func testPaxosBasic1() {
	t.sendAndListen(51)
}

func testPaxosBasic2() {
	t.sendAndListen(6)
}

func testPaxosBasic3() {
	t.sendAndListen(300)
}

func testPaxosDuelingLeaders() {
	reply := new(paxosrpc.GetServerReply)
	master.Call("Paxos.GetServers", new(paxosrpc.GetServerArgs), reply)
	nodes := (*reply).Nodes
	node1, _ := rpc.DialHTTP("tcp", nodes[1].HostPort)
	node2, _ := rpc.DialHTTP("tcp", nodes[2].HostPort)
	go t.send(50, node1)
	go t.send(50, node2)
	for i := 0; i < 300; i++ {
		<-t.messages
	}
	passCount++
	log.Println("Passed!")
}

func main() {
	flag.Parse()
	var err error
	master, err = rpc.DialHTTP("tcp", *masterHostPort)
	if err != nil {
		log.Fatalln("Failed to connect to the master server")
	}
	tests := []testFunc{
		//{"testPaxosBasic1", testPaxosBasic1},
		//{"testPaxosBasic2", testPaxosBasic2},
		//{"testPaxosBasic3", testPaxosBasic3},
		{"testPaxosDuelingLeaders", testPaxosDuelingLeaders},
	}
	// Run tests.
	rand.Seed(time.Now().Unix())
	myHostPort = "localhost:" + strconv.Itoa(10000+(rand.Int()%10000))
	t = tester{make(chan string, 1000), myHostPort, nil, 0}
	go t.acceptConnections()
	for _, t := range tests {
		if b, err := regexp.MatchString(*testRegex, t.name); b && err == nil {
			log.Printf("Running %s:\n", t.name)
			t.f()
			time.Sleep(time.Millisecond * 100)
		}
	}
	log.Printf("Passed (%d/%d) tests\n", passCount, passCount+failCount)
}
