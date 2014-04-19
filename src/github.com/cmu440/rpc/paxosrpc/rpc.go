package paxosrpc

type RemotePaxosServer interface {
	RecvPrepare(args *PrepareArgs, reply *PrepareReply) error

	RecvAccept(args *AcceptArgs, reply *AcceptReply) error

	RecvCommit(args *CommitArgs, reply *CommitReply) error

	GetServers(args *GetServerArgs, reply *GetServerReply) error

	ReplaceNode(oldNode *Node, newNode *Node) error

	MasterServer(args *GetMasterArgs, reply *GetMasterReply) error

	Propose(args *ProposeArgs, reply *ProposeReply) error
}

type PaxosServer struct {
	// Embed all methods into the struct. See the Effective Go section about
	// embedding for more details: golang.org/doc/effective_go.html#embedding
	RemotePaxosServer
}

type ValueSequence struct {
	Value    struct{}
	Sequence *Sequence
}

type Node struct {
	HostPort string
	NodeID   uint64
}

type Sequence struct {
	N      uint64
	NodeID uint64
}

// Wrap wraps s in a type-safe wrapper struct to ensure that only the desired
// StorageServer methods are exported to receive paxosrpcs.
func Wrap(s RemotePaxosServer) RemotePaxosServer {
	return &PaxosServer{s}
}
