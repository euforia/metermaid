package sink

import (
	"fmt"
	"os"

	"go.uber.org/zap"

	"github.com/euforia/metermaid/tsdb"
	"github.com/pkg/errors"
)

type Sink interface {
	Name() string
	Publish(...tsdb.Series) error
}

type StdoutSink struct{}

func (sink *StdoutSink) Name() string {
	return "stdout"
}

func (sink *StdoutSink) Publish(seri ...tsdb.Series) error {
	for _, s := range seri {
		fmt.Fprintf(os.Stdout, "%v\n", s)
	}
	return nil
}

func New(typ string) (Sink, error) {
	switch typ {
	case "datadog":
		return NewDataDogSink("", ""), nil
	case "stdout":
		return &StdoutSink{}, nil
	}

	return nil, fmt.Errorf("unsupported sink: %s", typ)
}

type MultiSink struct {
	sinks []Sink
	log   *zap.Logger
}

func NewMultiSink(logger *zap.Logger, sinks ...Sink) *MultiSink {
	return &MultiSink{
		sinks: sinks,
		log:   logger,
	}
}

func (sink *MultiSink) Name() string {
	str := "multi["
	for _, s := range sink.sinks {
		str += s.Name() + ","
	}
	return str[:len(str)-1] + "]"
}

func (sink *MultiSink) Register(s Sink) {
	sink.sinks = append(sink.sinks, s)
}

func (sink *MultiSink) Publish(seri ...tsdb.Series) error {
	errs := make([]error, 0, len(sink.sinks))
	for _, s := range sink.sinks {
		err := s.Publish(seri...)
		if err != nil {
			errs = append(errs, errors.Wrap(err, s.Name()))
		} else {
			sink.log.Debug("published", zap.String("sink", s.Name()), zap.Int("count", len(seri)))
		}
	}
	if len(errs) > 0 {
		var e string
		for _, err := range errs {
			e += err.Error() + "\n"
		}
		return errors.New(e)
	}
	return nil
}
