package paxos

import (
	"github.com/cmu440/rpc/paxosrpc"
)

type Paxos interface {
	RecvPrepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error

	RecvAccept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error

	RecvCommit(args *paxosrpc.CommitArgs, reply *paxosrpc.CommitReply) error

	GetServers(args *paxosrpc.GetServerArgs, reply *paxosrpc.GetServerReply) error

	ReplaceNode(oldNode *paxosrpc.ReplaceNodeArgs, newNode *paxosrpc.ReplaceNodeReply) error

	MasterServer(args *paxosrpc.GetMasterArgs, reply *paxosrpc.GetMasterReply) error

	Propose(args *paxosrpc.ProposeArgs, reply *paxosrpc.ProposeReply) error

	Quiesce(args *paxosrpc.QuiesceArgs, reply *paxosrpc.QuiesceReply) error

	Resume(args *paxosrpc.ResumeArgs, reply *paxosrpc.ResumeReply) error
}
