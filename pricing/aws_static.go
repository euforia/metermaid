package pricing

import (
	"errors"
	"fmt"
	"time"

	"github.com/euforia/metermaid/tsdb"
)

var (
	awsOnDemandPrices = map[string]map[string]float64{
		"us-west-2": awsOndemandUSWest2Prices,
	}
)

type AWSOnDemandStaticPricer struct{}

// Name ...
func (pp *AWSOnDemandStaticPricer) Name() string {
	return "aws-ondemand-static"
}

// History returns the point in time price for both start and end
func (pp *AWSOnDemandStaticPricer) History(start, end time.Time, filter map[string]string) (tsdb.DataPoints, error) {
	region := filter["Region"]
	regPrices, ok := awsOnDemandPrices[region]
	if ok {
		instType := filter["InstanceType"]
		price, ok := regPrices[instType]
		if ok {
			return tsdb.DataPoints{
				tsdb.DataPoint{Timestamp: uint64(start.UnixNano()), Value: price},
				tsdb.DataPoint{Timestamp: uint64(end.UnixNano()), Value: price},
			}, nil
		}
		return nil, fmt.Errorf("no pricing for instance: %s", instType)
	}
	return nil, errors.New("no pricing for region: " + region)
}

// us-west-2
// 03/11/2019
var awsOndemandUSWest2Prices = map[string]float64{
	"a1.medium":     0.0255,
	"a1.large":      0.051,
	"a1.xlarge":     0.102,
	"a1.2xlarge":    0.204,
	"a1.4xlarge":    0.408,
	"t3.nano":       0.0052,
	"t3.micro":      0.0104,
	"t3.small":      0.0208,
	"t3.medium":     0.0416,
	"t3.large":      0.0832,
	"t3.xlarge":     0.1664,
	"t3.2xlarge":    0.3328,
	"t2.nano":       0.0058,
	"t2.micro":      0.0116,
	"t2.small":      0.023,
	"t2.medium":     0.0464,
	"t2.large":      0.0928,
	"t2.xlarge":     0.1856,
	"t2.2xlarge":    0.3712,
	"m5.large":      0.096,
	"m5.xlarge":     0.192,
	"m5.2xlarge":    0.384,
	"m5.4xlarge":    0.768,
	"m5.12xlarge":   2.304,
	"m5.24xlarge":   4.608,
	"m5.metal":      4.608,
	"m5a.large":     0.086,
	"m5a.xlarge":    0.172,
	"m5a.2xlarge":   0.344,
	"m5a.4xlarge":   0.688,
	"m5a.12xlarge":  2.064,
	"m5a.24xlarge":  4.128,
	"m5d.large":     0.113,
	"m5d.xlarge":    0.226,
	"m5d.2xlarge":   0.452,
	"m5d.4xlarge":   0.904,
	"m5d.12xlarge":  2.712,
	"m5d.24xlarge":  5.424,
	"m5d.metal":     5.424,
	"m4.large":      0.10,
	"m4.xlarge":     0.20,
	"m4.2xlarge":    0.40,
	"m4.4xlarge":    0.80,
	"m4.10xlarge":   2.00,
	"m4.16xlarge":   3.20,
	"c5.large":      0.085,
	"c5.xlarge":     0.17,
	"c5.2xlarge":    0.34,
	"c5.4xlarge":    0.68,
	"c5.9xlarge":    1.53,
	"c5.18xlarge":   3.06,
	"c5d.large":     0.096,
	"c5d.xlarge":    0.192,
	"c5d.2xlarge":   0.384,
	"c5d.4xlarge":   0.768,
	"c5d.9xlarge":   1.728,
	"c5d.18xlarge":  3.456,
	"c5n.large":     0.108,
	"c5n.xlarge":    0.216,
	"c5n.2xlarge":   0.432,
	"c5n.4xlarge":   0.864,
	"c5n.9xlarge":   1.944,
	"c5n.18xlarge":  3.888,
	"c4.large":      0.10,
	"c4.xlarge":     0.199,
	"c4.2xlarge":    0.398,
	"c4.4xlarge":    0.796,
	"c4.8xlarge":    1.591,
	"p3.2xlarge":    3.06,
	"p3.8xlarge":    12.24,
	"p3.16xlarge":   24.48,
	"p3dn.24xlarge": 31.212,
	"p2.xlarge":     0.90,
	"p2.8xlarge":    7.20,
	"p2.16xlarge":   14.40,
	"g3.4xlarge":    1.14,
	"g3.8xlarge":    2.28,
	"g3.16xlarge":   4.56,
	"g3s.xlarge":    0.75,
	"x1.16xlarge":   6.669,
	"x1.32xlarge":   13.338,
	"x1e.xlarge":    0.834,
	"x1e.2xlarge":   1.668,
	"x1e.4xlarge":   3.336,
	"x1e.8xlarge":   6.672,
	"x1e.16xlarge":  13.344,
	"x1e.32xlarge":  26.688,
	"r5.large":      0.126,
	"r5.xlarge":     0.252,
	"r5.2xlarge":    0.504,
	"r5.4xlarge":    1.008,
	"r5.12xlarge":   3.024,
	"r5.24xlarge":   6.048,
	"r5.metal":      6.048,
	"r5a.large":     0.113,
	"r5a.xlarge":    0.226,
	"r5a.2xlarge":   0.452,
	"r5a.4xlarge":   0.904,
	"r5a.12xlarge":  2.712,
	"r5a.24xlarge":  5.424,
	"r5d.large":     0.144,
	"r5d.xlarge":    0.288,
	"r5d.2xlarge":   0.576,
	"r5d.4xlarge":   1.152,
	"r5d.12xlarge":  3.456,
	"r5d.24xlarge":  6.912,
	"r4.large":      0.133,
	"r4.xlarge":     0.266,
	"r4.2xlarge":    0.532,
	"r4.4xlarge":    1.064,
	"r4.8xlarge":    2.128,
	"r4.16xlarge":   4.256,
	"z1d.large":     0.186,
	"z1d.xlarge":    0.372,
	"z1d.2xlarge":   0.744,
	"z1d.3xlarge":   1.116,
	"z1d.6xlarge":   2.232,
	"z1d.12xlarge":  4.464,
	"z1d.metal":     4.464,
	"i3.large":      0.156,
	"i3.xlarge":     0.312,
	"i3.2xlarge":    0.624,
	"i3.4xlarge":    1.248,
	"i3.8xlarge":    2.496,
	"i3.16xlarge":   4.992,
	"i3.metal":      4.992,
	"h1.2xlarge":    0.468,
	"h1.4xlarge":    0.936,
	"h1.8xlarge":    1.872,
	"h1.16xlarge":   3.744,
	"d2.xlarge":     0.69,
	"d2.2xlarge":    1.38,
	"d2.4xlarge":    2.76,
	"d2.8xlarge":    5.52,
	"f1.2xlarge":    1.65,
	"f1.4xlarge":    3.30,
	"f1.16xlarge":   13.20,
}
