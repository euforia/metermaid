package main

import (
	"flag"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/euforia/metermaid"
	"github.com/euforia/metermaid/collector"
	"github.com/euforia/metermaid/node"
	"github.com/euforia/metermaid/types"
)

var (
	// bindAddr = flag.String("bind-addr", "127.0.0.1:8080", "")
	// advAddr  = flag.String("adv-addr", "", "")
	nodeMeta   = flag.String("node.meta", "", "node metadata key=value, ...")
	metricMeta = flag.String("metric.meta", "", "default metadata to add to all collections key=value, ...")
	// joinPeer = flag.String("join", "", "")
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

func main() {
	logger, _ := zap.NewDevelopment()
	nd := makeNode()
	logger.Info("node stats",
		zap.Uint64("cpu", nd.CPUShares),
		zap.Uint64("memory", nd.Memory),
		zap.Time("bootime", time.Unix(0, int64(nd.BootTime))),
	)

	eng := collector.NewEngine(logger)

	dc := &collector.DockerCollector{}
	dc.Init(map[string]interface{}{"labels": []string{"service"}})
	eng.Register(dc, 10*time.Second)

	nc := &collector.NodeCollector{}
	nc.Init(map[string]interface{}{"node": *nd})
	eng.Register(nc, 10*time.Second)

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
