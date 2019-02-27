package collector

import (
	"context"
	"time"

	"go.uber.org/zap"

	"github.com/euforia/metermaid/types"
)

const (
	// ResourceNode ...
	ResourceNode = "node"
	// ResourceContainer ...
	ResourceContainer = "container"
)

// RunStats holds the runtime of a given resource.
type RunStats struct {
	// Resource is the type of resource this Runtime represents
	Resource string
	// Cpu allocated or utilized by the process
	CPU uint64
	// Memory allocated or utilized by the process
	Memory uint64

	Meta types.Meta

	Start time.Time
	End   time.Time
}

// Duration returns the duration of the runtime
func (rt *RunStats) Duration() time.Duration {
	if rt.End.UnixNano() <= 0 {
		return time.Now().Sub(rt.Start)
	}
	return rt.End.Sub(rt.Start)
}

// Collector ...
type Collector interface {
	Name() string
	// Init should initialize the Collector returning an error on
	// any failure signaling the Collector not to be loaded
	Init(map[string]interface{}) error
	// Collection should return runtimes of resources.  End is not
	// used as it will be filled in with the current time on each
	// invocation
	Collect(context.Context) ([]RunStats, error)
}

type collector struct {
	interval time.Duration
	bc       Collector
	out      chan<- []RunStats

	log *zap.Logger
}

func (c *collector) run(ctx context.Context) {
	data, err := c.bc.Collect(ctx)
	if err == nil {
		c.out <- data
	} else {
		c.log.Info("collection failed", zap.String("name", c.bc.Name()), zap.Error(err))
	}

	timer := time.NewTimer(c.interval)
	c.log.Info("collector starting", zap.String("name", c.bc.Name()), zap.Duration("interval", c.interval))
	for {
		select {
		case <-timer.C:
			data, err := c.bc.Collect(ctx)
			if err == nil {
				c.out <- data
			} else {
				c.log.Info("collection failed", zap.String("name", c.bc.Name()), zap.Error(err))
			}
			// c.log.Debug("resetting", zap.String("name", c.bc.Name()))
			timer.Reset(c.interval)

		case <-ctx.Done():
			c.log.Info("collector scheduler stopped", zap.String("name", c.bc.Name()))
			timer.Stop()
			return
		}

	}
}

// schedule start the process of scheduling collections. Each collection
// is rescheduled upon each run
// func (c *collector) schedule() {
// 	data, err := c.bc.Collect(c.ctx)
// 	if err == nil {
// 		c.out <- data
// 	}

// 	select {
// 	case <-c.ctx.Done():

// 	default:
// 		time.AfterFunc(c.interval, c.schedule)
// 		c.log.Debug("collector rescheduled",
// 			zap.String("name", c.bc.Name()),
// 			zap.Duration("interval", c.interval),
// 		)

// 	}
// }
