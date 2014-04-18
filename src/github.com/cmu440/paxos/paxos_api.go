package paxos

import "github.com/cmu440/rpc"

type Paxos interface {
	CreateNode() (Paxos, error)
	Prepare(args *rpc.PrepareArgs, reply *rpc.PrepareReply) error
	Accept(args *rpc.AcceptArgs, reply *rpc.AcceptReply) error
	Commit(args *rpc.CommitArgs)
}
