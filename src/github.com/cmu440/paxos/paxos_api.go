package paxos

import (
	"github.com/cmu440/rpc/paxosrpc"
)

type Paxos interface {
	Prepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error

	Accept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error

	Commit(args *paxosrpc.CommitArgs) error

	GetServers(args *paxosrpc.GetServerArgs, reply *paxosrpc.GetServerReply) error

	AddNode(oldNode *paxosrpc.Node, newNode *paxosrpc.Node) error

	MasterServer(args *paxosrpc.GetMasterArgs, reply *paxosrpc.GetMasterReply) error

	Propose(args *paxosrpc.ProposeArgs, reply *paxosrpc.ProposeReply)
}
