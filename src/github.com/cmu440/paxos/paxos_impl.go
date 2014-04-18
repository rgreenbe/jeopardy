package paxos

import (
	"errors"
	"github.com/cmu440/rpc"
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
		err = rpc.RegisterName("Paxos", rpc.Wrap(p))
		if err == nil {
			break
		}
		time.Sleep(time.Millisecond * 200)
	}
	if master != nil { //It's the master
		var server *rpc.Client
		for {
			server, err = rpc.DialHTTP("tcp", master.HostPort)
			if err == nil {
				break
			}
			time.Sleep(time.Millisecond * 200)
		}
		reply := new(rpc.GetServerReply)
		args := &rpc.GetServerArgs{rpc.Node{hostPort, nodeID}}
		for {
			server.Call("Paxos.GetServers", args, reply)
			if reply.Status == rpc.OK {
				p.nodes = reply.Servers
				break
			}
			time.Sleep(time.Second)
		}
		p.master = server
	} else {
		p.nodes = make([]rpc.Node, 0, numNodes)
		p.nodes = append(ss.nodes, rpc.Node{hostPort, nodeID})
		if numNodes == 1 {
			p.ready = rpc.OK
		} else {
			p.ready = rpc.NOT_READY
		}
	}
	rpc.HandleHTTP()
	go http.Serve(listener, nil)
	return p, nil
}

func Prepare(args *rpc.PrepareArgs, reply *rpc.PrepareReply) error {

}

func Accept(args *rpc.AcceptArgs, reply *rpc.AcceptReply) error {

}

func Commit(args *rpc.CommitArgs) error {

}

func GetServers(args *rpc.GetServerArgs, reply *rpc.GetServerReply) error {

}

func AddNode(oldNode *rpc.Node, newNode *rpc.Node) error {

}

func MasterServer(args *rpc.GetMasterArgs, reply *rpc.GetMasterReply) error {

}

func Propose(args *rpc.ProposeArgs, reply *rpc.ProposeReply) {

}
