package paxos

import (
	"github.com/cmu440/rpc/paxosrpc"
	"net"
	"net/http"
	"net/rpc"
	"time"
)

type paxos struct {
	master *rpc.Client
	nodes  []paxosrpc.Node
	ready  paxosrpc.Status
}

func NewPaxos(masterHostPort string, numNodes int, hostPort string, nodeID uint32) (Paxos, error) {
	var listener net.Listener
	var err error
	for {
		listener, err = net.Listen("tcp", hostPort)
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 200) //Retry in a second
	}
	p := &paxos{nil, nil, paxosrpc.NOT_READY}
	for {
		err = rpc.RegisterName("Paxos", paxosrpc.Wrap(p))
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
	if masterHostPort != "" { //It's the master
		var server *rpc.Client
		for {
			server, err = rpc.DialHTTP("tcp", masterHostPort)
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
				p.nodes = reply.Nodes
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

func (p *paxos) Prepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error {
	return nil
}

func (p *paxos) Accept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error {
	return nil
}

func (p *paxos) Commit(args *paxosrpc.CommitArgs) error {
	return nil
}

func (p *paxos) GetServers(args *paxosrpc.GetServerArgs, reply *paxosrpc.GetServerReply) error {
	return nil
}

func (p *paxos) AddNode(oldNode *paxosrpc.Node, newNode *paxosrpc.Node) error {
	return nil
}

func (p *paxos) MasterServer(args *paxosrpc.GetMasterArgs, reply *paxosrpc.GetMasterReply) error {
	return nil
}

func (p *paxos) Propose(args *paxosrpc.ProposeArgs, reply *paxosrpc.ProposeReply) {
	return
}
