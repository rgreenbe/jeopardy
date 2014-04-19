package paxos

import (
	"container/list"
	"errors"
	"fmt"
	"github.com/cmu440/rpc/paxosrpc"
	"math"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

const (
	LESS int = iota
	EQUAL
	GREATER
)

type paxos struct {
	madeConnections bool
	numNodes        int
	nodeID          uint64
	masterID        uint64
	master          *rpc.Client
	proposalList    *list.List
	startPrepare    chan struct{}
	nodes           []paxosrpc.Node
	connections     []*rpc.Client
	highestSequence *paxosrpc.Sequence
	previous        *paxosrpc.ValueSequence
}

func NewPaxos(masterHostPort string, numNodes int, hostPort string, nodeID, masterID uint64) (Paxos, error) {
	var listener net.Listener
	var err error
	for {
		listener, err = net.Listen("tcp", hostPort)
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 200) //Retry in a second
	}
	p := &paxos{false, numNodes, nodeID, masterID, nil, list.New(), make(chan struct{}, 1000), nil,
		make([]*rpc.Client, 0, numNodes-1), nil, nil}
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
			time.Sleep(time.Second)
		}
		p.master = server
		go p.handlePrepare()
		go http.Serve(listener, nil)
	} else {
		p.nodes = make([]paxosrpc.Node, 0, numNodes)
		p.nodes = append(p.nodes, paxosrpc.Node{hostPort, nodeID})
		rpc.HandleHTTP()
		go p.handlePrepare()
		go http.Serve(listener, nil)
		return p, nil
	}
	return p, nil
}

func (p *paxos) RecvPrepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error {
	if p.compare(p.highestSequence, (*args).Sequence) == LESS {
		(*reply).Status = paxosrpc.OK
		(*reply).Previous = p.previous
		p.highestSequence = args.Sequence
	} else {
		(*reply).Status = paxosrpc.CANCEL
	}
	return nil
}

func (p *paxos) RecvAccept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error {
	if p.compare(p.highestSequence, (*args).Accept.Sequence) == LESS {
		(*reply).Status = paxosrpc.OK
		p.previous = (*args).Accept
	} else {
		(*reply).Status = paxosrpc.CANCEL
	}
	return nil
}

func (p *paxos) RecvCommit(args *paxosrpc.CommitArgs, reply *paxosrpc.CommitReply) error {
	fmt.Println("Committed")
	return nil
}

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
	if len(p.nodes) == p.numNodes {
		(*(reply)).Status = paxosrpc.OK
		(*(reply)).Nodes = p.nodes
		return nil
	}
	(*(reply)).Status = paxosrpc.NOT_READY
	return nil
}

func (p *paxos) ReplaceNode(oldNode *paxosrpc.Node, newNode *paxosrpc.Node) error {
	for index, node := range p.nodes {
		if node.NodeID == (*oldNode).NodeID {
			p.nodes[index] = *newNode
			server, err := rpc.DialHTTP("tcp", (*newNode).HostPort)
			if err != nil {
				p.connections[index] = server
			}
			break
		}
	}
	return errors.New("Old node does not exist")
}

func (p *paxos) MasterServer(args *paxosrpc.GetMasterArgs, reply *paxosrpc.GetMasterReply) error {
	if p.previous != nil {
		id := p.previous.Sequence.NodeID
		for _, node := range p.nodes {
			if node.NodeID == id {
				(*reply).Node = node
				return nil
			}
		}
	} else {
		var min uint64 = math.MaxUint32
		var minNode paxosrpc.Node
		for _, node := range p.nodes {
			if node.NodeID < min {
				min = node.NodeID
				minNode = node
			}
		}
		(*reply).Node = minNode
	}
	return nil
}

func (p *paxos) Propose(args *paxosrpc.ProposeArgs, reply *paxosrpc.ProposeReply) error {
	(*(reply)).Status = paxosrpc.OK
	p.proposalList.PushBack((*args).Proposal)
	p.startPrepare <- struct{}{}
	return nil
}

