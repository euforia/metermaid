package tsdb

// DataPoint holds a single data point in time along with any key value
// paired metadata
type DataPoint struct {
	// In nanoseconds
	Timestamp uint64
	//
	Value float64
}

// Equal returns true if all fields are equal
func (dp *DataPoint) Equal(data DataPoint) bool {
	return dp.Value == data.Value && dp.Timestamp == data.Timestamp
}
