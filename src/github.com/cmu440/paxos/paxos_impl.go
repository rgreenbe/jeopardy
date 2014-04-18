package paxos

import (
	"errors"
	"github.com/cmu440/rpc/paxosrpc"
	"net"
	"net/rpc"
)

type paxos struct {
	master *rpc.Client
	nodes  []Node
	ready  Status
}

func NewPaxos(master Node, numNodes int, hostPort string, nodeID uint32) (Paxos, error) {
	var listener net.Listener
	var err error
	for {
		listener, err = net.Listen("tcp", hostPort)
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 200) //Retry in a second
	}
	p := &paxos{}
	for {
		err = rpc.RegisterName("Paxos", paxosrpc.Wrap(p))
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
	if master != nil { //It's the master
		var server *paxosrpc.Client
		for {
			server, err = rpc.DialHTTP("tcp", master.HostPort)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 200)
		}
		reply := new(paxosrpc.GetServerReply)
		args := &paxosrpc.GetServerArgs{paxosrpc.Node{hostPort, nodeID}}
		for {
			server.Call("Paxos.GetServers", args, reply)
			if reply.Status == paxosrpc.OK {
				p.nodes = reply.Servers
				break
			}
			time.Sleep(time.Second)
		}
		p.master = server
	} else {
		p.nodes = make([]paxosrpc.Node, 0, numNodes)
		p.nodes = append(p.nodes, paxosrpc.Node{hostPort, nodeID})
		if numNodes == 1 {
			p.ready = paxosrpc.OK
		} else {
			p.ready = paxosrpc.NOT_READY
		}
	}
	rpc.HandleHTTP()
	go http.Serve(listener, nil)
	return p, nil
}

func Prepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error {

}

func Accept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error {

}

func Commit(args *paxosrpc.CommitArgs) error {

}

func GetServers(args *paxosrpc.GetServerArgs, reply *paxosrpc.GetServerReply) error {

}

func AddNode(oldNode *paxosrpc.Node, newNode *paxosrpc.Node) error {

}

func MasterServer(args *paxosrpc.GetMasterArgs, reply *paxosrpc.GetMasterReply) error {

}

func Propose(args *paxosrpc.ProposeArgs, reply *paxosrpc.ProposeReply) {

}
