package sink

import (
	"github.com/euforia/metermaid/tsdb"
	"github.com/pkg/errors"
	"go.uber.org/zap"
)

// MultiSink implements a wrapper Sink interface to multiple backend
// Sinks
type MultiSink struct {
	sinks []Sink
	log   *zap.Logger
}

// NewMultiSink returns a new MultiSink instance
func NewMultiSink(logger *zap.Logger, sinks ...Sink) *MultiSink {
	return &MultiSink{
		sinks: sinks,
		log:   logger,
	}
}

// Name satisfies the Sink interface
func (sink *MultiSink) Name() string {
	str := "multi["
	for _, s := range sink.sinks {
		str += s.Name() + ","
	}
	return str[:len(str)-1] + "]"
}

// Register registers a new backend sink. It is not thread-safe
func (sink *MultiSink) Register(s Sink) {
	sink.sinks = append(sink.sinks, s)
}

// Publish satisfies the Sink interface
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
