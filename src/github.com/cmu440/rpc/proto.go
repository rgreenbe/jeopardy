package rpc

type Status int

const (
	OK Status = iota + 1
	CANCEL
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