func (p *paxos) handlePrepare() {
	for {
		select {
		case <-p.startPrepare:
			p.sendPrepare()
		}
	}
}

func (p *paxos) sendPrepare() {
	var n uint64
	if p.highestSequence != nil {
		n = p.highestSequence.N + 1
	} else {
		n = 1
	}
	if !p.madeConnections {
		p.connectToNodes()
		p.madeConnections = true
	}
	args := &paxosrpc.PrepareArgs{&paxosrpc.Sequence{n, p.nodeID}}
	replyChan := make(chan paxosrpc.PrepareReply, p.numNodes-1)
	for _, connection := range p.connections {
		go p.rpcPrepare(connection, args, replyChan)
	}
	oldestPrepare := new(paxosrpc.ValueSequence)
	ok := 1
	cancel := 0
	for i := 1; i < p.numNodes; i++ {
		reply := <-replyChan
		if reply.Status == paxosrpc.OK {
			if reply.Previous != nil {
				if oldestPrepare == nil {
					oldestPrepare = reply.Previous
				} else if p.compare(oldestPrepare.Sequence, (*reply.Previous).Sequence) == LESS {
					oldestPrepare = reply.Previous
				}
			}
			ok++
			if (p.numNodes / 2) < ok {

				p.sendAccept(oldestPrepare)
				return
			}
		} else {
			cancel++
			if (p.numNodes / 2) < cancel {
				p.startPrepare <- struct{}{} //So sorry..try again
				return
			}
		}
	}
}

func (p *paxos) sendAccept(accept *paxosrpc.ValueSequence) {
	args := &paxosrpc.AcceptArgs{accept}
	replyChan := make(chan paxosrpc.Status, p.numNodes-1)
	for _, connection := range p.connections {
		go p.rpcAccept(connection, args, replyChan)
	}
	ok := 1
	cancel := 0
	for i := 1; i < p.numNodes; i++ {
		reply := <-replyChan
		if reply == paxosrpc.OK {
			ok++
			if (p.numNodes / 2) < ok {
				p.sendCommit(accept)
			}
		} else {
			cancel++
			if (p.numNodes / 2) < cancel {
				p.startPrepare <- struct{}{} //So sorry..try again
				return
			}
		}
	}
}

func (p *paxos) sendCommit(commit *paxosrpc.ValueSequence) {
	args := &paxosrpc.CommitArgs{commit}
	for _, connection := range p.connections {
		go p.rpcCommit(connection, args)
	}
	p.startPrepare <- struct{}{}
}

func (p *paxos) rpcCommit(server *rpc.Client, args *paxosrpc.CommitArgs) {
	reply := new(paxosrpc.CommitReply)
	server.Call("Paxos.RecvCommit", args, reply)
}

func (p *paxos) rpcAccept(server *rpc.Client, args *paxosrpc.AcceptArgs, replyChan chan paxosrpc.Status) {
	reply := new(paxosrpc.AcceptReply)
	server.Call("Paxos.RecvAccept", args, reply)
	replyChan <- (*reply).Status
}

func (p *paxos) rpcPrepare(server *rpc.Client, args *paxosrpc.PrepareArgs, replyChan chan paxosrpc.PrepareReply) {
	reply := new(paxosrpc.PrepareReply)
	server.Call("Paxos.RecvPrepare", args, reply)
	replyChan <- *reply
}

func (p *paxos) compare(highest, prepare *paxosrpc.Sequence) int {
	if highest.N < prepare.N {
		return LESS
	} else if highest.N == prepare.N {
		if prepare.NodeID == highest.NodeID {
			return EQUAL
		} else if prepare.NodeID < highest.NodeID {
			return LESS
		}
		return GREATER
	}
	return GREATER
}

func (p *paxos) connectToNodes() {
	for _, node := range p.nodes {
		if node.NodeID != p.nodeID {
			server, err := rpc.DialHTTP("tcp", node.HostPort)
			if err != nil {
				p.connections = append(p.connections, server)
			}
		}
	}
}
