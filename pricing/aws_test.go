package pricing

import (
	"testing"
	"time"

	"go.uber.org/zap"

	"github.com/euforia/metermaid/node"
	"github.com/stretchr/testify/assert"
)

func Test_AWSPriceProvider_SpotHistory(t *testing.T) {
	p := NewAWSSpotPricer()

	now := time.Now()
	start := time.Unix(now.Unix()-14400, 0)
	pricing, err := p.SpotHistory(start, now, map[string]string{
		"Region":       "us-west-2",
		"InstanceType": "m3.xlarge",
	})
	assert.Nil(t, err)
	assert.NotEqual(t, 0, len(pricing))
	t.Log(pricing)
}

func Test_AWSPriceProvider_OnDemandHistory(t *testing.T) {
	p := NewAWSOnDemandPricer()
	start, _ := time.Parse("2006-01-02T15:04", "2019-02-19T00:00")
	end, _ := time.Parse("2006-01-02T15:04", "2019-02-21T00:00")
	pricing, err := p.History(start, end, map[string]string{
		"Region":       "us-west-2",
		"InstanceType": "r4.xlarge",
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(pricing))
}

func Test_Pricer(t *testing.T) {
	start, _ := time.Parse("2006-01-02T15:04", "2019-02-19T00:00")
	end, _ := time.Parse("2006-01-02T15:04", "2019-02-21T00:00")
	pricer := NewPricer(
		NewAWSSpotPricer(),
		node.Node{Meta: map[string]string{
			"Region":           "us-west-2",
			"InstanceType":     "r4.xlarge",
			"AvailabilityZone": "us-west-2b",
		}},
		zap.NewExample(),
	)

	assert.NotEqual(t, 0, len(pricer.cache))

	prices, err := pricer.History(start, end)
	assert.Nil(t, err)
	assert.EqualValues(t, end.UnixNano(), prices.Last().Timestamp)
}
