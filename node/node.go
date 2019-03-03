package node

import (
	"github.com/euforia/metermaid/types"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

// Nodes implements helper functions for groups of nodes
type Nodes []Node

// GroupBy returns a map of nodes grouped by the given Meta key
func (nodes Nodes) GroupBy(key string) map[string][]Node {
	grouped := make(map[string][]Node)
	for _, n := range nodes {
		keyV, ok := n.Meta[key]
		if !ok {
			continue
		}
		if vals, ok := grouped[keyV]; ok {
			grouped[keyV] = append(vals, n)
		} else {
			grouped[keyV] = []Node{n}
		}
	}
	return grouped
}

const (
	// PlatformAmazon is an amazon ec2 instance
	PlatformAmazon = "amazon"
)

// Platform holds the nodes platform information
type Platform struct {
	Name    string
	Family  string
	Version string
}

// Node holds information about a given matchine
type Node struct {
	// Node name
	Name string
	// Accessible address
	Address string
	// Total cpu shares in Hz
	CPUShares uint64
	// Total memory in bytes
	Memory uint64
	// Time the system booted
	BootTime uint64
	// OS and harware info
	Platform Platform
	// Arbitrary node metadata including things like instance type. These
	// are used for grouping and aggregation queries.
	Meta types.Meta
}

// CPUPercent returns the percent ratio of the given shares relative to the
// node. If the input is zero 1 is returned ie 100%
func (n *Node) CPUPercent(shares uint64) float64 {
	if shares != 0 {
		return float64(shares) / float64(n.CPUShares)
	}
	return 1
}

// MemoryPercent returns the percent ratio of the given mem relative to the
// node. If the input is zero 1 is returned ie 100%
func (n *Node) MemoryPercent(mem uint64) float64 {
	if mem != 0 {
		return float64(mem) / float64(n.Memory)
	}
	return 1
}

// IsAWSSpot returns true if this nodes metadata contains the appropriate
// spot tag key
func (n *Node) IsAWSSpot() bool {
	_, ok := n.Meta[SpotTag]
	return ok
}

// New computes the total cpu shares and memory of the system
// and returns a new Node instance
func New() (*Node, error) {
	cpus, _ := cpu.Info()
	// Get total for all cpus and cores
	var mhz float64
	for _, c := range cpus {
		// Convert to Hz
		mhz += c.Mhz * 1e6 * float64(c.Cores)
	}

	node := &Node{
		CPUShares: uint64(mhz),
	}
	// Total memory for node
	m, _ := mem.VirtualMemory()
	node.Memory = m.Total

	info, err := host.Info()
	if err == nil {
		// Convert to nanoseconds like everything else
		node.Name = info.HostID
		node.BootTime = info.BootTime * 1e9
		node.Platform = Platform{
			Name:    info.Platform,
			Family:  info.PlatformFamily,
			Version: info.PlatformVersion,
		}

		if mp := NewMetaProvider(node.Platform.Name); mp != nil {
			node.Meta, err = mp.Meta()
		}
	}

	return node, err
}

// NewWithMeta returns a new Node with the given metadata
func NewWithMeta(meta types.Meta) (*Node, error) {
	nd, err := New()
	if err == nil {
		if nd.Meta == nil {
			nd.Meta = make(types.Meta)
		}
		for k, v := range meta {
			nd.Meta[k] = v
		}
	}

	return nd, err
}

// NewWithMetaString returns a new node parsing the input metadata string
func NewWithMetaString(metastr string) (*Node, error) {
	if metastr != "" {
		meta := types.ParseMetaFromString(metastr)
		return NewWithMeta(meta)
	}
	return New()
}
