package rpc

type RemotePaxosServer interface {
	Prepare(args *rpc.PrepareArgs, reply *rpc.PrepareReply) error

	Accept(args *rpc.AcceptArgs, reply *rpc.AcceptReply) error

	Commit(args *rpc.CommitArgs) error

	GetServers(args *rpc.GetServerArgs, reply *rpc.GetServerReply) error

	AddNode(oldNode *rpc.Node, newNode *rpc.Node) error

	MasterServer(args *rpc.GetMasterArgs, reply *rpc.GetMasterReply) error

	Propose(args *rpc.ProposeArgs, reply *rpc.ProposeReply)
}

type PaxosServer struct {
	// Embed all methods into the struct. See the Effective Go section about
	// embedding for more details: golang.org/doc/effective_go.html#embedding
	RemotePaxosServer
}

// Wrap wraps s in a type-safe wrapper struct to ensure that only the desired
// StorageServer methods are exported to receive RPCs.
func Wrap(s RemotePaxosServer) RemotePaxosServer {
	return &PaxosServer{s}
}
