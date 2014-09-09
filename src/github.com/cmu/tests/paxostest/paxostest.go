package main

import (
	"flag"
	"github.com/cmu440/backend"
	"github.com/cmu440/paxos"
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
	testType       = flag.String("type", "regular", "The type of tests to run...this could be dead node tests, regular tests, or catchup tests")
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
		err := master.Call("Paxos.Propose", &paxosrpc.ProposeArgs{message}, new(paxosrpc.ProposeReply))
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
		err := node.Call("Paxos.Propose", &paxosrpc.ProposeArgs{message}, new(paxosrpc.ProposeReply))
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

/*This is 51 messages because one of the tests has a node continually drop
*odd numbered command slots..thus in order for it to catchup to the last command, we
*have to submit an odd number of commands (and the other one drops the first 50 commits)
 */
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
	index := 1
	/* We want to make sure that we are that we are connecting to the right node
	* (the one we don't kill in certain tests) so this coordinates with paxostest.sh
	* in order to ensure communication with an alive process. This assumes that in a 3 node
	* system, one of the non-leader nodes gets killed
	 */
	if nodes[index].NodeID != 1 {
		index = 2
	}
	node, _ := rpc.DialHTTP("tcp", nodes[index].HostPort)
	go t.send(50, master)
	go t.send(50, node)
	for i := 0; i < 100*(*numNodes); i++ {
		<-t.messages
	}
	passCount++
	log.Println("Passed!")
}

func testPaxosReplaceNode() {
	t.send(10, master)
	for i := 0; i < (*numNodes-1)*10; i++ { //All but the dead node reply
		<-t.messages
	}
	reply := new(paxosrpc.GetServerReply)
	master.Call("Paxos.GetServers", new(paxosrpc.GetServerArgs), reply)
	nodes := (*reply).Nodes
	hostPort := "localhost:" + strconv.Itoa(10000+(rand.Int()%10000))
	stub := backend.NewStub()
	_, err := paxos.NewPaxos(nodes[0].HostPort, *numNodes, hostPort, uint64(*numNodes), 0, stub, paxosrpc.NONE)
	if err != nil {
		log.Println(err)
	}
	var replaceNode paxosrpc.Node
	for _, node := range nodes {
		if node.NodeID == 2 { //Again, we always assume 2 is killed..this is dependent on paxostest.sh
			replaceNode = node
		}
	}
	args := &paxosrpc.QuiesceArgs{nodes[0], replaceNode, paxosrpc.Node{hostPort, uint64(*numNodes)}}
	for _, node := range nodes {
		if node.NodeID != 2 {
			conn, err := rpc.DialHTTP("tcp", node.HostPort)
			if err != nil {
				log.Println(err)
			}
			conn.Call("Paxos.Quiesce", args, new(paxosrpc.QuiesceReply))
		}
	}
	for i := 0; i < 10; i++ {
		<-t.messages
	}
	master.Call("Paxos.GetServers", new(paxosrpc.GetServerArgs), reply)
	nodes = (*reply).Nodes
	for _, node := range nodes {
		conn, err := rpc.DialHTTP("tcp", node.HostPort)
		if err != nil {
			log.Println(err)
		}
		conn.Call("Paxos.Resume", new(paxosrpc.ResumeArgs), new(paxosrpc.ResumeReply))
	}
	t.sendAndListen(20)
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
		{"testPaxosBasic3", testPaxosBasic3},
		{"testPaxosDuelingLeaders", testPaxosDuelingLeaders},
	}
	if *testType == "dead" {
		*numNodes--
	} else if *testType == "replace" {
		tests = []testFunc{
			{"testPaxosReplaceNode", testPaxosReplaceNode},
		}
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
