package tsdb

import "github.com/euforia/metermaid/fl"

// DataPoint holds a single data point in time along with any key value
// paired metadata
type DataPoint struct {
	// In nanoseconds
	Timestamp uint64
	//
	Value float64
	//
	Meta map[string]string
}

// Equal returns true if all fields are equal
func (dp *DataPoint) Equal(data DataPoint) bool {
	if dp.Value == data.Value && dp.Timestamp == data.Timestamp {
		for k, v := range dp.Meta {
			if data.Meta[k] != v {
				return false
			}
		}
		return true
	}
	return false
}

func (dp *DataPoint) Match(query fl.Query) bool {
	for k, q := range query {
		if !dp.MatchField(k, q...) {
			return false
		}
	}
	return true
}

func (dp *DataPoint) MatchField(name string, filters ...fl.Filter) bool {
	switch name {
	case "Timestamp":
		return dp.MatchTimestamp(filters...)
	case "Value":
		return dp.MatchValue(filters...)
	default:
		for _, filter := range filters {
			if dp.MatchMeta(name, filter) {
				continue
			}
			return false
		}
	}
	return true
}

func (dp *DataPoint) MatchTimestamp(filters ...fl.Filter) bool {
	for _, filter := range filters {
		if fl.MatchInt64(int64(dp.Timestamp), filter) {
			continue
		}
		return false
	}
	return true
}

func (dp *DataPoint) MatchValue(filters ...fl.Filter) bool {
	for _, filter := range filters {
		if fl.MatchFloat64(dp.Value, filter) {
			continue
		}
		return false
	}
	return true
}

func (dp *DataPoint) MatchMeta(key string, filter fl.Filter) bool {
	val, ok := dp.Meta[key]
	if !ok {
		return false
	}
	return fl.MatchString(val, filter)
}
