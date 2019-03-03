package collector

import (
	"context"
	"strings"

	"go.uber.org/zap"
)

// Engine manages registered collectors
type Engine struct {
	collectors map[string]*collector
	out        chan []RunStats

	ctx    context.Context
	cancel context.CancelFunc

	log *zap.Logger
}

// NewEngine returns a new Engine instance
func NewEngine(logger *zap.Logger) *Engine {
	eng := &Engine{
		out:        make(chan []RunStats, 32),
		collectors: make(map[string]*collector),
		log:        logger,
	}
	eng.ctx, eng.cancel = context.WithCancel(context.Background())

	return eng
}

// Register registers a new Collector run at the given interval
func (eng *Engine) Register(c Collector, conf *Config) error {
	err := conf.Validate()
	if err != nil {
		return err
	}

	conf.Logger = eng.log
	if err = c.Init(conf); err != nil {
		return err
	}

	eng.collectors[c.Name()] = &collector{
		interval: conf.Interval,
		out:      eng.out,
		bc:       c,
		log:      eng.log,
		done:     make(chan struct{}),
	}

	eng.log.Info("collector registered", zap.String("name", c.Name()))
	return nil
}

// RunStats returns a channel containing newly available runtimes
func (eng *Engine) RunStats() <-chan []RunStats {
	return eng.out
}

// Start starts each registered collector
func (eng *Engine) Start() {
	for _, c := range eng.collectors {
		go c.run(eng.ctx)
	}

	var kstr string
	for k := range eng.collectors {
		kstr += k + ","
	}
	kstr = strings.TrimSuffix(kstr, ",")
	eng.log.Info("engine starting", zap.String("collectors", kstr))
}

// Stop signals all collectors to stop and waits for them to exit before
// returning
func (eng *Engine) Stop() {
	select {
	case <-eng.ctx.Done():
		// already stopped
	default:
		eng.cancel()
		for _, c := range eng.collectors {
			c.stop()
		}
		close(eng.out)
	}
}
