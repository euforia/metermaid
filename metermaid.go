package metermaid

import (
	"errors"
	"time"

	"github.com/euforia/metermaid/node"
	"github.com/euforia/metermaid/pricing"
	"github.com/euforia/metermaid/storage"
	"github.com/euforia/metermaid/types"
	"go.uber.org/zap"
)

// Metermaid inmplements the canonical pricing interface
type Metermaid interface {
	PriceReport(start, end time.Time) (*pricing.Report, error)
	Containers() storage.Containers
}

type Config struct {
	Node             *node.Node
	ContainerStorage storage.Containers
	Pricer           pricing.Provider
	Collector        CCollector
	Logger           *zap.Logger
}

type meterMaid struct {
	node *node.Node

	pp *pricing.Pricer

	cpuWeight float64
	memWeight float64

	cstore storage.Containers
	log    *zap.Logger
}

// New returns a new Metermaid instance
func New(conf *Config) Metermaid {
	mm := &meterMaid{
		node:      conf.Node,
		cpuWeight: 0.5,
		memWeight: 0.5,
		pp:        pricing.NewPricer(conf.Pricer, *conf.Node, conf.Logger),
		cstore:    conf.ContainerStorage,
		log:       conf.Logger,
	}

	go mm.run(conf.Collector.Updates())

	return mm
}

func (mm *meterMaid) Containers() storage.Containers {
	return mm.cstore
}

func (mm *meterMaid) PriceReport(start, end time.Time) (*pricing.Report, error) {
	history, err := mm.pp.History(start, end)
	// history, err := mm.priceHistory(start, end)
	if err == nil {
		// per, _ := time.ParseDuration("1h")
		return pricing.NewReport(history), nil
	}
	return nil, err
}

func (mm *meterMaid) run(updates <-chan types.Container) {
	// This loop will exit once the collector closes the above channel
	// If select is used then the validity of the read must be checked.
	var err error
	for c := range updates {
		c.UnitsBurned, err = mm.computeContainerPrice(c)
		if err != nil {
			mm.log.Info("failed to compute price", zap.Error(err))
		}

		mm.cstore.Set(c)
		mm.log.Info("update",
			zap.String("id", c.ID),
			zap.Duration("runtime", c.RunTime()),
			zap.Duration("alloctime", c.AllocatedTime()),
			zap.Float64("burned", c.UnitsBurned),
		)
	}
}

// func (mm *meterMaid) priceHistory(start, end time.Time) (tsdb.DataPoints, error) {
// 	bt := int64(mm.node.BootTime)
// 	if start.UnixNano() < bt {
// 		start = time.Unix(0, bt)
// 	}
// 	d, err := mm.pp.History(start, end)
// 	return d, err
// }

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

// computeContainerPrice computes the price of the container using the percent of the total
// price for the node
func (mm *meterMaid) computeContainerPrice(update types.Container) (float64, error) {
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

	prices, err := mm.pp.History(start, end)
	if err != nil {
		return 0, err
	}

	if len(prices) > 0 {
		cprices := prices.Scale(mm.cpuWeight * rCPU)
		mprices := prices.Scale(mm.memWeight * rMem)
		return cprices.SumPerHour() + mprices.SumPerHour(), nil
	}

	return 0, errors.New("no price history")
}

// end defines how long the last price should be applied for
// func computePriceOverTime(prices tsdb.DataPoints, cpuWeight, memWeight float64) (cpuPrice, memPrice float64) {
// var (
// 	l = len(prices) - 1
// 	d time.Duration
// )

// for i, p := range prices[:l] {
// 	d = time.Duration(prices[i+1].Timestamp - p.Timestamp)
// 	// Add each cpu and mem cost
// 	cpuPrice += cprices[i].Value * d.Hours()
// 	memPrice += mprices[i].Value * d.Hours()
// }

// d = time.Duration(end.UnixNano() - int64(prices[l].Timestamp))
// cpuPrice += cprices[l].Value * d.Hours()
// memPrice += mprices[l].Value * d.Hours()
// return
// }

// // computeTotalPriceHours computes the total cost of the given timeline.
// // It assumes the pricing to be on a per hour basis
// func computeTotalPriceHours(prices tsdb.DataPoints) (total float64) {
// 	var (
// 		l = len(prices) - 1
// 		d time.Duration
// 	)
// 	for i, p := range prices[:l] {
// 		d = time.Duration(prices[i+1].Timestamp - p.Timestamp)
// 		// Add each cpu and mem cost
// 		total += prices[i].Value * d.Hours()
// 	}
// 	return total
// }
