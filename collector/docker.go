package collector

import (
	"context"
	"os"
	"time"

	"go.uber.org/zap"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"github.com/docker/docker/client"
	"github.com/pkg/errors"

	"github.com/euforia/metermaid/types"
)

const dockerAPIVersion = "1.37"

// DockerCollector implements a Collector for container run times
type DockerCollector struct {
	client *client.Client
	// meta keys to in the container label to add to the RunStats
	meta []string

	ctx    context.Context
	cancel context.CancelFunc

	out  chan RunStats
	done chan struct{}

	log *zap.Logger
}

// Name satisfies the Collector interface
func (dc *DockerCollector) Name() string {
	return "docker"
}

// Init satisfies the Collector interface
func (dc *DockerCollector) Init(config *Config) error {

	dc.done = make(chan struct{})
	dc.out = make(chan RunStats, 16)
	dc.ctx, dc.cancel = context.WithCancel(context.Background())
	dc.log = config.Logger

	conf := config.Config
	if conf == nil {
		return nil
	}

	if str, ok := conf["DOCKER_API_VERSION"]; ok {
		if apiVersion, ok := str.(string); ok {
			os.Setenv("DOCKER_API_VERSION", apiVersion)
		} else {
			return errors.New("DOCKER_API_VERSION must be a string")
		}
	} else {
		os.Setenv("DOCKER_API_VERSION", dockerAPIVersion)
	}

	if tags, ok := conf["labels"]; ok {
		out, err := ifaceSliceToStringSlice(tags)
		if err != nil {
			return errors.Wrap(err, "labels")
		}
		dc.meta = out
	}

	cli, err := client.NewEnvClient()
	if err == nil {
		dc.client = cli
		go dc.events()
	}
	return err
}

// Collect satisfies the Collector interface
func (dc *DockerCollector) Collect(ctx context.Context) ([]RunStats, error) {
	opts := dtypes.ContainerListOptions{All: true}
	list, err := dc.client.ContainerList(ctx, opts)
	if err != nil {
		return nil, err
	}

	rts := make([]RunStats, len(list))
	for i, item := range list {
		rt, err := dc.makeRunStat(ctx, item.ID)
		if err == nil {
			rts[i] = rt
			continue
		}
		dc.log.Info("skipping stat", zap.String("container", item.ID), zap.Error(err))
	}
	return rts, nil
}

func (dc *DockerCollector) events() {
	events, errs := dc.client.Events(dc.ctx, dtypes.EventsOptions{})
	for {
		select {
		case event := <-events:
			dc.handleEvent(event)

		case err := <-errs:
			dc.log.Info("event loop exiting", zap.Error(err))
			goto EXIT_LOOP

		case <-dc.ctx.Done():
			dc.log.Info("event loop exiting")
			goto EXIT_LOOP

		}
	}

EXIT_LOOP:
	close(dc.out)
	close(dc.done)
	return
}

func (dc *DockerCollector) handleEvent(event events.Message) {
	if event.Type != "container" {
		return
	}
	// Action and status will be equal if the action succeeded?? We skip over
	// failed actions
	if event.Action != event.Status {
		return
	}

	var (
		rs  RunStats
		err error
	)
	switch event.Action {
	case "create":
		rs, err = dc.makeRunStat(dc.ctx, event.Actor.ID)

	case "die":
		rs, err = dc.makeRunStat(dc.ctx, event.Actor.ID)
		rs.End = time.Now()

	// case "destroy":
	// TODO:

	default:
		return
	}

	if err == nil {
		dc.out <- rs
		return
	}

	dc.log.Info("failed to make stat", zap.Error(err), zap.String("container", event.Actor.ID))
}

// makeRunStat gets the necessary container info based on id.
func (dc *DockerCollector) makeRunStat(ctx context.Context, cid string) (RunStats, error) {
	rt := RunStats{
		Resource: ResourceContainer,
		Meta:     types.Meta{"container": cid},
	}

	cont, err := dc.client.ContainerInspect(ctx, cid)
	if err != nil {
		return rt, err
	}

	rt.CPU = uint64(cont.HostConfig.Memory)
	rt.Memory = uint64(cont.HostConfig.CPUShares)

	labels := cont.Config.Labels
	// Add specified labels
	for _, k := range dc.meta {
		if val, ok := labels[k]; ok {
			rt.Meta[k] = val
		}
	}

	rt.Start, err = time.Parse(time.RFC3339Nano, cont.Created)
	return rt, err
}

func (dc *DockerCollector) Updates() <-chan RunStats {
	return dc.out
}

// Stop stops the event handler loop
func (dc *DockerCollector) Stop() {
	dc.cancel()
	<-dc.done
	dc.client.Close()
}

func ifaceSliceToStringSlice(tags interface{}) ([]string, error) {
	switch typ := tags.(type) {
	case []string:
		return typ, nil

	case []interface{}:
		out := make([]string, 0, len(typ))
		for _, k := range typ {
			if l, ok := k.(string); ok {
				out = append(out, l)
				continue
			}
			return nil, errors.New("must be a list of strings")
		}
		return out, nil

	}

	return nil, errors.New("must be a list of strings")
}
