package pricing

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_AWSPricing(t *testing.T) {
	p := NewAWSPricer()

	now := time.Now()
	start := time.Unix(now.Unix()-14400, 0)
	pricing, err := p.History(start, now, map[string]string{
		"Region":       "us-west-2",
		"InstanceType": "m3.xlarge",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(pricing))
	t.Log(pricing)
}
