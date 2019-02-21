package tsdb

type Series struct {
	Name string
	Data DataPoints
	Meta map[string]string
}
