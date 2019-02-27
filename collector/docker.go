package collector

import (
	"context"
	"errors"
	"os"
	"time"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/euforia/metermaid/types"
)

const dockerAPIVersion = "1.37"

type DockerCollector struct {
	client *client.Client
	// meta keys to in the container label to add to the RunStats
	meta []string
}

func (dc *DockerCollector) Name() string {
	return "docker"
}

func (dc *DockerCollector) Init(conf map[string]interface{}) error {
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
		if taglist, ok := tags.([]string); ok {
			dc.meta = taglist
		} else {
			return errors.New("tags must be a list of strings")
		}
	}

	cli, err := client.NewEnvClient()
	if err == nil {
		dc.client = cli
	}
	return err
}

func (dc *DockerCollector) Collect(ctx context.Context) ([]RunStats, error) {
	opts := dtypes.ContainerListOptions{All: true}
	list, err := dc.client.ContainerList(ctx, opts)
	if err != nil {
		return nil, err
	}

	rts := make([]RunStats, len(list))
	for i, item := range list {
		rt := RunStats{
			Resource: ResourceContainer,
			Start:    time.Unix(item.Created, 0),
			Meta: types.Meta{
				"container": item.ID,
			},
		}
		cont, err := dc.client.ContainerInspect(ctx, item.ID)
		if err == nil {
			rt.CPU = uint64(cont.HostConfig.Memory)
			rt.Memory = uint64(cont.HostConfig.CPUShares)
		}

		for _, k := range dc.meta {
			if val, ok := item.Labels[k]; ok {
				rt.Meta[k] = val
			}
		}
		rts[i] = rt
	}
	return rts, nil
}
