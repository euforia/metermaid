package metermaid

import (
	"context"
	"os"
	"time"

	"github.com/docker/docker/client"

	"github.com/euforia/metermaid/types"
)

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

// ContainerStats returns stats for a container by the given id
func (client *DockerClient) ContainerStats(ctx context.Context, id string) (*types.Container, error) {
	details, err := client.ContainerInspect(ctx, id)
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
		startedAt, _ := time.Parse(time.RFC3339Nano, details.State.StartedAt)
		cont.Start = startedAt.UnixNano()
		if !details.State.Running {
			finishAt, _ := time.Parse(time.RFC3339Nano, details.State.FinishedAt)
			cont.Stop = finishAt.UnixNano()
		}
		return cont, nil
	}
	return nil, err
}
