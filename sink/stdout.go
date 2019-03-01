package sink

import (
	"fmt"
	"os"

	"github.com/euforia/metermaid/tsdb"
)

// StdoutSink implements a Sink interface that writes to stdout
type StdoutSink struct{}

// Name satisfies the Sink interface
func (sink *StdoutSink) Name() string {
	return "stdout"
}

// Publish satisfies the Sink interface
func (sink *StdoutSink) Publish(seri ...tsdb.Series) error {
	for _, s := range seri {
		fmt.Fprintf(os.Stdout, "%v\n", s)
	}
	return nil
}
