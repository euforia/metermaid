package tsdb

import (
	"fmt"

	"github.com/euforia/metermaid/types"
)

// Series is a grouping of data points by a unique name
// key-value pair metadata
type Series struct {
	Name string
	Meta types.Meta
	Data DataPoints
}

// ID returns the unique name and meta string
func (s Series) ID() string {
	return s.Name + "{" + s.Meta.String() + "}"
}

func (s Series) String() string {
	return fmt.Sprintf("%s{%s} %v", s.Name, s.Meta.String(), s.Data)
}

// Start returns the timestamp of the first point
func (s *Series) Start() uint64 {
	return s.Data[0].Timestamp
}

// End returns the timestamp of the last point
func (s *Series) End() uint64 {
	return s.Data.Last().Timestamp
}

// AddDataPoints adds the given points to the series
func (s *Series) AddDataPoints(dps ...DataPoint) {
	s.Data = s.Data.Insert(dps...)
}
