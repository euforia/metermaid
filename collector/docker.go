package collector

import (
	"context"
	"os"
	"time"

	dtypes "github.com/docker/docker/api/types"
	"github.com/docker/docker/client"
	"github.com/euforia/metermaid/types"
	"github.com/pkg/errors"
)

const dockerAPIVersion = "1.37"

type DockerCollector struct {
	client *client.Client
	// meta keys to in the container label to add to the RunStats
	meta []string
}

// Name satisfies the Collector interface
func (dc *DockerCollector) Name() string {
	return "docker"
}

// Init satisfies the Collector interface
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
		out, err := ifaceSliceToStringSlice(tags)
		if err != nil {
			return errors.Wrap(err, "labels")
		}
		dc.meta = out
	}

	cli, err := client.NewEnvClient()
	if err == nil {
		dc.client = cli
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
