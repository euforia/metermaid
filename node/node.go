package node

import (
	"bytes"
	"encoding/binary"
	"strings"

	"github.com/hashicorp/memberlist"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

// Node holds information about a given matchine
type Node struct {
	// Node name
	Name string
	// Accessible address
	Address string
	// Total cpu shares in MHz
	CPUShares uint64
	// Total memory in bytes
	Memory uint64
	// Arbitrary node metadata including things like instance type. These
	// are used for grouping and aggregation queries.
	Meta map[string]string
}

// New computes the total cpu shares and memory of the system
// and returns a new Node instance
func New() *Node {
	cpus, _ := cpu.Info()
	// Get total for all cpus and cores
	var mhz float64
	for _, c := range cpus {
		mhz += c.Mhz * float64(c.Cores)
	}
	// Total memory for node
	m, _ := mem.VirtualMemory()

	return &Node{
		CPUShares: uint64(mhz),
		Memory:    m.Total,
	}
}

// NewFromMemberlistNode returns a Node constructed from a memberlist.Node
func NewFromMemberlistNode(in *memberlist.Node) *Node {
	n := &Node{
		Name:      in.Name,
		Address:   in.Address(),
		CPUShares: binary.BigEndian.Uint64(in.Meta[:8]),
		Memory:    binary.BigEndian.Uint64(in.Meta[8:16]),
	}
	if len(in.Meta[16:]) > 0 {
		taglist := bytes.Split(in.Meta[16:], []byte("\n"))
		if len(taglist) > 0 {
			n.Meta = make(map[string]string)
			for _, tagpair := range taglist {
				if len(tagpair) > 0 {
					kv := strings.Split(string(tagpair), "=")
					n.Meta[kv[0]] = kv[1]
				}
			}
		}
	}
	return n
}
