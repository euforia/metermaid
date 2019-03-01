package collector

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func Test_Engine(t *testing.T) {
	logger := zap.NewExample()
	eng := NewEngine(logger)

	dc := &DockerCollector{}
	eng.Register(dc, &Config{
		Interval: 10 * time.Second,
		Config:   map[string]interface{}{},
	})

	eng.Start()

	stats := <-eng.RunStats()
	assert.NotEqual(t, 0, len(stats))
	eng.Stop()
}
