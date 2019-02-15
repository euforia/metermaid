package metermaid

import (
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/mem"
)

type NodeMetaProvider interface {
	Get(id string) map[string]string
}

type NodeCostProvider interface{}

// Node holds information about a given matchine
type Node struct {
	// CPUShares in MHz
	CPUShares uint64
	// Memory in bytes
	Memory uint64
	// Units is the abstact unit for the cost of the node
	UnitsPerMin float64
}

// NewNode computes the total cpu shares and memory of the system
// and returns a new Node instance
func NewNode() *Node {
	cpus, _ := cpu.Info()

	var mhz float64
	for _, c := range cpus {
		mhz += c.Mhz * float64(c.Cores)
	}

	m, _ := mem.VirtualMemory()

	return &Node{
		CPUShares: uint64(mhz),
		Memory:    m.Total,
	}
}
