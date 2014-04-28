package paxosrpc

type Status int

const (
	OK Status = iota + 1
	CANCEL
	NOT_READY
)

type PrepareArgs struct {
	Sequence *Sequence
}

type PrepareReply struct {
	Status   Status
	Previous *ValueSequence
}

type AcceptArgs struct {
	Accept *ValueSequence
}

type AcceptReply struct {
	Status Status
}

type CommitArgs struct {
	Committed *ValueSequence
}

type CommitReply struct {
	//do nothing
}

type GetServerArgs struct {
	Node *Node
}

type GetServerReply struct {
	Nodes  []Node
	Status Status
}

type GetMasterArgs struct {
	//do nothing
}

type GetMasterReply struct {
	Node Node
}

type ProposeArgs struct {
	Proposal *[]byte
}

type ProposeReply struct {
	Status Status
}
type ReplaceNodeArgs struct {
	Update  bool
	OldNode Node
	NewNode Node
}
type ReplaceNodeReply struct {
	Done bool
}

type QuiesceArgs struct {
}

type QuiesceReply struct {
}
