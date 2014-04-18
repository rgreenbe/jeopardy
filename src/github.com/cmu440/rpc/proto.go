package rpc

type Status int

const (
	OK Status = iota + 1
	CANCEL
	NOT_READY
)

type PrepareArgs struct {
	ProposerSequence int
}

type PrepareReply struct {
	Status           Status
	Uncommitted      struct{}
	ProposerSequence int
}

type AcceptArgs struct {
	ProposerSequence int
	Uncommitted      struct{}
}

type AcceptReply struct {
	Status Status
}

type CommitArgs struct {
	Committed struct{}
}

type GetServerArgs struct {
	Node Node
}

type GetServerReply struct {
	Nodes []Node
}

type GetMasterArgs struct {
	//do nothing
}

type GetMasterReply struct {
	Node Node
}

type ProposeArgs struct {
	Proposal struct{}
}

type ProposeReply struct {
	Status Status
}

type Node struct {
	HostPort string
	NodeID   uint32
}
