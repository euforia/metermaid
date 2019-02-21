package tsdb

import (
	"sort"
	"time"
)

// DataPoints are a set of time sotable data points
type DataPoints []DataPoint

// Encompasses returns true if the DataPoints encompass the given start and stop
// timestamps
func (c DataPoints) Encompasses(start, end uint64) bool {
	l := len(c) - 1
	if l >= 0 {
		return start >= c[0].Timestamp && start < c[l].Timestamp &&
			end > c[0].Timestamp && end <= c[l].Timestamp
	}
	return false
}
func (c DataPoints) Sum() (total float64) {
	for _, i := range c {
		total += i.Value
	}
	return
}

// Get returns the set of prices from start to end timestamp in epoch nanoseconds
func (c DataPoints) Get(start, end uint64) DataPoints {
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
			ei = i
			break
		}
	}
	return c[si:ei]
}

// Add adds the given prices to the history, sorts and returns the new
// complete history
func (c DataPoints) Add(dps ...DataPoint) DataPoints {
	return append(c, dps...)
}

// Dedup returns a dedupped copy of the series
func (c DataPoints) Dedup() DataPoints {
	l := len(c)
	out := make(DataPoints, 0, l)
	out = append(out, c[0])

	for i := 1; i < l; i++ {
		if c[i].Equal(out[len(out)-1]) {
			continue
		}
		out = append(out, c[i])
	}
	return out
}

func (c DataPoints) Len() int {
	return len(c)
}

func (c DataPoints) Less(i, j int) bool {
	return c[i].Timestamp < c[j].Timestamp
}

func (c DataPoints) Swap(i, j int) {
	c[i], c[j] = c[j], c[i]
}

// Scale multiplies each value by the given multiplier returning a new
// set of datapoints
func (c DataPoints) Scale(multiplier float64) DataPoints {
	out := make(DataPoints, len(c))
	for i := range c {
		out[i] = DataPoint{
			Timestamp: c[i].Timestamp,
			Value:     c[i].Value * multiplier,
			Meta:      c[i].Meta,
		}
	}
	return out
}

// Per returns DataPoints that are filled in per the given interval
func (c DataPoints) Per(dur time.Duration) DataPoints {
	var (
		l    = len(c) - 1
		d    = uint64(dur)
		gend = make(DataPoints, 0)
	)

	for i, p := range c[:l] {
		delta := c[i+1].Timestamp - p.Timestamp
		slots := delta / d
		for i := uint64(1); i < slots; i++ {
			gend = gend.Add(DataPoint{
				Timestamp: p.Timestamp + (i * d),
				Value:     p.Value,
				Meta:      p.Meta,
			})
		}
	}

	if len(gend) > 0 {
		gend = gend.Add(c...)
		sort.Sort(gend)
		return gend
	}

	return c
}
