package paxos

import (
	"github.com/cmu440/rpc/paxosrpc"
)

func (p *paxos) RecvCommit(args *paxosrpc.CommitArgs, reply *paxosrpc.CommitReply) error {
	if p.simulateNetworkError((*args).Committed.Sequence.Round) {
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
		p.learner.RecvCommit((*args).Committed.Value, false)
		p.contestedRound++
		p.noopRound = p.contestedRound
		p.catchup()
	} else if (*args).Committed.Sequence.Round > p.contestedRound {
		round, ok := p.commits[(*args).Committed.Sequence.Round]
		if ok {
			(*round).committed = true
			(*round).previous = (*args).Committed
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
