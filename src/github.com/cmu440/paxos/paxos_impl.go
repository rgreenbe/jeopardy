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
	"strconv"
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
}

func NewPaxos(masterHostPort string, numNodes int, hostPort string, nodeID, masterID uint64, learner backend.Backend) (Paxos, error) {
	var listener net.Listener
	var err error
	for {
		listener, err = net.Listen("tcp", hostPort)
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 200) //Retry in a bit
	}
	p := &paxos{nil, 0, 0, false, numNodes, nodeID, masterID, nil, list.New(), make(chan struct{}, 1000), nil,
		make([]*rpc.Client, 0, numNodes-1), learner, make(map[uint64]*Round), new(sync.Mutex), new(sync.Mutex)}
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

func (p *paxos) RecvPrepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error {
	if p.nodeID == 2 && (*args).Sequence.Round%2 == 1 {
		return nil
	}
	p.dataLock.Lock()
	defer p.dataLock.Unlock()
	roundNum := (*args).Sequence.Round
	round, ok := p.commits[roundNum]
	if ok { //Somone else has prepared this round
		if p.valid(round.highestSequence, (*args).Sequence) {
			(*reply).Status = paxosrpc.OK
			(*reply).Previous = (*round).previous
			(*round).highestSequence = (*args).Sequence
			if p.valid(p.highestSequence, (*args).Sequence) {
				p.highestSequence = (*args).Sequence
			}
		} else {
			(*reply).Status = paxosrpc.CANCEL
		}
	} else {
		p.commits[roundNum] = &Round{(*args).Sequence, nil, false}
		(*reply).Status = paxosrpc.OK
	}
	return nil
}

func (p *paxos) RecvAccept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error {
	if p.nodeID == 2 && (*args).Accept.Sequence.Round%2 == 1 {
		return nil
	}
	p.dataLock.Lock()
	defer p.dataLock.Unlock()
	roundNum := (*args).Accept.Sequence.Round
	round, ok := p.commits[roundNum]
	if ok {
		if p.valid(round.highestSequence, (*args).Accept.Sequence) {
			(*reply).Status = paxosrpc.OK
			round.previous = (*args).Accept
			round.highestSequence = (*args).Accept.Sequence
			p.commits[roundNum] = round
			if p.valid(p.highestSequence, (*args).Accept.Sequence) {
				p.highestSequence = (*args).Accept.Sequence
			}
		} else {
			(*reply).Status = paxosrpc.CANCEL
		}
	}
	return nil
}

