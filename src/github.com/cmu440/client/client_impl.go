package client

import (
	"errors"
	"github.com/cmu440/rpc/paxosrpc"
)

type jeopardyClient struct {
	client *rpc.Client
}

func NewJeopardyClient(serverHost string, serverPort int) (jeopardyClient, error) {
	cli, err := rpc.DialHTTP("tcp", net.JoinHostPort(serverHost, strconv.Itoa(serverPort)))
	if err != nil {
		return nil, err
	}
	return &tribClient{client: cli}, nil
}
func (jc *jeopardyClient) SendMessage(userMessage string) {
	args := paxosrpc.Propose{[]byte(userMessage)}

}
