package backend

type Backend interface {
	RecvCommit([]byte) error
}
