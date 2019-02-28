package sink

import (
	"os"

	"github.com/euforia/metermaid/tsdb"
	"github.com/zorkian/go-datadog-api"
)

type DataDogSink struct {
	client *datadog.Client

	// default tags added to all metrics
	// tags map[string]string
}

func NewDataDogSink(apiKey, appKey string) *DataDogSink {
	if apiKey == "" {
		apiKey = os.Getenv("DD_API_KEY")
	}
	if appKey == "" {
		appKey = os.Getenv("DD_APP_KEY")
	}

	return &DataDogSink{client: datadog.NewClient(apiKey, appKey)}
}

func (dd *DataDogSink) Name() string {
	return "datadog"
}

func (dd *DataDogSink) Publish(seri ...tsdb.Series) error {
	metrics := make([]datadog.Metric, len(seri))
	for i, s := range seri {
		metric := datadog.Metric{
			Metric: &s.Name,
			Tags:   formatTags(s.Meta),
			Points: make([]datadog.DataPoint, len(s.Data)),
		}
		for j, dp := range s.Data {
			ts := float64(dp.Timestamp) / 1e9
			val := dp.Value
			metric.Points[j] = datadog.DataPoint{&ts, &val}
		}
		metrics[i] = metric
	}
	return dd.client.PostMetrics(metrics)
}

func formatTags(tags map[string]string) []string {
	list := make([]string, 0, len(tags))
	for k, v := range tags {
		list = append(list, k+":"+v)
	}
	return list
}
