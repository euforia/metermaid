package collector

import (
	"context"
	"errors"
	"fmt"
	"time"

	"go.uber.org/zap"

	"github.com/euforia/metermaid/node"
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
	// Arbitrary metadata used for querying
	Meta types.Meta
	// Process start time
	Start time.Time
	// Optional end time
	End time.Time
}

// Config holds a collector configuration
type Config struct {
	// Node information
	Node *node.Node
	// Run interval
	Interval time.Duration
	// Custom config
	Config map[string]interface{}
	// logger populated by engine
	Logger *zap.Logger
}

// Validate returns an error if any fields are not as expected
func (c *Config) Validate() error {
	if c.Interval < 0 {
		return errors.New("invalid interval")
	} else if c.Interval == 0 {
		return errors.New("interval must be greater than 0")
	}
	return nil
}

// Collector implements a Collector to capture the run time of a process
type Collector interface {
	Name() string
	// Init should initialize the Collector returning an error on
	// any failure signaling the Collector not to be loaded
	Init(*Config) error
	// Updates returns a channel with runstats. This is used when a collector
	// wants to push data to the engine rather than waiting to be called
	// periodically.  This is useful when listening to events and pushing
	// run stats.
	Updates() <-chan RunStats

	// Collection should return runtimes of resources.  End is not
	// used as it will be filled in with the current time on each
	// invocation
	Collect(context.Context) ([]RunStats, error)

	//
	Stop()
}

// New returns a new Collector of the given type or an error if the type is
// not supported
func New(typ string) (cltr Collector, err error) {
	switch typ {
	case "node":
		cltr = &NodeCollector{}
		// conf[typ].Config["node"] = *nd
	case "docker":
		cltr = &DockerCollector{}
	default:
		err = fmt.Errorf("unsupported collector: %s", typ)
	}

	return
}

type collector struct {
	interval time.Duration
	bc       Collector
	out      chan<- []RunStats
	done     chan struct{}
	log      *zap.Logger
}

func (c *collector) run(ctx context.Context) {
	// Run initial collection
	data, err := c.bc.Collect(ctx)
	if err == nil {
		c.out <- data
	} else {
		c.log.Info("collection failed", zap.String("name", c.bc.Name()), zap.Error(err))
	}

	var (
		updates = c.bc.Updates()
		timer   = time.NewTimer(c.interval)
	)
	c.log.Info("collector starting", zap.String("name", c.bc.Name()), zap.Duration("interval", c.interval))
	if updates == nil {
		c.log.Info("collector push not supported", zap.String("name", c.bc.Name()))
	}

	for {
		select {
		case <-timer.C:
			data, err := c.bc.Collect(ctx)
			if err == nil {
				c.out <- data
			} else {
				c.log.Info("collection failed", zap.String("name", c.bc.Name()), zap.Error(err))
			}
			timer.Reset(c.interval)

		case evt := <-updates:
			c.out <- []RunStats{evt}

		case <-ctx.Done():
			c.log.Info("collector scheduler stopped", zap.String("name", c.bc.Name()))
			timer.Stop()
			close(c.done)
			return

		}
	}
}

func (c *collector) stop() {
	c.bc.Stop()
	<-c.done
}
