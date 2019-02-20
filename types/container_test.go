package types

import (
	"testing"
	"time"

	"github.com/euforia/metermaid/fl"
	"github.com/stretchr/testify/assert"
)

var testC = &Container{
	Name:   "bar",
	Create: time.Now().UnixNano(),
	Memory: 100,
	Labels: map[string]string{
		"service":   "my-service",
		"component": "my-component",
	},
}

func Test_Container(t *testing.T) {
	assert.False(t, testC.MatchField("service", fl.Filter{Operator: fl.NoOp, Values: []string{"foo"}}))
	assert.False(t, testC.MatchField("service", fl.Filter{Operator: fl.NoOp, Values: []string{"my", "foo"}}))

	assert.True(t, testC.MatchField("service", fl.Filter{Operator: fl.NoOp, Values: []string{"my-service"}}))
	assert.True(t, testC.MatchField("service", fl.Filter{Operator: fl.NoOp, Values: []string{"my-service", "foo"}}))

	assert.False(t, testC.MatchField("Name", fl.Filter{Operator: fl.OpNotEqual, Values: []string{"bar", "foo"}}))
	assert.False(t, testC.MatchField("Name", fl.Filter{Operator: fl.OpNotEqual, Values: []string{"bar"}}))
	assert.True(t, testC.MatchField("Name", fl.Filter{Operator: fl.NoOp, Values: []string{"bar"}}))

	assert.True(t, testC.MatchField("Name", fl.Filter{Operator: fl.NoOp, Values: []string{"bar"}}))
}

func Test_MatchTime(t *testing.T) {
	assert.True(t, testC.MatchField("Create", fl.Filter{
		Operator: fl.OpLessEqual,
		Values:   []string{time.Now().Format(time.RFC3339Nano)},
	}))
	assert.True(t, testC.MatchField("Create", fl.Filter{
		Operator: fl.OpLess,
		Values:   []string{time.Now().Format(time.RFC3339Nano)},
	}))
	assert.False(t, testC.MatchField("Create", fl.Filter{
		Operator: fl.NoOp,
		Values:   []string{time.Now().Format(time.RFC3339Nano)},
	}))
	assert.False(t, testC.MatchField("Create", fl.Filter{
		Operator: fl.OpGreater,
		Values:   []string{time.Now().Format(time.RFC3339Nano)},
	}))
	assert.False(t, testC.MatchField("Create", fl.Filter{
		Operator: fl.OpGreaterEqual,
		Values:   []string{time.Now().Format(time.RFC3339Nano)},
	}))

}
