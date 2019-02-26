package pricing

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/euforia/metermaid/tsdb"
)

// AWSSpotPricer provides aws pricing history
type AWSSpotPricer struct{}

// NewAWSSpotPricer returns a new instance of AWSSpotPricer
func NewAWSSpotPricer() *AWSSpotPricer {
	return &AWSSpotPricer{}
}

// Name returns the name of the provider
func (pp *AWSSpotPricer) Name() string {
	return "aws-spot"
}

// History returns the price history given the filter. Region is a required
// filter key.
func (pp *AWSSpotPricer) History(start, end time.Time, filter map[string]string) (tsdb.DataPoints, error) {
	return pp.SpotHistory(start, end, filter)
}

// SpotHistory returns the spot price history given the filter. Region is a
// required filter key.
func (pp *AWSSpotPricer) SpotHistory(start, end time.Time, filter map[string]string) (tsdb.DataPoints, error) {
	sess := session.New(&aws.Config{Region: aws.String(filter["Region"])})
	svc := ec2.New(sess)

	input := &ec2.DescribeSpotPriceHistoryInput{
		EndTime: &end,
		ProductDescriptions: []*string{
			// TODO: make part of the filter. Could be Windows SUSE Linux based on
			// what AWS has
			aws.String("Linux/UNIX"),
		},
		StartTime: &start,
	}

	if itype, ok := filter["InstanceType"]; ok {
		input.InstanceTypes = []*string{
			aws.String(itype),
		}
	}

	if az, ok := filter["AvailabilityZone"]; ok {
		input.AvailabilityZone = aws.String(az)
	}

	result, err := svc.DescribeSpotPriceHistory(input)
	if err != nil {
		return nil, err
	}

	prices := make(tsdb.DataPoints, 0, len(result.SpotPriceHistory))
	for _, sp := range result.SpotPriceHistory {
		pf, er := strconv.ParseFloat(*sp.SpotPrice, 64)
		if er != nil {
			err = er
			continue
		}

		prices = prices.Insert(tsdb.DataPoint{
			Value:     pf,
			Timestamp: uint64(sp.Timestamp.UnixNano()),
			// Meta: map[string]string{
			// 	"Region":             filter["Region"],
			// 	"AvailabilityZone":   *sp.AvailabilityZone,
			// 	"ProductDescription": *sp.ProductDescription,
			// 	"InstanceType":       *sp.InstanceType,
			// },
		})
		// prices = append(prices, price)
	}
	return prices, err
}
