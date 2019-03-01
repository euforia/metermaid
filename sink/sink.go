package sink

import (
	"fmt"

	"github.com/euforia/metermaid/tsdb"
)

// Sink is an interface that implements writing data to a backend
type Sink interface {
	Name() string
	// Should publish the series data
	Publish(...tsdb.Series) error
}

// New returns a new Sink of the given type or error if the type is not
// supported
func New(typ string) (Sink, error) {
	switch typ {
	case "datadog":
		return NewDataDogSink("", ""), nil
	case "stdout":
		return &StdoutSink{}, nil
	}

	return nil, fmt.Errorf("unsupported sink: %s", typ)
}
