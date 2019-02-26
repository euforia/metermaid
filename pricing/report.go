package pricing

import "github.com/euforia/metermaid/tsdb"

// Report is a price report
type Report struct {
	Total   float64
	Min     float64
	Max     float64
	Average float64

	History tsdb.DataPoints
}

// NewReport returns a new Price computing the per interval and total
func NewReport(data tsdb.DataPoints) *Report {
	out := &Report{History: data}
	out.Total = data.SumPerHour()
	out.Min = data.Min()
	out.Max = data.Max()
	out.Average = data.Sum() / float64(len(data))

	return out
}
