package main

import (
	"encoding/binary"

	"github.com/euforia/metermaid"
	"github.com/hashicorp/memberlist"
	"go.uber.org/zap"
)

// type Proto uint8

// const (
// 	UDPProto Proto = iota + 1
// 	TCPProto
// 	HTTPProto
// 	HTTPSProto
// )

type GossipDelegate struct {
	node metermaid.Node
	log  *zap.Logger
}

// LocalState satisfies the gossip.Delegate interface
func (del *GossipDelegate) LocalState(join bool) []byte { return nil }

// MergeRemoteState satisfies the gossip.Delegate interface
func (del *GossipDelegate) MergeRemoteState(data []byte, join bool) {}

// NodeMeta satisfies the gossip.Delegate interface
func (del *GossipDelegate) NodeMeta(overhead int) []byte {
	cpu := make([]byte, 8)
	binary.BigEndian.PutUint64(cpu, del.node.CPUShares)
	mem := make([]byte, 8)
	binary.BigEndian.PutUint64(mem, del.node.Memory)
	return append(cpu, mem...)
}

// NotifyMsg satisfies the gossip.Delegate interface
func (del *GossipDelegate) NotifyMsg(msg []byte) {}

// NotifyJoin satisfies the memberlist.EventDelegate interface
func (del *GossipDelegate) NotifyJoin(node *memberlist.Node) {
	del.log.Info("node joined",
		zap.String("name", node.Name),
		zap.String("addr", node.Address()),
	)
}

// NotifyLeave satisfies the memberlist.EventDelegate interface
func (del *GossipDelegate) NotifyLeave(node *memberlist.Node) {}

// NotifyUpdate satisfies the memberlist.EventDelegate interface
func (del *GossipDelegate) NotifyUpdate(node *memberlist.Node) {
	del.log.Info("node updated")
}
