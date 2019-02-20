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

	pp        pricing.PriceProvider
	cpuWeight float64
	memWeight float64

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
			rCPU, rMem := mm.utilizationPercent(c)
			c.UnitsBurned, _ = mm.computePrice(c, rCPU, rMem)
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
	cpu = float64(c.CPUShares) / float64(mm.node.CPUShares)
	if cpu == 0 {
		// Full node price if no cpu set
		cpu = 1
	}
	mem = float64(c.Memory) / float64(mm.node.Memory)
	if mem == 0 {
		// Full node price if no mem set
		mem = 1
	}
	return
}

// computePrice computes the price of the container using the percent of the total
// price for the node
func (mm *meterMaid) computePrice(update types.Container, rCPU, rMem float64) (float64, error) {
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
		d     float64
	)

	for i, p := range prices[:l] {
		d = float64(time.Duration(prices[i+1].Timestamp - p.Timestamp))
		// Add each cpu and mem cost
		total += (((p.Price * mm.cpuWeight) / 3600e9) * d * rCPU) + (((p.Price * mm.memWeight) / 3600e9) * d * rMem)
	}

	p := prices[l]
	d = float64(time.Duration(end.UnixNano() - int64(p.Timestamp)))
	return total + (((p.Price * mm.cpuWeight) / 3600e9) * d * rCPU) + (((p.Price * mm.memWeight) / 3600e9) * d * rMem), nil
}
