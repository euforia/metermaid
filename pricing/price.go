package pricing

import (
	"time"
)

// Price holds the price of on object at the given time.
type Price struct {
	Price float64
	// In nanoseconds
	Timestamp uint64
	// Things like AZ, InstanceType, product etc. to filter on
	Meta map[string]string
}

type sortedPrices []*Price

func (p sortedPrices) Len() int {
	return len(p)
}

func (p sortedPrices) Less(i, j int) bool {
	return p[i].Timestamp < p[j].Timestamp
}

func (p sortedPrices) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

// PriceProvider implments an interface to return pricing information
type PriceProvider interface {
	History(start, end time.Time, filter map[string]string) ([]*Price, error)
}

type Cache []*Price

func (c Cache) history(start, end uint64) []*Price {
	var (
		si int
		ei int
	)

	for i, price := range c {
		if price.Timestamp > start {
			si = i - 1
			break
		}
	}

	l := len(c) - 1
	for i := l; i > -1; i-- {
		if c[i].Timestamp < end {
			ei = i + 1
			break
		}
	}
	return c[si:ei]
}
