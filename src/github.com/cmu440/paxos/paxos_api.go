package paxos

import (
	"github.com/cmu440/rpc/paxosrpc"
)

type Paxos interface {
	RecvPrepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error

	RecvAccept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error

	RecvCommit(args *paxosrpc.CommitArgs) error

	GetServers(args *paxosrpc.GetServerArgs, reply *paxosrpc.GetServerReply) error

	ReplaceNode(oldNode *paxosrpc.Node, newNode *paxosrpc.Node) error

	MasterServer(args *paxosrpc.GetMasterArgs, reply *paxosrpc.GetMasterReply) error

	Propose(args *paxosrpc.ProposeArgs, reply *paxosrpc.ProposeReply)
}
