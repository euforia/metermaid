package main

import (
	"time"

	"github.com/euforia/metermaid/node"
	"github.com/euforia/metermaid/pricing"
	"github.com/euforia/metermaid/storage"
	"github.com/euforia/metermaid/types"
	"go.uber.org/zap"
)

type meterMaid struct {
	node *node.Node

	pp pricing.PriceProvider

	cstore storage.Containers
	log    *zap.Logger
}

func (mm *meterMaid) BurnHistory(start, end time.Time) ([]*pricing.Price, error) {
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
			percent := mm.utilizationPercent(c)
			c.UnitsBurned, _ = mm.computePrice(c, percent)
			mm.cstore.Set(c)

			mm.log.Info("update",
				zap.String("id", c.ID[:12]),
				zap.Float64("ratio", percent),
				zap.Duration("runtime", c.RunTime()),
				zap.Duration("alloctime", c.AllocatedTime()),
				zap.Float64("burned", c.UnitsBurned),
			)
		}
	}
}

func (mm *meterMaid) utilizationPercent(c types.Container) float64 {
	pcpu := float64(c.CPUShares) / float64(mm.node.CPUShares)
	if pcpu == 0 {
		// Full node price if no cpu set
		pcpu = 1
	}
	pmem := float64(c.Memory) / float64(mm.node.Memory)
	if pmem == 0 {
		// Full node price if no mem set
		pmem = 1
	}
	if pcpu > pmem {
		return pcpu
	}
	return pmem
}

// computePrice computes the price of the container using the percent of the total
// price for the node
func (mm *meterMaid) computePrice(update types.Container, percent float64) (float64, error) {
	var (
		start = time.Unix(0, update.Create)
		end   time.Time
	)
	if update.Destroy > 0 {
		end = time.Unix(0, update.Destroy)
	} else if update.Stop > 0 {
		end = time.Unix(0, update.Stop)
	} else {
		end = time.Now()
	}

	prices, err := mm.BurnHistory(start, end)
	if err != nil {
		return 0, err
	}

	var (
		l     = len(prices) - 1
		total float64
		d     time.Duration
		np    float64
	)

	for i, p := range prices[:l] {
		np = p.Price / 3600e9
		d = time.Duration(prices[i+1].Timestamp - p.Timestamp)
		total += np * percent * float64(d)
	}

	np = prices[l].Price / 3600e9
	d = time.Duration(end.UnixNano() - int64(prices[l].Timestamp))
	return total + np*percent*float64(d), nil
}
