package backend

type Backend interface {
	RecvCommit([]byte, bool) error
}
