package node

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"strings"

	"github.com/euforia/metermaid/fl"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"
	"github.com/shirou/gopsutil/mem"
)

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
	// Total cpu shares in MHz
	CPUShares uint64
	// Total memory in bytes
	Memory uint64
	// Time the system booted
	BootTime uint64
	// OS and harware info
	Platform Platform
	// Arbitrary node metadata including things like instance type. These
	// are used for grouping and aggregation queries.
	Meta map[string]string
}

func (n *Node) Match(query fl.Query) bool {
	for k, filters := range query {
		if n.MatchMeta(k, filters...) {
			continue
		}
		return false
	}
	return true
}

func (n *Node) MatchMeta(name string, filters ...fl.Filter) bool {
	val, ok := n.Meta[name]
	if !ok {
		return false
	}

	for _, filter := range filters {
		if fl.MatchString(val, filter) {
			continue
		}
		return false
	}
	return true
}

func (n *Node) MarshalMeta() []byte {
	bt := make([]byte, 8)
	binary.BigEndian.PutUint64(bt, n.BootTime)
	cpu := make([]byte, 8)
	binary.BigEndian.PutUint64(cpu, n.CPUShares)
	mem := make([]byte, 8)
	binary.BigEndian.PutUint64(mem, n.Memory)

	data := append(cpu, mem...)
	data = append(bt, data...)

	data = append(data, []byte(n.Platform.Name+"\n")...)
	data = append(data, []byte(n.Platform.Family+"\n")...)
	data = append(data, []byte(n.Platform.Version+"\n")...)

	for k, v := range n.Meta {
		data = append(data, []byte(k+"="+v+"\n")...)
	}

	return data
}

// CPUPercent returns the percent ratio of the given shares relative to the
// node
func (n *Node) CPUPercent(shares uint64) float64 {
	return float64(shares) / float64(n.CPUShares)
}

// MemoryPercent returns the percent ratio of the given mem relative to the
// node
func (n *Node) MemoryPercent(mem uint64) float64 {
	return float64(mem) / float64(n.Memory)
}

func (n *Node) UnmarshalMeta(meta []byte) {
	n.BootTime = binary.BigEndian.Uint64(meta[:8])
	n.CPUShares = binary.BigEndian.Uint64(meta[8:16])
	n.Memory = binary.BigEndian.Uint64(meta[16:24])

	if len(meta[24:]) > 0 {
		list := bytes.Split(meta[24:], []byte("\n"))
		n.Platform.Name = string(list[0])
		n.Platform.Family = string(list[1])
		n.Platform.Version = string(list[2])

		if len(list) > 3 {
			n.Meta = make(map[string]string)
			for _, tagpair := range list[3:] {
				if len(tagpair) > 0 {
					kv := strings.Split(string(tagpair), "=")
					n.Meta[kv[0]] = kv[1]
				}
			}
		}
	}
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

	node := &Node{
		CPUShares: uint64(mhz),
	}
	// Total memory for node
	m, _ := mem.VirtualMemory()
	node.Memory = m.Total

	info, err := host.Info()
	if err == nil {
		// Convert to nanoseconds like everything else
		fmt.Println(info.BootTime)
		node.BootTime = info.BootTime * 1e9
		node.Platform = Platform{
			Name:    info.Platform,
			Family:  info.PlatformFamily,
			Version: info.PlatformVersion,
		}
	}

	return node
}
