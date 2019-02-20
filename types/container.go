package types

import (
	"time"

	"github.com/euforia/metermaid/fl"
)

// Container holds information about a container that is or was running
type Container struct {
	ID        string
	Name      string
	Create    int64 // epoch nano
	Start     int64 // epoch nano
	Stop      int64 // epoch nano
	Destroy   int64 // epoch nano
	Memory    int64 // bytes
	CPUShares int64 // MHz?
	Labels    map[string]string
	Tags      map[string]string
	// Units used.  This can be dollars or any other
	// virtual unit
	UnitsBurned float64
}

// Destroyed returns true if the container has been destroyed
func (cont *Container) Destroyed() bool {
	return cont.Destroy > 0
}

// RunTime returns duration for which the container was actually running
func (cont *Container) RunTime() time.Duration {
	return delta(cont.Stop, cont.Start)
}

// AllocatedTime returns the amount of time container resources were allocated
// i.e from the time it was created to the time it was completely destroyed
func (cont *Container) AllocatedTime() time.Duration {
	return delta(cont.Destroy, cont.Create)
}

func delta(end, start int64) time.Duration {
	if d := end - start; d > -1 {
		return time.Duration(d)
	}
	return time.Duration(0)
}

func (cont *Container) Match(query fl.Query) bool {
	for k, q := range query {
		if !cont.MatchField(k, q...) {
			return false
		}
	}
	return true
}

// MatchField returns true if the name of the given field match the filters.  filters are
// treated as AND's
func (cont *Container) MatchField(name string, filters ...fl.Filter) bool {
	switch name {
	case "Name":
		for _, filter := range filters {
			if fl.MatchString(cont.Name, filter) {
				continue
			}
			return false
		}
		return true
	case "Create":
		for _, filter := range filters {
			if fl.MatchTime(time.Unix(0, cont.Create), filter) {
				continue
			}
			return false
		}
		return true

	case "Start":
		for _, filter := range filters {
			if fl.MatchTime(time.Unix(0, cont.Start), filter) {
				continue
			}
			return false
		}
		return true
	case "Stop":
		for _, filter := range filters {
			if fl.MatchTime(time.Unix(0, cont.Stop), filter) {
				continue
			}
			return false
		}
		return true
	case "Destroy":
		for _, filter := range filters {
			if fl.MatchTime(time.Unix(0, cont.Destroy), filter) {
				continue
			}
			return false
		}
		return true
	case "CPUShares":
		for _, filter := range filters {
			if fl.MatchInt64(cont.CPUShares, filter) {
				continue
			}
			return false
		}
		return true
	case "Memory":
		for _, filter := range filters {
			if fl.MatchInt64(cont.Memory, filter) {
				continue
			}
			return false
		}
		return true

	default:
		for _, filter := range filters {
			if cont.MatchLabel(name, filter) {
				continue
			}
			return false
		}
		return true
	}
}

func (cont *Container) MatchLabel(name string, filter fl.Filter) bool {
	val, ok := cont.Labels[name]
	if !ok {
		return false
	}
	return fl.MatchString(val, filter)
}
