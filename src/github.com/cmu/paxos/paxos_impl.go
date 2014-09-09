package paxos

import (
	"container/list"
	"errors"
	"github.com/cmu440/backend"
	"github.com/cmu440/rpc/paxosrpc"
	"log"
	"math"
	"net"
	"net/http"
	"net/rpc"
	"sync"
	"time"
)

const (
	LESS int = iota
	EQUAL
	GREATER
)

type Round struct {
	highestSequence *paxosrpc.Sequence
	previous        *paxosrpc.ValueSequence
	committed       bool
}

type paxos struct {
	quiesce         bool
	highestSequence *paxosrpc.Sequence
	contestedRound  uint64
	noopRound       uint64 //Indicator of what round we are proposing noops through, if catching up
	madeConnections bool
	numNodes        int
	nodeID          uint64
	masterID        uint64
	master          *rpc.Client
	proposalList    *list.List
	startPrepare    chan struct{}
	nodes           []paxosrpc.Node
	connections     []*rpc.Client
	learner         backend.Backend
	commits         map[uint64]*Round
	listLock        *sync.Mutex
	dataLock        *sync.Mutex
	quiesceLock     *sync.Mutex
	debug           paxosrpc.Debug
}

/*
* Start a new paxos server, register to receive rpc calls, connect to the master, and get the list of other nodes
 */
func NewPaxos(masterHostPort string, numNodes int, hostPort string, nodeID, masterID uint64, learner backend.Backend, debug paxosrpc.Debug) (Paxos, error) {
	var listener net.Listener
	var err error
	for {
		listener, err = net.Listen("tcp", hostPort)
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 200) //Retry in a bit
	}
	p := &paxos{false, nil, 0, 0, false, numNodes, nodeID, masterID, nil, list.New(), make(chan struct{}, 1000), nil,
		make([]*rpc.Client, 0, numNodes-1), learner, make(map[uint64]*Round), new(sync.Mutex), new(sync.Mutex), new(sync.Mutex), debug}
	for {
		err = rpc.RegisterName("Paxos", paxosrpc.Wrap(p))
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
	if masterHostPort != "" {
		var server *rpc.Client
		for {
			server, err = rpc.DialHTTP("tcp", masterHostPort)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 200)
		}
		reply := new(paxosrpc.GetServerReply)
		args := &paxosrpc.GetServerArgs{&paxosrpc.Node{hostPort, nodeID}}
		for {
			server.Call("Paxos.GetServers", args, reply)
			if reply.Status == paxosrpc.OK {
				p.nodes = reply.Nodes
				break
			}
			time.Sleep(time.Millisecond * 200)
		}
		p.master = server
	} else {
		p.nodes = make([]paxosrpc.Node, 0, numNodes)
		p.nodes = append(p.nodes, paxosrpc.Node{hostPort, nodeID})
	}
	rpc.HandleHTTP()
	go p.handlePrepare()
	go http.Serve(listener, nil)
	return p, nil
}

//See API
func (p *paxos) GetServers(args *paxosrpc.GetServerArgs, reply *paxosrpc.GetServerReply) error {
	seen := false
	if (*args).Node != nil {
		for _, node := range p.nodes {
			if node.NodeID == (*((*args).Node)).NodeID {
				seen = true
			}
		}
		if !seen {
			p.nodes = append(p.nodes, *((*args).Node))
		}
	}
	if len(p.nodes) >= p.numNodes {
		(*(reply)).Status = paxosrpc.OK
		(*(reply)).Nodes = p.nodes
		return nil
	}
	(*(reply)).Status = paxosrpc.NOT_READY
	return nil
}

//See API
func (p *paxos) ReplaceNode(args *paxosrpc.ReplaceNodeArgs, reply *paxosrpc.ReplaceNodeReply) error {
	oldNode := args.OldNode
	newNode := args.NewNode
	seen := false
	//Remove the old node if it exists and add in the new node if it doesn't
	for index, node := range p.nodes {
		if node.NodeID == oldNode.NodeID {
			p.nodes = append(p.nodes[:index], p.nodes[index+1:]...)
		} else if node.NodeID == newNode.NodeID {
			seen = true
		}
	}
	if !seen {
		p.nodes = append(p.nodes, newNode)
	}
	p.connectToNodes()

	(*reply).Done = true
	return errors.New("Old node does not exist")
}

//See API
func (p *paxos) MasterServer(args *paxosrpc.GetMasterArgs, reply *paxosrpc.GetMasterReply) error {
	var min uint64 = math.MaxUint32
	var minNode paxosrpc.Node
	for _, node := range p.nodes {
		if node.NodeID < min {
			min = node.NodeID
			minNode = node
		}
	}
	(*reply).Node = minNode
	return nil
}

//See API
func (p *paxos) Quiesce(args *paxosrpc.QuiesceArgs, reply *paxosrpc.QuiesceReply) error {
	p.quiesceLock.Lock()
	p.quiesce = true
	p.quiesceLock.Unlock()
L:
	for {
		select {
		case <-time.After(time.Millisecond * 100): //Ping to see commits are flushed every 100ms
			p.listLock.Lock()
			if p.proposalList.Len() == 0 {
				p.listLock.Unlock()
				break L
			}
			p.listLock.Unlock()
		}
	}
	//Now that we have flushed out all the old commits, add the new node
	oldNode := args.OldNode
	newNode := args.NewNode
	seen := false
	//Remove the old node if it exists and add in the new node if it doesn't
	for index, node := range p.nodes {
		if node.NodeID == oldNode.NodeID {
			p.nodes = append(p.nodes[:index], p.nodes[index+1:]...)
		} else if node.NodeID == newNode.NodeID {
			seen = true
		}
	}
	if !seen {
		p.nodes = append(p.nodes, newNode)
	}
	var nodeConn *rpc.Client
	p.connections = make([]*rpc.Client, 0, p.numNodes-1)
	for _, node := range p.nodes {
		if node.NodeID != p.nodeID {
			server, err := rpc.DialHTTP("tcp", node.HostPort)
			if node.NodeID == newNode.NodeID {
				nodeConn = server
			}
			if err == nil {
				p.connections = append(p.connections, server)
			} else {
				log.Println(err)
			}
		}
	}
	if p.nodeID == (*args).Update.NodeID { //Caller selects active node to push updates
		p.dataLock.Lock()
		var i uint64
		for i = 0; i < p.contestedRound; i++ {
			round, ok := p.commits[i]
			if !ok {
				log.Println("Log is inconsistent")
			} else {
				nodeConn.Call("Paxos.RecvCommit", &paxosrpc.CommitArgs{(*round).previous}, new(paxosrpc.CommitReply))
			}
		}
		p.dataLock.Unlock()
	}
	return nil
}

//See API
func (p *paxos) Resume(args *paxosrpc.ResumeArgs, reply *paxosrpc.ResumeReply) error {
	p.quiesceLock.Lock()
	defer p.quiesceLock.Unlock()
	p.quiesce = false
	return nil
}
