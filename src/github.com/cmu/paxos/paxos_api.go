package paxos

import (
	"github.com/cmu440/rpc/paxosrpc"
)

type Paxos interface {

	/*
	*Standard RecvPrepare implementation
	 */
	RecvPrepare(args *paxosrpc.PrepareArgs, reply *paxosrpc.PrepareReply) error
	/*
	*Standard RevcAccept implementation
	 */
	RecvAccept(args *paxosrpc.AcceptArgs, reply *paxosrpc.AcceptReply) error
	/*
	*RecvCommit implemtation that also builds in funcationality for catchup if a node is behind
	 */
	RecvCommit(args *paxosrpc.CommitArgs, reply *paxosrpc.CommitReply) error
	/*
	*Returns a list of active Paxos servers
	 */
	GetServers(args *paxosrpc.GetServerArgs, reply *paxosrpc.GetServerReply) error
	/*
	*Replaces the old node with the new node
	 */
	ReplaceNode(oldNode *paxosrpc.ReplaceNodeArgs, newNode *paxosrpc.ReplaceNodeReply) error
	/*
		Returns the current master server -- the active server with the lowest nodeID
	*/
	MasterServer(args *paxosrpc.GetMasterArgs, reply *paxosrpc.GetMasterReply) error
	/*
		Standard propose implementation. Does not block.
	*/
	Propose(args *paxosrpc.ProposeArgs, reply *paxosrpc.ProposeReply) error

	/*
	*Forces Paxos to stop accepting proposals and periodically polls the node
	* until it has flushed all commits. The node will then replace the old node
	* specified in the arguments with the new node specified in the arguments
	* and establish a connection. At that point, if the node is instructed
	* in the arguments (see rpc/paxosrpc/proto.go) to catchup the new node,
	* it will send commits to the new node until it is caught up. Upon which it
	* will return
	*
	 */
	Quiesce(args *paxosrpc.QuiesceArgs, reply *paxosrpc.QuiesceReply) error

	/*
	* Resume will end the quiesce period and allow nodes to start accepting proposals
	 */
	Resume(args *paxosrpc.ResumeArgs, reply *paxosrpc.ResumeReply) error
}
