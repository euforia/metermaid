package collector

import (
	"context"
	"strings"
	"time"

	"go.uber.org/zap"
)

type Engine struct {
	collectors map[string]*collector
	out        chan []RunStats

	ctx    context.Context
	cancel context.CancelFunc

	log *zap.Logger
}

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
func (eng *Engine) Register(c Collector, interval time.Duration) {
	cltr := &collector{
		interval: interval,
		out:      eng.out,
		bc:       c,
		log:      eng.log,
	}
	eng.collectors[c.Name()] = cltr
	eng.log.Info("collector registered", zap.String("name", c.Name()))
}

// RunStats returns a channel containing newly available runtimes
func (eng *Engine) RunStats() <-chan []RunStats {
	return eng.out
}

// Start starts each registered collector
func (eng *Engine) Start() {
	var kstr string
	for k := range eng.collectors {
		kstr += k + ","
	}
	kstr = strings.TrimSuffix(kstr, ",")
	eng.log.Info("engine starting", zap.String("collectors", kstr))

	for _, c := range eng.collectors {
		go c.run(eng.ctx)
	}
}

// Stop signals all collectors to stop and waits for them to exit before
// returning
func (eng *Engine) Stop() {
	select {
	case <-eng.ctx.Done():
		// already stopped
	default:
		eng.cancel()
		close(eng.out)
	}
}
