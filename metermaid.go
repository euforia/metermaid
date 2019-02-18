package metermaid

import (
	"context"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"go.uber.org/zap"

	"github.com/euforia/metermaid/types"
)

// MeterMaid  ...
type MeterMaid interface {
	Updates() <-chan types.Container
	Stop() error
}

// ContainerProvider implements a container data provider
type ContainerProvider interface {
	// should return a list of known containers.
	Containers(context.Context) ([]*types.Container, error)
	// should return a contianer by the given id
	Container(ctx context.Context, id string) (*types.Container, error)
	// should clean up as needed
	Close() error
}

type meterMaid struct {
	// container info provider
	cp ContainerProvider

	// Containers currently running
	containers map[string]*types.Container

	// Outbound channel for container updates
	out chan types.Container

	cancel context.CancelFunc
	done   chan struct{}

	log *zap.Logger
}

// New returns a new MeterMaid interface
func New(logger *zap.Logger) (MeterMaid, error) {
	dockerClient, err := NewDockerClient("")
	if err != nil {
		return nil, err
	}

	mm := &meterMaid{
		cp:         dockerClient,
		containers: make(map[string]*types.Container),
		out:        make(chan types.Container, 32),
		done:       make(chan struct{}, 1),
		log:        logger,
	}

	if mm.log == nil {
		mm.log, _ = zap.NewDevelopment()
	}

	go mm.run()
	return mm, nil
}

func (mm *meterMaid) run() {
	ctx := context.Background()
	ctx, mm.cancel = context.WithCancel(ctx)

	mm.seedWithRunning(ctx)

	events, errs := mm.cp.(*DockerClient).Events(ctx, dtypes.EventsOptions{})
	mm.log.Info("listening for events")
	for {
		select {
		case event := <-events:
			mm.handleEvent(event)

		case err := <-errs:
			mm.log.Info("docker event error", zap.Error(err))

		case <-ctx.Done():
			mm.log.Info("event loop exiting")
			close(mm.out)
			close(mm.done)
			return

		}
	}
}

func (mm *meterMaid) Updates() <-chan types.Container {
	return mm.out
}

func (mm *meterMaid) handleEvent(event events.Message) {
	if event.Type != "container" {
		return
	}

	// Action and status will be equal if the action succeeded?? We skip over
	// failed actions
	if event.Action != event.Status {
		return
	}

	var cont *types.Container

	switch event.Action {
	case "create":
		var err error
		cont, err = mm.cp.Container(context.Background(), event.Actor.ID)
		if err == nil {
			mm.containers[event.Actor.ID] = cont
			mm.log.Debug("tracking", zap.String("id", event.Actor.ID), zap.String("action", "create"))
		} else {
			mm.log.Info("failed to get container details",
				zap.String("id", event.Actor.ID),
				zap.Error(err),
			)
			return
		}
	// case "attach":
	case "start":
		var ok bool
		if cont, ok = mm.containers[event.Actor.ID]; ok {
			cont.Start = event.TimeNano
		}
	// case "resize":
	case "die":
		var ok bool
		if cont, ok = mm.containers[event.Actor.ID]; ok {
			cont.Stop = event.TimeNano
			mm.log.Debug("container died",
				zap.String("id", cont.ID),
				zap.Duration("runtime", cont.RunTime()),
			)
		}
	case "destroy":
		var ok bool
		if cont, ok = mm.containers[event.Actor.ID]; ok {
			cont.Destroy = event.TimeNano
			// Once destroyed we stop tracking the container
			delete(mm.containers, cont.ID)
			mm.log.Debug("container destroyed",
				zap.String("id", cont.ID),
				zap.Duration("alloctime", cont.AllocatedTime()),
			)
		}
	default:
		return
	}

	if cont != nil {
		mm.out <- *cont
	}

}

//  seedWithRunning gets the list of running containers and populates
// the initial state.  This is meant to be called once on startup.
func (mm *meterMaid) seedWithRunning(ctx context.Context) {
	list, _ := mm.cp.Containers(ctx)
	mm.log.Info("seeding with running containers", zap.Int("count", len(list)))

	for _, cont := range list {
		mm.log.Info("tracking", zap.String("id", cont.ID), zap.String("action", "seed"))
		mm.containers[cont.ID] = cont
		mm.out <- *cont
	}
}

func (mm *meterMaid) Stop() error {
	mm.log.Info("stopping")
	// Stop main loop
	mm.cancel()
	// Wait for shutdown
	<-mm.done
	// Close docker connection
	err := mm.cp.Close()
	mm.log.Info("stopped")
	return err
}