func (p *paxos) RecvCommit(args *paxosrpc.CommitArgs, reply *paxosrpc.CommitReply) error {
	if p.nodeID == 2 && (*args).Committed.Sequence.Round%2 == 1 {
		return nil
	}
	p.dataLock.Lock()
	defer p.dataLock.Unlock()
	if p.highestSequence == nil || p.compare(p.highestSequence, (*args).Committed.Sequence) == LESS {
		p.highestSequence = (*args).Committed.Sequence
	}
	if (*args).Committed.Sequence.Round == p.contestedRound {
		round, ok := p.commits[(*args).Committed.Sequence.Round]
		if ok {
			(*round).committed = true
		} else {
			p.commits[(*args).Committed.Sequence.Round] = &Round{(*args).Committed.Sequence, (*args).Committed, true}
		}
		p.learner.RecvCommit([]byte(string((*args).Committed.Value) + " " + strconv.FormatUint(p.contestedRound, 10) + " " + string(p.nodeID)))
		p.contestedRound++
		p.noopRound = p.contestedRound
		p.catchup()
	} else if (*args).Committed.Sequence.Round > p.contestedRound {
		round, ok := p.commits[(*args).Committed.Sequence.Round]
		if ok {
			(*round).committed = true
		} else {
			p.commits[(*args).Committed.Sequence.Round] = &Round{(*args).Committed.Sequence, (*args).Committed, true}
		}
		if (*args).Committed.Sequence.Round == p.noopRound {
			//We will record the commit and send to learner when caught up
			p.noopRound++
		} else {
			//We need to do some catchup because we are behind
			var i uint64
			var diff uint64 = (*args).Committed.Sequence.Round - p.noopRound
			for i = 0; i < diff; i++ {
				p.noopRound++
				go p.Propose(new(paxosrpc.ProposeArgs), new(paxosrpc.ProposeReply))
			}
			p.noopRound++
		}
	}
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

func (p *paxos) ReplaceNode(args *paxosrpc.ReplaceNodeArgs, reply *paxosrpc.ReplaceNodeReply) error {
	oldNode := args.OldNode
	newNode := args.NewNode
	dummyReply := new(paxosrpc.CommitReply)
	server, err := rpc.DialHTTP("tcp", newNode.HostPort)
	if err != nil {
		(*reply).Done = false
	}
	for index, node := range p.nodes {
		if node.NodeID == (oldNode).NodeID {
			p.nodes[index] = newNode
			for _, commit := range p.commits {
				server.Call("Paxos.RecvCommit", commit, dummyReply)
			}
			(*reply).Done = true
			return nil
		}
	}
	(*reply).Done = false
	return errors.New("Old node does not exist")
}

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

func (p *paxos) Propose(args *paxosrpc.ProposeArgs, reply *paxosrpc.ProposeReply) error {
	(*(reply)).Status = paxosrpc.OK
	proposal := (*args).Proposal
	p.listLock.Lock()
	if proposal != nil {
		p.proposalList.PushBack(*proposal)
	} else {
		//p.proposalList.PushFront(nil)
	}
	p.listLock.Unlock()
	p.startPrepare <- struct{}{}
	return nil
}

func (p *paxos) Quiesce(args *paxosrpc.QuiesceArgs, reply *paxosrpc.QuiesceReply) error {
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

func (p *paxos) getN() uint64 {
	if p.highestSequence != nil {
		return p.highestSequence.N + 1
	} else {
		return 1
	}
}

func (p *paxos) sendPrepare() {
	p.dataLock.Lock()
	sequence := &paxosrpc.Sequence{p.contestedRound, p.getN(), p.nodeID}
	if !p.madeConnections {
		p.connectToNodes()
	}
	args := &paxosrpc.PrepareArgs{sequence}
	replyChan := make(chan paxosrpc.PrepareReply, p.numNodes-1)
	for _, connection := range p.connections {
		go p.rpcPrepare(connection, args, replyChan)
	}
	go p.rpcPrepare(nil, args, replyChan)
	var oldestPrepare *paxosrpc.ValueSequence
	ok := 0
	cancel := 0
	p.dataLock.Unlock()
	for i := 0; i < p.numNodes; i++ {
		reply := <-replyChan
		if reply.Status == paxosrpc.OK {
			if reply.Previous != nil {
				if oldestPrepare == nil {
					oldestPrepare = reply.Previous
				} else if p.compare(oldestPrepare.Sequence, (*reply.Previous).Sequence) == GREATER {
					oldestPrepare = reply.Previous
				}
			}
			ok++
			if (p.numNodes / 2) < ok {
				if oldestPrepare == nil {
					p.listLock.Lock()
					var value []byte
					front := p.proposalList.Front()
					if front != nil {
						value = front.Value.([]byte)
					} else {
						return
					}
					p.listLock.Unlock()
					p.sendAccept(&paxosrpc.ValueSequence{value, sequence}, true)
				} else {
					(*oldestPrepare).Sequence = sequence
					p.sendAccept(oldestPrepare, false)
				}
				return
			}
		} else {
			cancel++
			if (p.numNodes / 2) < cancel {
				time.Sleep(time.Millisecond * 100)
				p.startPrepare <- struct{}{} //So sorry..try again
				return
			}
		}
	}
}

func (p *paxos) sendAccept(accept *paxosrpc.ValueSequence, ownValue bool) {
	p.dataLock.Lock()
	args := &paxosrpc.AcceptArgs{accept}
	replyChan := make(chan paxosrpc.Status, p.numNodes-1)
	p.dataLock.Unlock()
	for _, connection := range p.connections {
		go p.rpcAccept(connection, args, replyChan)
	}
	go p.rpcAccept(nil, args, replyChan)
	ok := 0
	cancel := 0
	for i := 0; i < p.numNodes; i++ {
		reply := <-replyChan
		if reply == paxosrpc.OK {
			ok++
			if (p.numNodes / 2) < ok {
				p.sendCommit(accept, ownValue)
				return
			}
		} else {
			cancel++
			if (p.numNodes / 2) < cancel {
				time.Sleep(time.Millisecond * 100)
				p.startPrepare <- struct{}{} //So sorry..try again
				return
			}
		}
	}
}

func (p *paxos) sendCommit(commit *paxosrpc.ValueSequence, ownValue bool) {
	args := &paxosrpc.CommitArgs{commit}
	for _, connection := range p.connections {
		go p.rpcCommit(connection, args)
	}
	p.listLock.Lock()
	if ownValue {
		p.proposalList.Remove(p.proposalList.Front())
	} else if p.proposalList.Len() > 0 {
		p.startPrepare <- struct{}{}
	}
	p.listLock.Unlock()
	p.dataLock.Lock()
	if p.contestedRound == (*commit).Sequence.Round {
		p.learner.RecvCommit([]byte(string((*args).Committed.Value) + " proposer" + strconv.FormatUint(p.nodeID, 10) + " round " + strconv.FormatUint((*args).Committed.Sequence.Round, 10)))
		p.commits[(*commit).Sequence.Round] = &Round{(*commit).Sequence, commit, true}
		p.contestedRound++
		p.catchup()
	}
	p.dataLock.Unlock()
}

func (p *paxos) catchup() {
	for {
		round, ok := p.commits[p.contestedRound]
		if ok && round.committed {
			p.contestedRound++
			p.learner.RecvCommit(round.previous.Value)
		} else {
			return
		}
	}
}

func (p *paxos) rpcCommit(server *rpc.Client, args *paxosrpc.CommitArgs) {
	reply := new(paxosrpc.CommitReply)
	server.Call("Paxos.RecvCommit", args, reply)
}

func (p *paxos) rpcAccept(server *rpc.Client, args *paxosrpc.AcceptArgs, replyChan chan paxosrpc.Status) {
	reply := new(paxosrpc.AcceptReply)
	if server != nil {
		server.Call("Paxos.RecvAccept", args, reply)
	} else {
		p.RecvAccept(args, reply)
	}
	replyChan <- (*reply).Status
}

func (p *paxos) rpcPrepare(server *rpc.Client, args *paxosrpc.PrepareArgs, replyChan chan paxosrpc.PrepareReply) {
	reply := new(paxosrpc.PrepareReply)
	if server != nil {
		server.Call("Paxos.RecvPrepare", args, reply)
	} else {
		p.RecvPrepare(args, reply)
	}
	replyChan <- *reply
}

func (p *paxos) valid(highest, prepare *paxosrpc.Sequence) bool {
	if highest == nil {
		return true
	} else if p.compare(highest, prepare) == LESS {
		return true
	} else if p.compare(highest, prepare) == EQUAL {
		return true
	}
	return false
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
	p.connections = make([]*rpc.Client, 0, p.numNodes-1)
	for _, node := range p.nodes {
		if node.NodeID != p.nodeID {
			server, err := rpc.DialHTTP("tcp", node.HostPort)
			if err == nil {
				p.connections = append(p.connections, server)
			} else {
				log.Println(err)
			}
		}
	}
	p.madeConnections = true
}
