package paxos

import (
	"github.com/cmu440/rpc/paxosrpc"
)

func (p *paxos) RecvPrepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error {
	if p.simulateNetworkError((*args).Sequence.Round) {
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
	if p.simulateNetworkError((*args).Accept.Sequence.Round) {
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
