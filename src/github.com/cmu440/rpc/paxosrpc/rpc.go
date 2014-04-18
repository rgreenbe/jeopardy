package paxosrpc

type RemotePaxosServer interface {
	Prepare(args *PrepareArgs, reply *PrepareReply) error

	Accept(args *AcceptArgs, reply *AcceptReply) error

	Commit(args *CommitArgs) error

	GetServers(args *GetServerArgs, reply *GetServerReply) error

	AddNode(oldNode *Node, newNode *Node) error

	MasterServer(args *GetMasterArgs, reply *GetMasterReply) error

	Propose(args *ProposeArgs, reply *ProposeReply)
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
