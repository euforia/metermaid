package metermaid

import (
	"context"
	"fmt"

	"github.com/docker/docker/api/types"
	"github.com/docker/docker/api/types/events"
	"go.uber.org/zap"
)

// MeterMaid  ...
type MeterMaid interface {
	Stop() error
}

type meterMaid struct {
	docker *DockerClient

	containers map[string]*Container

	cancel context.CancelFunc
	done   chan struct{}

	log *zap.Logger
}

// New returns a new MeterMaid interface
func New() (MeterMaid, error) {
	client, err := NewDockerClient("")
	if err != nil {
		return nil, err
	}

	mm := &meterMaid{
		docker:     client,
		containers: make(map[string]*Container),
		log:        zap.NewExample(),
		done:       make(chan struct{}, 1),
	}
	go mm.run()
	return mm, nil
}

func (mm *meterMaid) run() {
	ctx := context.Background()
	ctx, mm.cancel = context.WithCancel(ctx)

	mm.inspectRunningContainers(ctx)

	events, errs := mm.docker.Events(ctx, types.EventsOptions{})
	mm.log.Info("listening for events")
	for {
		select {
		case event := <-events:
			mm.handleEvent(event)

		case err := <-errs:
			mm.log.Info("docker event error", zap.Error(err))

		case <-ctx.Done():
			close(mm.done)
			return

		}
	}
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

	switch event.Action {
	case "create":
		cont, err := mm.docker.ContainerStats(context.Background(), event.Actor.ID)
		if err == nil {
			mm.containers[event.Actor.ID] = cont
			mm.log.Info("tracking", zap.String("id", event.Actor.ID))
			return
		}
		mm.log.Info("failed to get container details",
			zap.String("id", event.Actor.ID),
			zap.Error(err),
		)

	// case "attach":
	case "start":
		if cont, ok := mm.containers[event.Actor.ID]; ok {
			cont.Start = event.TimeNano
		}
	// case "resize":
	case "die":
		if cont, ok := mm.containers[event.Actor.ID]; ok {
			cont.Stop = event.TimeNano
			mm.log.Info("container",
				zap.String("id", cont.ID),
				zap.Duration("runtime", cont.RunTime()),
			)
		}
	case "destroy":
		if cont, ok := mm.containers[event.Actor.ID]; ok {
			cont.Destroy = event.TimeNano
			mm.log.Info("container",
				zap.String("id", cont.ID),
				zap.Duration("alloctime", cont.AllocatedTime()),
			)
		}
		// TODO
		// ensure we are tracking this container
		// if died then do nothing.
	default:
		return
	}

	// log.Printf("---> %+v", event)
}

func (mm *meterMaid) inspectRunningContainers(ctx context.Context) {
	opts := types.ContainerListOptions{All: true}
	list, _ := mm.docker.ContainerList(ctx, opts)

	for _, cont := range list {
		scont, err := mm.docker.ContainerStats(ctx, cont.ID)
		if err == nil {
			mm.log.Info("tracking", zap.String("id", cont.ID))
			fmt.Printf("%+v\n", scont)
			mm.containers[scont.ID] = scont
		}
	}
}

func (mm *meterMaid) Stop() error {
	mm.log.Info("stopping")
	mm.cancel()
	// Wait for shutdown
	<-mm.done
	// Close docker connection
	err := mm.docker.Close()
	mm.log.Info("stopped")
	return err
}
