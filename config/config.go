package config

import (
	"io/ioutil"
	"time"

	"github.com/hashicorp/hcl"
)

type Config struct {
	Collectors map[string]*CollectorConfig
	Sinks      map[string]*SinkConfig
}

func ParseFile(filename string) (*Config, error) {
	var conf Config
	b, err := ioutil.ReadFile(filename)
	if err == nil {
		err = hcl.Unmarshal(b, &conf)
	}
	return &conf, err
}

type SinkConfig struct{}

type CollectorConfig struct {
	Interval string
	Config   map[string]interface{}
}

func (cc *CollectorConfig) IntervalDuration() time.Duration {
	if d, err := time.ParseDuration(cc.Interval); err == nil {
		return d
	}
	return -1
}
