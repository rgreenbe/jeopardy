package paxos

import "github.com/cmu440/rpc"

type Paxos interface {
	CreateNode() (*Paxos, error)
	Prepare(args *rpc.PrepareArgs, reply *rpc.PrepareReply) error
	Accept(args *rpc.AcceptArgs, reply *rpc.AcceptReply) error
	Commit(args *rpc.CommitArgs) error
	ActiveNodes(args *rpc.GetServerArgs, reply *rpc.GetServerReply) error
	AddNode(oldNode *rpc.Node, newNode *rpc.Node) error
	MasterServer(args *rpc.GetMasterArgs, reply *rpc.GetMasterReply) error
}
