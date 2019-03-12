package pricing

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/pricing"
	"github.com/euforia/metermaid/tsdb"
)

// On demand pricing api needs the name instead of id
var regionMap = map[string]string{
	"us-west-2": "US West (Oregon)",
	"us-east-1": "US East (N. Virginia)",
}

// AWSOnDemandPricer provides aws pricing history
type AWSOnDemandPricer struct{}

// NewAWSOnDemandPricer returns a new instance of AWSOnDemandPricer
func NewAWSOnDemandPricer() *AWSOnDemandPricer {
	return &AWSOnDemandPricer{}
}

// Name returns the name of the pricer
func (pp *AWSOnDemandPricer) Name() string {
	return "aws-ondemand"
}

// History returns the price history given the filter. Region is a required
// filter key.
func (pp *AWSOnDemandPricer) History(start, end time.Time, filter map[string]string) (tsdb.DataPoints, error) {
	return pp.OnDemandHistory(filter)
}

// OnDemandHistory returns the ondemand price for the node
func (pp *AWSOnDemandPricer) OnDemandHistory(filter map[string]string) (tsdb.DataPoints, error) {
	priceList, err := getNonSpotPrice(filter)
	if err == nil {
		var dps tsdb.DataPoints
		if dps, err = parseOnDemandPriceData(priceList); err == nil {
			return dps, nil
		}
	}
	return nil, err
}

func getNonSpotPrice(filter map[string]string) ([]aws.JSONValue, error) {
	// Service is currently only available in us-east-1
	sess := session.New(&aws.Config{Region: aws.String("us-east-1")})
	prc := pricing.New(sess)

	filters := make([]*pricing.Filter, 0, 6)
	if region, ok := filter["Region"]; ok {
		filters = append(filters, &pricing.Filter{
			Type:  aws.String(pricing.FilterTypeTermMatch),
			Field: aws.String("location"),
			Value: aws.String(regionMap[region]),
		})
	}

	if val, ok := filter["InstanceType"]; ok {
		filters = append(filters, &pricing.Filter{
			Type:  aws.String(pricing.FilterTypeTermMatch),
			Field: aws.String("InstanceType"),
			Value: aws.String(val),
		})
	}

	// Default is shared tenancy
	tenancy := "Shared"
	if val, ok := filter["Tenancy"]; ok {
		tenancy = val
	}
	// Add defaults including tenancy
	filters = append(filters,
		&pricing.Filter{
			Type:  aws.String(pricing.FilterTypeTermMatch),
			Field: aws.String("Tenancy"),
			Value: aws.String(tenancy),
		},
		&pricing.Filter{
			Type:  aws.String(pricing.FilterTypeTermMatch),
			Field: aws.String("CapacityStatus"),
			Value: aws.String("Used"),
		},
		&pricing.Filter{
			Type:  aws.String(pricing.FilterTypeTermMatch),
			Field: aws.String("Operation"),
			Value: aws.String("RunInstances"),
		},
		&pricing.Filter{
			Type:  aws.String(pricing.FilterTypeTermMatch),
			Field: aws.String("OperatingSystem"),
			Value: aws.String("Linux"),
		},
	)
	// fmt.Println(filters)
	resp, err := prc.GetProducts(&pricing.GetProductsInput{
		ServiceCode: aws.String("AmazonEC2"),
		Filters:     filters,
	})
	if err == nil {
		return resp.PriceList, nil
	}
	return nil, err
}

func parseOnDemandPriceData(priceList []aws.JSONValue) (dps tsdb.DataPoints, err error) {
	for _, j := range priceList {
		terms := j["terms"].(map[string]interface{})
		ondemand := terms["OnDemand"].(map[string]interface{})
		for _, iv := range ondemand {
			d, err := parseNonSpotJSON(iv.(map[string]interface{}))
			if err == nil {
				dps = dps.Insert(d...)
				continue
			}
			return nil, err
		}
	}
	return
}

func parseNonSpotJSON(v map[string]interface{}) (tsdb.DataPoints, error) {
	effDate, err := time.Parse("2006-01-02T15:04:05Z", v["effectiveDate"].(string))
	if err != nil {
		return nil, err
	}

	ts := uint64(effDate.UnixNano())
	dps := make(tsdb.DataPoints, 0)

	pdim := v["priceDimensions"].(map[string]interface{})
	for _, ivv := range pdim {
		a := ivv.(map[string]interface{})
		ppu := a["pricePerUnit"].(map[string]interface{})
		value, err := strconv.ParseFloat(ppu["USD"].(string), 64)
		// return
		if err == nil {
			dps = dps.Insert(tsdb.DataPoint{Timestamp: ts, Value: value})
			continue
		}
		return nil, err
	}
	return dps, nil
}
