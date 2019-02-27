package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/euforia/metermaid"
	"github.com/euforia/metermaid/collector"
	"github.com/euforia/metermaid/config"
	"github.com/euforia/metermaid/node"
	"github.com/euforia/metermaid/types"
)

var (
	nodeMeta   = flag.String("node.meta", "", "additional node metadata key=value, ...")
	metricMeta = flag.String("metric.meta", "", "default metadata to add to all collections key=value, ...")
	confFile   = flag.String("conf", "", "path to config file")
)

func init() {
	flag.Parse()
}

func makeNode() *node.Node {
	nd := node.New()
	if *nodeMeta != "" {
		meta := types.ParseMetaFromString(*nodeMeta)
		if nd.Meta == nil {
			nd.Meta = make(types.Meta)
		}
		for k, v := range meta {
			nd.Meta[k] = v
		}
	}
	return nd
}

func makeCollectors(nd *node.Node, eng *collector.Engine, conf map[string]*config.CollectorConfig) error {
	for typ, c := range conf {
		var cltr collector.Collector
		switch typ {
		case "node":
			cltr = &collector.NodeCollector{}
			conf[typ].Config["node"] = *nd

		case "docker":
			cltr = &collector.DockerCollector{}
		default:
			return fmt.Errorf("unsupported collector: %s", typ)
		}

		err := cltr.Init(c.Config)
		if err == nil {
			interval := c.IntervalDuration()
			if interval < 0 {
				return fmt.Errorf("invalid interval: %s", c.Interval)
			}
			eng.Register(cltr, interval)
			continue
		}
		return err
	}
	return nil
}

func main() {

	logger, _ := zap.NewDevelopment()
	var (
		conf *config.Config
		err  error
	)
	if *confFile != "" {
		conf, err = config.ParseFile(*confFile)
		if err != nil {
			logger.Fatal("failed to parse config", zap.Error(err))
		}
	}

	nd := makeNode()
	logger.Info("node stats",
		zap.String("meta", nd.Meta.String()),
		zap.Uint64("cpu", nd.CPUShares),
		zap.Uint64("memory", nd.Memory),
		zap.Time("bootime", time.Unix(0, int64(nd.BootTime))),
	)

	eng := collector.NewEngine(logger)
	if err = makeCollectors(nd, eng, conf.Collectors); err != nil {
		logger.Fatal("failed to initialize collectors", zap.Error(err))
	}

	// dc := &collector.DockerCollector{}
	// dc.Init(map[string]interface{}{"labels": []string{"service"}})
	// eng.Register(dc, 10*time.Second)

	// nc := &collector.NodeCollector{}
	// nc.Init(map[string]interface{}{"node": *nd})
	// eng.Register(nc, 10*time.Second)

	eng.Start()

	var defTags types.Meta
	if *metricMeta != "" {
		defTags = types.ParseMetaFromString(*metricMeta)
	}
	_ = metermaid.NewMetermaid(*nd, eng, nil, defTags, logger)

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	eng.Stop()
}
