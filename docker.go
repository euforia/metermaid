package metermaid

import (
	"context"
	"errors"
	"os"
	"time"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"

	"github.com/euforia/metermaid/types"
)

// ErrEventsNotSupported is returned when a provider does not support
// events
var ErrEventsNotSupported = errors.New("events not supported")

// DockerClient ...
type DockerClient struct {
	*client.Client
}

// NewDockerClient returns a new instance of DockerClient using defaults
// if necessary
func NewDockerClient(apiVersion string) (*DockerClient, error) {
	if apiVersion == "" {
		os.Setenv("DOCKER_API_VERSION", "1.37")
	} else {
		os.Setenv("DOCKER_API_VERSION", apiVersion)
	}

	cli, err := client.NewEnvClient()
	if err == nil {
		return &DockerClient{cli}, nil
	}
	return nil, err
}

// Container returns stats for a container by the given id
func (client *DockerClient) Container(ctx context.Context, id string) (*types.Container, error) {
	details, err := client.Client.ContainerInspect(ctx, id)
	if err == nil {
		cont := &types.Container{
			ID:        details.ID,
			Name:      details.Name,
			Labels:    details.Config.Labels,
			CPUShares: details.HostConfig.CPUShares,
			Memory:    details.HostConfig.Memory,
		}

		createdAt, _ := time.Parse(time.RFC3339Nano, details.Created)
		cont.Create = createdAt.UnixNano()

		if startedAt, err := time.Parse(time.RFC3339Nano, details.State.StartedAt); err == nil {
			cont.Start = startedAt.UnixNano()
		}

		if !details.State.Running {
			if finishAt, err := time.Parse(time.RFC3339Nano, details.State.FinishedAt); err == nil {
				cont.Stop = finishAt.UnixNano()
			}
		}

		return cont, nil
	}
	return nil, err
}

// Containers returns a list of running containers.  A complete or partial list
// is returned depending on the error. The error returned is the last error occurred
func (client *DockerClient) Containers(ctx context.Context) ([]*types.Container, error) {
	opts := dtypes.ContainerListOptions{All: true}
	list, err := client.Client.ContainerList(ctx, opts)
	if err != nil {
		return nil, err
	}

	containers := make([]*types.Container, 0, len(list))
	for i := range list {
		cont, er := client.Container(ctx, list[i].ID)
		if er == nil {
			containers = append(containers, cont)
		} else {
			err = er
		}
	}
	return containers, err
}
