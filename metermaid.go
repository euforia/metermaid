package metermaid

import (
	"fmt"
	"time"

	"github.com/euforia/metermaid/collector"
	"github.com/euforia/metermaid/node"
	"github.com/euforia/metermaid/pricing"
	"github.com/euforia/metermaid/sink"
	"github.com/euforia/metermaid/tsdb"
	"github.com/euforia/metermaid/types"
	"go.uber.org/zap"
)

const (
	metricNameCostNode      = "cost.node"
	metricNameCostContainer = "cost.container"
)

// Config holds the metermaid config
type Config struct {
	Node *node.Node
	// Meta tags to add to all metrics
	DefaultMeta types.Meta
	Logger      *zap.Logger
}

// MeterMaid is the canonical interface for price metering
type MeterMaid interface {
	RegisterCollector(string, *collector.Config) error
	RegisterSink(string) error
	Start() error
	Shutdown()
}

type metermaid struct {
	node *node.Node

	defaultTags types.Meta

	cpuWeight float64
	memWeight float64

	pricer    *pricing.Pricer
	collector *collector.Engine
	sink      *sink.MultiSink

	done chan struct{}
	log  *zap.Logger
}

// New returns a new MeterMaid instance based on the config
func New(conf *Config) MeterMaid {
	mm := &metermaid{
		node:        conf.Node,
		cpuWeight:   0.5,
		memWeight:   0.5,
		defaultTags: conf.DefaultMeta,
		done:        make(chan struct{}),
		log:         conf.Logger,
	}

	mm.init()

	return mm
}

func (mm *metermaid) init() {
	mm.log.Info("node",
		zap.Uint64("cpu", mm.node.CPUShares),
		zap.Uint64("memory", mm.node.Memory),
		zap.Time("bootime", time.Unix(0, int64(mm.node.BootTime))),
		zap.String("meta", mm.node.Meta.String()),
	)

	mm.collector = collector.NewEngine(mm.log)
	mm.sink = sink.NewMultiSink(mm.log)
	mm.pricer = pricing.NewPricer(*mm.node, mm.log)

	mm.log.Info("pricer loaded", zap.String("backend", mm.pricer.Name()))
}

func (mm *metermaid) RegisterSink(typ string) error {
	snk, err := sink.New(typ)
	if err == nil {
		mm.sink.Register(snk)
		mm.log.Info("sink registered", zap.String("name", snk.Name()))
	}
	return err
}

func (mm *metermaid) RegisterCollector(typ string, conf *collector.Config) error {
	cltr, err := collector.New(typ)
	if err == nil {
		conf.Node = mm.node
		err = mm.collector.Register(cltr, conf)
	}
	return err
}

func (mm *metermaid) Start() error {
	err := mm.pricer.Initialize()
	if err == nil {
		mm.collector.Start()
		mm.log.Info("sinks loaded", zap.String("name", mm.sink.Name()))
		go mm.run(mm.collector.RunStats())
	}
	return err
}

func (mm *metermaid) run(ch <-chan []collector.RunStats) {
	for runStats := range ch {
		seri := []tsdb.Series{}
		now := time.Now()
		for _, rs := range runStats {
			if rs.End.Unix() <= 0 {
				rs.End = now
			}
			for k, v := range mm.defaultTags {
				rs.Meta[k] = v
			}

			s, err := mm.makePriceSeries(rs)
			if err == nil {
				seri = append(seri, s)
				continue
			}
			mm.log.Info("failed to make series", zap.Error(err))
		}

		if len(seri) > 0 {
			if err := mm.sink.Publish(seri...); err != nil {
				mm.log.Info("failed to publish", zap.Error(err))
			} else {
				for _, s := range seri {
					mm.log.Debug("published", zap.String("name", s.ID()))
				}
				mm.log.Info("published", zap.Int("count", len(seri)))
			}
		}
	}

	close(mm.done)
}

func (mm *metermaid) makePriceSeries(rs collector.RunStats) (tsdb.Series, error) {
	s := tsdb.Series{Meta: rs.Meta}
	prices, err := mm.pricer.History(rs.Start, rs.End)
	if err != nil {
		return s, err
	}

	switch rs.Resource {
	case collector.ResourceNode:
		s.Name = metricNameCostNode
		s.Data = tsdb.DataPoints{tsdb.DataPoint{
			Timestamp: uint64(time.Now().UnixNano()),
			Value:     prices.SumPerHour(),
		}}

	case collector.ResourceContainer:
		s.Name = metricNameCostContainer

		var (
			cu tsdb.DataPoints
			mu tsdb.DataPoints
		)

		if rs.CPU > 0 {
			cu = prices.Scale(mm.cpuWeight * float64(rs.CPU) / float64(mm.node.CPUShares))
		} else {
			cu = prices.Scale(mm.cpuWeight)
		}

		if rs.Memory > 0 {
			mu = prices.Scale(mm.memWeight * float64(rs.Memory) / float64(mm.node.Memory))
		} else {
			mu = prices.Scale(mm.memWeight)
		}

		s.Data = tsdb.DataPoints{tsdb.DataPoint{
			Timestamp: uint64(time.Now().UnixNano()),
			Value:     cu.SumPerHour() + mu.SumPerHour(),
		}}

	default:
		return s, fmt.Errorf("unknown resource: %s", rs.Resource)

	}

	return s, nil
}

func (mm *metermaid) Shutdown() {
	mm.collector.Stop()
	<-mm.done
}
