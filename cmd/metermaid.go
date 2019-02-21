package main

import (
	"time"

	"github.com/euforia/metermaid/tsdb"

	"github.com/euforia/metermaid/node"
	"github.com/euforia/metermaid/pricing"
	"github.com/euforia/metermaid/storage"
	"github.com/euforia/metermaid/types"
	"go.uber.org/zap"
)

type meterMaid struct {
	node *node.Node

	pp *pricing.Pricer

	cpuWeight float64
	memWeight float64

	cstore storage.Containers
	log    *zap.Logger
}

func newMeterMaid(nd *node.Node, cstore storage.Containers, pricer *pricing.Pricer, logger *zap.Logger) *meterMaid {
	return &meterMaid{
		node:      nd,
		cpuWeight: 0.5,
		memWeight: 0.5,
		pp:        pricer,
		cstore:    cstore,
		log:       logger,
	}
}

func (mm *meterMaid) BurnHistory(start, end time.Time) (tsdb.DataPoints, error) {
	bt := int64(mm.node.BootTime)
	if start.UnixNano() < bt {
		start = time.Unix(0, bt)
	}
	return mm.pp.History(start, end, mm.node.Meta)
}

func (mm *meterMaid) run(updates <-chan types.Container) {
	for {
		select {
		case c := <-updates:
			c.UnitsBurned, _ = mm.computePrice(c)
			mm.cstore.Set(c)

			mm.log.Info("update",
				zap.String("id", c.ID[:12]),
				zap.Duration("runtime", c.RunTime()),
				zap.Duration("alloctime", c.AllocatedTime()),
				zap.Float64("burned", c.UnitsBurned),
			)
		}
	}
}

func (mm *meterMaid) utilizationPercent(c types.Container) (cpu float64, mem float64) {
	cpu = mm.node.CPUPercent(uint64(c.CPUShares))
	if cpu == 0 {
		// Full utilization if no cpu set
		cpu = 1
	}

	mem = mm.node.MemoryPercent(uint64(c.Memory))
	if mem == 0 {
		// Full utilization if no mem set
		mem = 1
	}

	return
}

// computePrice computes the price of the container using the percent of the total
// price for the node
func (mm *meterMaid) computePrice(update types.Container) (float64, error) {
	var (
		rCPU, rMem = mm.utilizationPercent(update)
		start      = time.Unix(0, update.Create)
		end        time.Time
	)

	if update.Destroy > 0 {
		end = time.Unix(0, update.Destroy)
	} else if update.Stop > 0 {
		end = time.Unix(0, update.Stop)
	} else {
		end = time.Now()
	}

	prices, err := mm.BurnHistory(start, end)
	if err == nil {
		cpuPrice, memPrice := computePriceOverTime(prices, end, mm.cpuWeight*rCPU, mm.memWeight*rMem)
		return cpuPrice + memPrice, nil
	}
	return 0, err
}

// end defines how long the last price should be applied for
func computePriceOverTime(prices tsdb.DataPoints, end time.Time, cpuWeight, memWeight float64) (cpuPrice, memPrice float64) {
	var (
		l = len(prices) - 1
		d time.Duration
	)
	// Prices is are per hour
	cprices := prices.Scale(cpuWeight)
	mprices := prices.Scale(memWeight)

	for i, p := range prices[:l] {
		d = time.Duration(prices[i+1].Timestamp - p.Timestamp)
		// Add each cpu and mem cost
		cpuPrice += cprices[i].Value * d.Hours()
		memPrice += mprices[i].Value * d.Hours()
	}

	d = time.Duration(end.UnixNano() - int64(prices[l].Timestamp))
	cpuPrice += cprices[l].Value * d.Hours()
	memPrice += mprices[l].Value * d.Hours()
	return
}
