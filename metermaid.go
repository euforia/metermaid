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

type CollectorEngine interface {
	RunStats() <-chan []collector.RunStats
}

const (
	metricNameCostNode      = "cost.node"
	metricNameCostContainer = "cost.container"
)

type metermaid struct {
	node node.Node

	pricer *pricing.Pricer

	cpuWeight float64
	memWeight float64

	sink sink.Sink

	defaultTags types.Meta

	log *zap.Logger
}

func NewMetermaid(nd node.Node, eng CollectorEngine, snk sink.Sink, defaulTags types.Meta, logger *zap.Logger) *metermaid {
	mm := &metermaid{
		node:        nd,
		cpuWeight:   0.5,
		memWeight:   0.5,
		sink:        snk,
		defaultTags: defaulTags,
		pricer:      pricing.NewPricer(nd, logger),
		log:         logger,
	}
	if mm.sink == nil {
		mm.sink = &sink.StdoutSink{}
	}

	mm.log.Info("sinks loaded", zap.String("name", mm.sink.Name()))
	go mm.run(eng.RunStats())

	return mm
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
