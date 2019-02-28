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

// Collector implements the collection of runtime information
type Collector interface {
	RunStats() <-chan []collector.RunStats
}

// Config holds the metermaid config
type Config struct {
	Node        *node.Node
	Collector   *collector.Engine
	Sink        sink.Sink
	DefaultMeta types.Meta
	Logger      *zap.Logger
}

type MeterMaid interface {
	Start() error
	Shutdown()
}

type metermaid struct {
	node node.Node

	defaultTags types.Meta

	cpuWeight float64
	memWeight float64

	pricer    *pricing.Pricer
	collector *collector.Engine
	sink      sink.Sink

	done chan struct{}
	log  *zap.Logger
}

func New(conf *Config) MeterMaid {
	mm := &metermaid{
		node:        *conf.Node,
		cpuWeight:   0.5,
		memWeight:   0.5,
		collector:   conf.Collector,
		sink:        conf.Sink,
		pricer:      pricing.NewPricer(*conf.Node, conf.Logger),
		defaultTags: conf.DefaultMeta,
		log:         conf.Logger,
		done:        make(chan struct{}),
	}

	// if err := mm.pricer.Initialize(); err != nil {
	// 	mm.log.Fatal("pricer failed to initialize", zap.Error(err))
	// }

	if mm.sink == nil {
		mm.sink = &sink.StdoutSink{}
	}
	mm.log.Info("sink loaded", zap.String("name", mm.sink.Name()))

	// go mm.run(conf.Collector.RunStats())

	return mm
}

func (pc *metermaid) Start() error {
	err := pc.pricer.Initialize()
	if err == nil {
		pc.collector.Start()
		go pc.run(pc.collector.RunStats())
	}
	return err
}

func (pc *metermaid) run(ch <-chan []collector.RunStats) {
	for runStats := range ch {
		seri := []tsdb.Series{}
		now := time.Now()
		for _, rs := range runStats {
			if rs.End.Unix() <= 0 {
				rs.End = now
			}
			for k, v := range pc.defaultTags {
				rs.Meta[k] = v
			}

			s, err := pc.makePriceSeries(rs)
			if err == nil {
				seri = append(seri, s)
				continue
			}
			pc.log.Info("failed to make series", zap.Error(err))
		}

		if len(seri) > 0 {
			if err := pc.sink.Publish(seri...); err != nil {
				pc.log.Info("failed to publish", zap.Error(err))
			} else {
				for _, s := range seri {
					pc.log.Debug("published", zap.String("name", s.ID()))
				}
				pc.log.Info("published", zap.Int("count", len(seri)))
			}
		}
	}

	close(pc.done)
}

func (pc *metermaid) makePriceSeries(rs collector.RunStats) (tsdb.Series, error) {
	s := tsdb.Series{Meta: rs.Meta}
	prices, err := pc.pricer.History(rs.Start, rs.End)
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
			cu = prices.Scale(pc.cpuWeight * float64(rs.CPU) / float64(pc.node.CPUShares))
		} else {
			cu = prices.Scale(pc.cpuWeight)
		}

		if rs.Memory > 0 {
			mu = prices.Scale(pc.memWeight * float64(rs.Memory) / float64(pc.node.Memory))
		} else {
			mu = prices.Scale(pc.memWeight)
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

func (pc *metermaid) Shutdown() {
	pc.collector.Stop()
	<-pc.done
}
