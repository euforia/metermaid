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
	var str string
	for _, s := range sink.sinks {
		str += s.Name() + ","
	}
	if l := len(str); l > 0 {
		str = str[:l-1]
	}
	return "multi[" + str + "]"
}

// Register registers a new backend sink. It is not thread-safe
func (sink *MultiSink) Register(s Sink) {
	sink.sinks = append(sink.sinks, s)
}

// Publish satisfies the Sink interface
func (sink *MultiSink) Publish(seri ...tsdb.Series) error {
	l := len(sink.sinks)
	if l == 0 {
		sink.log.Info("sink: no backend")
		return nil
	}
	errs := make([]error, 0, l)
	for _, s := range sink.sinks {
		err := s.Publish(seri...)
		if err != nil {
			errs = append(errs, errors.Wrap(err, s.Name()))
		} else {
			sink.log.Info("published",
				zap.String("sink", s.Name()),
				zap.Int("count", len(seri)))
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
