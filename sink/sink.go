package sink

import (
	"fmt"
	"os"

	"github.com/euforia/metermaid/tsdb"
)

type Sink interface {
	Publish(...tsdb.Series) error
}

type StdoutSink struct{}

func (sink *StdoutSink) Publish(seri ...tsdb.Series) error {
	for _, s := range seri {
		fmt.Fprintf(os.Stdout, "%v\n", s)
	}
	return nil
}
