package client

import (
	"github.com/cmu440/rpc/paxosrpc"
)

type JeopardyClient interface {
	SendMessage(message string) error
}
