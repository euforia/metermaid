package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"go.uber.org/zap"

	"github.com/euforia/metermaid"
	"github.com/euforia/metermaid/collector"
	"github.com/euforia/metermaid/node"
	"github.com/euforia/metermaid/types"
)

var (
	nodeMeta   = flag.String("node-meta", "", "additional node metadata key=value, ...")
	metricMeta = flag.String("metric-meta", "", "metadata to add to all collections key=value, ...")
	confFile   = flag.String("conf", "config.hcl", "path to config file")
	debug      = flag.Bool("debug", false, "")
)

func init() {
	flag.Parse()
}

func registerCollectors(mm metermaid.MeterMaid, conf map[string]*collectorConfig) (err error) {
	for k, v := range conf {
		err = mm.RegisterCollector(k, &collector.Config{
			Config:   v.Config,
			Interval: v.IntervalDuration(),
		})
		if err == nil {
			continue
		}
		return err
	}
	return nil
}

func registerSinks(mm metermaid.MeterMaid, conf map[string]*sinkConfig) (err error) {
	for k := range conf {
		if err = mm.RegisterSink(k); err == nil {
			continue
		}
		return
	}

	if *debug {
		err = mm.RegisterSink("stdout")
	}
	return
}

func main() {
	userConf, err := parseConfigFile()
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}

	logger, _ := zap.NewDevelopment()
	nd, err := node.NewWithMetaString(*nodeMeta)
	if err != nil {
		logger.Info("node metadata partially loaded", zap.Error(err))
	}

	conf := &metermaid.Config{
		Node:        nd,
		DefaultMeta: types.ParseMetaFromString(*metricMeta),
		Logger:      logger,
	}
	mm := metermaid.New(conf)

	// Register collectors
	if err = registerCollectors(mm, userConf.Collectors); err != nil {
		logger.Fatal("initialize collector", zap.Error(err))
	}

	// Register sinks
	if err = registerSinks(mm, userConf.Sinks); err != nil {
		logger.Fatal("initialize sink", zap.Error(err))
	}

	// Start
	if err = mm.Start(); err != nil {
		logger.Fatal("start", zap.Error(err))
	}

	// Setup signal handler
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	<-sigs
	mm.Shutdown()
}
