package pricing

import (
	"sort"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/euforia/metermaid/node"
	"github.com/euforia/metermaid/tsdb"
)

// Provider implments an interface to return pricing information
type Provider interface {
	Name() string
	History(start, end time.Time, filter map[string]string) (tsdb.DataPoints, error)
}

// Pricer is a the canonical interface to interact with pricing data
// for the node. It implements caching on top of the Provider
type Pricer struct {
	node node.Node
	pp   Provider

	mu          sync.RWMutex
	cache       tsdb.DataPoints
	lastFetched uint64

	log *zap.Logger
}

// NewPricer returns a new Pricer backed by the given provider
func NewPricer(provider Provider, nd node.Node, logger *zap.Logger) *Pricer {
	pr := &Pricer{
		pp:   provider,
		node: nd,
		log:  logger,
	}
	start := time.Unix(0, int64(nd.BootTime))
	pr.fetchHistory(start, start, time.Now())
	logger.Info("pricer",
		zap.String("backend", pr.pp.Name()),
		zap.Time("cache.start", time.Unix(0, int64(pr.cache[0].Timestamp))),
		zap.Int("cache.size", len(pr.cache)),
	)
	return pr
}

// History satisfies the Provider interface
func (pr *Pricer) History(start, end time.Time) (tsdb.DataPoints, error) {
	pr.log.Debug("price history request", zap.Time("start", start), zap.Time("end", end))
	// Set to boot time if start is less than that
	bt := int64(pr.node.BootTime)
	if start.UnixNano() < bt {
		start = time.Unix(0, bt)
	}

	prices, err := pr.history(start, end)
	if err != nil {
		return nil, err
	}
	pr.log.Debug("price history request",
		zap.Time("start", start), zap.Time("end", end), zap.Int("count", len(prices)))

	if len(prices) > 0 {
		// Add end marker for proper price calculation
		e := uint64(end.UnixNano())
		if last := prices.Last(); last.Timestamp < e {
			prices = prices.Insert(tsdb.DataPoint{
				Timestamp: e, Value: last.Value,
			})
		}
	}
	return prices, nil
}

func (pr *Pricer) history(start, end time.Time) (tsdb.DataPoints, error) {
	s := uint64(start.UnixNano())
	e := uint64(end.UnixNano())
	// We only check end as we always should have all data since
	// the boot time.
	pr.mu.RLock()
	// 5 min since last fetch
	if e <= pr.lastFetched+300e9 {
		prices := pr.cache.Get(s, e)
		pr.mu.RUnlock()
		return prices, nil
	}

	last := pr.cache.Last()
	pr.mu.RUnlock()
	return pr.fetchHistory(start, time.Unix(0, int64(last.Timestamp)), end)
}

// reqStart is the request start time. start is the start of the fetch. reqStart is used
// to return the query response to avoid an addtional lock/unlock cycle
func (pr *Pricer) fetchHistory(reqStart, start, end time.Time) (tsdb.DataPoints, error) {

	prices, err := pr.pp.History(start, end, pr.node.Meta)
	if err == nil {
		pr.log.Debug("fetched price history",
			zap.Time("start", start), zap.Time("end", end),
			zap.Int("count", len(prices)))

		pr.mu.Lock()
		prices = pr.cache.Insert(prices...)
		sort.Sort(prices)
		pr.cache = prices.Dedup()

		// pr.log.Debug("new price history",
		// 	zap.Int("count", len(pr.cache)))

		pr.lastFetched = uint64(time.Now().UnixNano())
		prices = pr.cache.Get(uint64(reqStart.UnixNano()), uint64(end.UnixNano()))
		pr.mu.Unlock()
	}
	// pr.log.Debug("price history result", zap.Int("size", len(prices)))
	return prices, err
}
