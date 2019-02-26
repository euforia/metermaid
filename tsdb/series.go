package tsdb

import "github.com/euforia/metermaid/types"

// Series is a grouping of data points by a unique name
// key-value pair metadata
type Series struct {
	Name string
	Meta types.Meta
	Data DataPoints
}

func (s *Series) Start() uint64 {
	return s.Data[0].Timestamp
}

func (s *Series) End() uint64 {
	return s.Data.Last().Timestamp
}
