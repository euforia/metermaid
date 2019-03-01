package main

import (
	"errors"
	"io/ioutil"
	"time"

	"github.com/hashicorp/hcl"
)

type config struct {
	Collectors map[string]*collectorConfig
	Sinks      map[string]*sinkConfig
}

func parseConfigFile() (*config, error) {
	if *confFile == "" {
		return nil, errors.New("config file required")
	}

	var conf config
	b, err := ioutil.ReadFile(*confFile)
	if err == nil {
		err = hcl.Unmarshal(b, &conf)
	}

	return &conf, err
}

type sinkConfig struct {
	Config map[string]interface{}
}

type collectorConfig struct {
	Interval string
	Config   map[string]interface{}
}

func (cc *collectorConfig) IntervalDuration() time.Duration {
	if d, err := time.ParseDuration(cc.Interval); err == nil {
		return d
	}
	return -1
}
