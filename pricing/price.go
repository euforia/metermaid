package pricing

import (
	"sort"
	"time"

	"github.com/euforia/metermaid/tsdb"
)

// PriceHistory ...
type PriceHistory struct {
	Total   float64
	History tsdb.DataPoints
}

// NewPriceHistory returns a new Price computing the per interval and total
func NewPriceHistory(data tsdb.DataPoints, per time.Duration) *PriceHistory {
	out := &PriceHistory{History: data}
	list := data.Per(per)
	out.Total = list.Sum()
	return out
}

// Provider implments an interface to return pricing information
type Provider interface {
	History(start, end time.Time, filter map[string]string) (tsdb.DataPoints, error)
}

// Pricer is a the canonical interface to interact with pricing data
// for the node
type Pricer struct {
	pp     Provider
	prices tsdb.DataPoints
}

func NewPricer(provider Provider) *Pricer {
	return &Pricer{pp: provider, prices: make(tsdb.DataPoints, 0)}
}

func (pr *Pricer) History(start, end time.Time, meta map[string]string) (tsdb.DataPoints, error) {
	s := uint64(start.UnixNano())
	e := uint64(end.UnixNano())
	if !pr.prices.Encompasses(s, e) {
		h, err := pr.pp.History(start, end, meta)
		if err != nil {
			return nil, err
		}
		prices := pr.prices.Add(h...)
		sort.Sort(prices)
		pr.prices = prices.Dedup()
	}
	return pr.prices.Get(s, e), nil
}
