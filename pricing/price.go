package pricing

import "time"

// Price holds the price of on object at the given time.
type Price struct {
	Price float64
	// In nanoseconds
	Timestamp uint64
	// Things like AZ, InstanceType, product etc. to filter on
	Meta map[string]string
}

// Pricer implments an interface to return pricing information
type Pricer interface {
	History(start, end time.Time, filter map[string]string) ([]*Price, error)
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
