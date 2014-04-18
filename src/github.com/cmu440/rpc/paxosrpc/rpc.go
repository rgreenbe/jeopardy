package paxosrpc

type RemotePaxosServer interface {
	Prepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error

	Accept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error

	Commit(args *paxosrpc.CommitArgs) error

	GetServers(args *paxosrpc.GetServerArgs, reply *paxosrpc.GetServerReply) error

	AddNode(oldNode *paxosrpc.Node, newNode *paxosrpc.Node) error

	MasterServer(args *paxosrpc.GetMasterArgs, reply *paxosrpc.GetMasterReply) error

	Propose(args *paxosrpc.ProposeArgs, reply *paxosrpc.ProposeReply)
}

type PaxosServer struct {
	// Embed all methods into the struct. See the Effective Go section about
	// embedding for more details: golang.org/doc/effective_go.html#embedding
	RemotePaxosServer
}

// Wrap wraps s in a type-safe wrapper struct to ensure that only the desired
// StorageServer methods are exported to receive paxosrpcs.
func Wrap(s RemotePaxosServer) RemotePaxosServer {
	return &PaxosServer{s}
}
