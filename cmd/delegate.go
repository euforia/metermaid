package main

import (
	"bytes"

	"github.com/euforia/metermaid/tsdb"

	"github.com/euforia/metermaid/node"
	"github.com/euforia/metermaid/types"
	"github.com/hashicorp/memberlist"
	"go.uber.org/zap"
)

type GossipDelegate struct {
	node node.Node
	log  *zap.Logger
}

// LocalState satisfies the gossip.Delegate interface
func (del *GossipDelegate) LocalState(join bool) []byte {
	header := []byte(del.node.Meta.String() + "\n")
	if join {
		// return del.localStateOnJoin()
	}
	// return del.localStateUpdate()

	return header
}

// MergeRemoteState satisfies the gossip.Delegate interface
func (del *GossipDelegate) MergeRemoteState(data []byte, join bool) {
	i := bytes.IndexRune(data, '\n')
	if i < 0 {
		return
	}

	header := string(data[:i])
	meta := types.ParseMetaFromString(header)
	series := tsdb.Series{
		Meta: meta,
	}
	if series.Data == nil {

	}
	// if !meta.Equal(del.node.Meta) {
	// 	return
	// }

	if join {
		// del.mergeRemoteStateOnJoin()
	} else {
		// del.mergeRemoteStateUpdate()
	}
}

// NodeMeta satisfies the gossip.Delegate interface
func (del *GossipDelegate) NodeMeta(overhead int) []byte {
	return del.node.MarshalMeta()
}

// NotifyMsg satisfies the gossip.Delegate interface
func (del *GossipDelegate) NotifyMsg(msg []byte) {}

// NotifyJoin satisfies the memberlist.EventDelegate interface
func (del *GossipDelegate) NotifyJoin(nd *memberlist.Node) {
	n := newNode(nd)

	del.log.Info("node joined",
		zap.String("name", n.Name),
		zap.String("addr", n.Address),
		zap.String("tags", n.Meta.String()),
	)
}

// NotifyLeave satisfies the memberlist.EventDelegate interface
func (del *GossipDelegate) NotifyLeave(node *memberlist.Node) {}

// NotifyUpdate satisfies the memberlist.EventDelegate interface
func (del *GossipDelegate) NotifyUpdate(nd *memberlist.Node) {
	n := newNode(nd)
	del.log.Info("node updated",
		zap.String("name", n.Name),
		zap.String("addr", n.Address),
		zap.String("tags", n.Meta.String()),
	)
}

func newNode(in *memberlist.Node) *node.Node {
	n := &node.Node{
		Name:    in.Name,
		Address: in.Address(),
	}
	n.UnmarshalMeta(in.Meta)
	return n
}
