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
	"github.com/euforia/metermaid/sink"
	"github.com/euforia/metermaid/types"
)

var (
	nodeMeta   = flag.String("node.meta", "", "additional node metadata key=value, ...")
	metricMeta = flag.String("metric.meta", "", "default metadata to add to all collections key=value, ...")
	confFile   = flag.String("conf", "", "path to config file")
	debug      = flag.Bool("debug", false, "")
)

func init() {
	flag.Parse()
}

func makeNode(metastr string) *node.Node {
	nd := node.New()
	if metastr != "" {
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
		cltr, err := collector.New(typ)
		if err != nil {
			return err
		}

		c.Config["node"] = *nd
		err = cltr.Init(c.Config)
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

func makeSink(logger *zap.Logger, cont map[string]*config.SinkConfig) (sink.Sink, error) {
	msink := sink.NewMultiSink(logger)
	for k := range cont {
		snk, err := sink.New(k)
		if err == nil {
			msink.Register(snk)
			continue
		}
		return nil, err
	}

	if *debug {
		s, _ := sink.New("stdout")
		msink.Register(s)
	}

	return msink, nil
}

func getDefaultMeta() types.Meta {
	if *metricMeta != "" {
		return types.ParseMetaFromString(*metricMeta)
	}
	return nil
}

func main() {
	var (
		userConf  *config.Config
		err       error
		logger, _ = zap.NewDevelopment()
	)

	if *confFile != "" {
		userConf, err = config.ParseFile(*confFile)
		if err != nil {
			logger.Fatal("failed to parse config", zap.Error(err))
		}
	}

	nd := makeNode(*nodeMeta)
	logger.Info("node stats",
		zap.String("meta", nd.Meta.String()),
		zap.Uint64("cpu", nd.CPUShares),
		zap.Uint64("memory", nd.Memory),
		zap.Time("bootime", time.Unix(0, int64(nd.BootTime))),
	)

	eng := collector.NewEngine(logger)
	if err = makeCollectors(nd, eng, userConf.Collectors); err != nil {
		logger.Fatal("failed to initialize collectors", zap.Error(err))
	}

	conf := &metermaid.Config{
		Node:        nd,
		Collector:   eng,
		Logger:      logger,
		DefaultMeta: getDefaultMeta(),
	}

	// eng.Start()

	snk, err := makeSink(logger, userConf.Sinks)
	if err != nil {
		logger.Fatal("failed to initialize sink", zap.Error(err))
	}
	conf.Sink = snk

	mm := metermaid.New(conf)
	if err = mm.Start(); err != nil {
		logger.Fatal("failed to start", zap.Error(err))
	}

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	mm.Shutdown()
}
