package paxos

import (
	"github.com/cmu440/rpc/paxosrpc"
	"log"
	"net/rpc"
	"time"
)

func (p *paxos) Propose(args *paxosrpc.ProposeArgs, reply *paxosrpc.ProposeReply) error {
	p.quiesceLock.Lock()
	if p.quiesce == true {
		(*reply).Status = paxosrpc.QUIESCE
		p.quiesceLock.Unlock()
		return nil
	}
	p.quiesceLock.Unlock()
	(*(reply)).Status = paxosrpc.OK
	proposal := (*args).Proposal
	p.listLock.Lock()
	if len(proposal) > 0 {
		p.proposalList.PushBack(proposal)
	}
	p.listLock.Unlock()
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
	cancel := 0 + ((p.numNodes - 1) - len(p.connections)) //If we have a dead node with a dead connection, we cound it out
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
	cancel := 0 + ((p.numNodes - 1) - len(p.connections)) //If we have a dead node with a dead connection, we cound it out
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
		p.learner.RecvCommit((*args).Committed.Value, ownValue)
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
			p.learner.RecvCommit(round.previous.Value, false)
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

func (p *paxos) simulateNetworkError(round uint64) bool {
	if p.debug == paxosrpc.NONE {
		return false
	} else if p.nodeID%2 == 0 {
		return false //Even numbered nodes won't drop messages
	}
	p.dataLock.Lock()
	defer p.dataLock.Unlock()
	if p.debug == paxosrpc.DROPODD && round%2 == 1 {
		return true
	} else if p.debug == paxosrpc.DROPSTART && round < 50 {
		return true
	}
	return false
}
