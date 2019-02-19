package pricing

import (
	"errors"
	"sort"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
)

// AWSPriceProvider provides aws pricing history
type AWSPriceProvider struct{}

// NewAWSPriceProvider returns a new instance of AWSPriceProvider
func NewAWSPriceProvider() *AWSPriceProvider {
	return &AWSPriceProvider{}
}

// History returns the price history given the filter. Region is a required
// filter key.
func (pricing *AWSPriceProvider) History(start, end time.Time, filter map[string]string) ([]*Price, error) {
	// TODO:
	// - determine if the instance is spot, reserved or on-demand
	//
	return pricing.SpotHistory(start, end, filter)
}

// OnDemand returns the ondemand price for the node
func (pricing *AWSPriceProvider) OnDemand() (*Price, error) {
	return nil, errors.New("To be implemented")
}

// SpotHistory returns the spot price history given the filter. Region is a
// required filter key.
func (pricing *AWSPriceProvider) SpotHistory(start, end time.Time, filter map[string]string) ([]*Price, error) {
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

	prices := make([]*Price, 0, len(result.SpotPriceHistory))
	for _, sp := range result.SpotPriceHistory {
		pf, er := strconv.ParseFloat(*sp.SpotPrice, 64)
		if er != nil {
			err = er
			continue
		}

		price := &Price{
			Price:     pf,
			Timestamp: uint64(sp.Timestamp.UnixNano()),
			Meta: map[string]string{
				"AvailabilityZone":   *sp.AvailabilityZone,
				"ProductDescription": *sp.ProductDescription,
				"InstanceType":       *sp.InstanceType,
			},
		}
		prices = append(prices, price)
	}

	sort.Sort(sortedPrices(prices))
	return prices, err
}
