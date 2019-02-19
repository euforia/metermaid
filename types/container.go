package types

import (
	"time"
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
