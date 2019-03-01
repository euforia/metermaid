package tsdb

import (
	"time"
)

// DataPoints are a set of time sotable data points
type DataPoints []DataPoint

// SumPerHour returns the sum of all values per hour.  The last value
// is taken to fill in the gaps
func (c DataPoints) SumPerHour() (total float64) {
	var (
		l = len(c) - 1
		d time.Duration
	)
	if l < 0 {
		return
	}

	for i, p := range c[:l] {
		d = time.Duration(c[i+1].Timestamp - p.Timestamp)
		// Add cost per hour times the number of hours
		total += c[i].Value * d.Hours()
	}
	return total
}

// Sum returns the sum of all values
func (c DataPoints) Sum() (total float64) {
	for _, i := range c {
		total += i.Value
	}
	return
}

// Max returns the max of the DataPoints
func (c DataPoints) Max() (max float64) {
	for _, p := range c {
		if p.Value > max {
			max = p.Value
		}
	}
	return
}

// Min returns the min of the DataPoints
func (c DataPoints) Min() (min float64) {
	l := len(c)
	switch l {
	case 0:
		return
	case 1:
		return c[0].Value
	}

	min = c[0].Value
	for _, p := range c[1:] {
		if p.Value < min {
			min = p.Value
		}
	}
	return
}

// Get returns the set of prices from start to end timestamp in epoch nanoseconds.
// It assumes the datapoints are sorted.  Unsorted data points will result in
// invalid data.
func (c DataPoints) Get(start, end uint64) DataPoints {
	var (
		si int
		ei = len(c) - 1 // set to last index
	)

	if ei == -1 || start > c[ei].Timestamp {
		return nil
	}

	if start > c[0].Timestamp {
		for i := 1; i <= ei; i++ {
			if start <= c[i].Timestamp {
				si = i
				break
			}
		}
	}

	if end < c[ei].Timestamp {
		for i := ei; i > -1; i-- {
			if end >= c[i].Timestamp {
				ei = i
				break
			}
		}
	}

	return c[si : ei+1]
}

// Clone returns a copy of all the data points
func (c DataPoints) Clone() DataPoints {
	clone := make(DataPoints, len(c))
	copy(clone, c)
	return clone
}

// Last returns the last data point in the set
func (c DataPoints) Last() DataPoint {
	return c[len(c)-1]
}

// Insert adds the given prices to the history, sorts and returns the new
// complete history
func (c DataPoints) Insert(dps ...DataPoint) DataPoints {
	return append(c, dps...)
}

// Dedup returns a dedupped copy of the series
func (c DataPoints) Dedup() DataPoints {
	l := len(c)
	if l <= 1 {
		return c
	}

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

// Scale multiplies each value by the given multiplier returning a new
// set of datapoints
func (c DataPoints) Scale(multiplier float64) DataPoints {
	out := make(DataPoints, len(c))
	for i := range c {
		out[i] = DataPoint{
			Timestamp: c[i].Timestamp,
			Value:     c[i].Value * multiplier,
			// Meta:      c[i].Meta,
		}
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
