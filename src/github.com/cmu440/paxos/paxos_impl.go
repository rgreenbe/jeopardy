package paxos

import (
	"container/list"
	"errors"
	"github.com/cmu440/rpc/paxosrpc"
	"math"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

type paxos struct {
	numNodes        int
	master          *rpc.Client
	proposalList    *list.List
	startPrepare    chan struct{}
	nodes           []paxosrpc.Node
	highestSequence *paxosrpc.Sequence
	previous        *paxosrpc.ValueSequence
}

func NewPaxos(masterHostPort string, numNodes int, hostPort string, nodeID uint32) (Paxos, error) {
	var listener net.Listener
	var err error
	for {
		listener, err = net.Listen("tcp", hostPort)
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 200) //Retry in a second
	}
	p := &paxos{numNodes, nil, make(chan struct{}, 1000), list.New(), make(chan struct{}), nil, nil, nil}
	for {
		err = rpc.RegisterName("Paxos", paxosrpc.Wrap(p))
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
	if masterHostPort != "" { //It's the master
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
	} else {
		p.nodes = make([]paxosrpc.Node, 0, numNodes)
		p.nodes = append(p.nodes, paxosrpc.Node{hostPort, nodeID})
	}
	rpc.HandleHTTP()
	go p.submitPrepare()
	go http.Serve(listener, nil)
	return p, nil
}

func (p *paxos) RecvPrepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error {
	if p.isGreater((*args).Sequence) {
		(*reply).Status = paxosrpc.OK
		(*reply).Previous = p.previous
		p.highestSequence = args.Sequence
	} else {
		(*reply).Status = paxosrpc.CANCEL
	}
	return nil
}

func (p *paxos) RecvAccept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error {
	if p.isGreater((*args).Accept.Sequence) {
		(*reply).Status = paxosrpc.OK
		p.previous = (*args).Accept
	} else {
		(*reply).Status = paxosrpc.CANCEL
	}
	return nil
}

func (p *paxos) RecvCommit(args *paxosrpc.CommitArgs) error {
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
			break
		}
	}
	return errors.New("Node does not exist")
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
		var min uint32 = math.MaxUint32
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

func (p *paxos) Propose(args *paxosrpc.ProposeArgs, reply *paxosrpc.ProposeReply) {
	proposal := args.Proposal
	(*(reply)).Status = paxosrpc.OK
	p.proposalChan <- proposal
	p.startPrepare <- struct{}{}
	return
}

func (p *paxos) submitPrepare() {
	for {
		select {
		case <-p.startPrepare:

		}
	}
}

func (p *paxos) isGreater(prepare *paxosrpc.Sequence) bool {
	if p.highestSequence.N < prepare.N {
		return prepare.NodeID < p.highestSequence.NodeID
	}
	return false
}
