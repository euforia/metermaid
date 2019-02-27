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
	// ctx := context.Background()

	dc := &DockerCollector{}
	dc.Init(map[string]interface{}{})
	eng.Register(dc, 10*time.Second)

	eng.Start()

	stats := <-eng.RunStats()
	assert.NotEqual(t, 0, len(stats))
	eng.Stop()
}
