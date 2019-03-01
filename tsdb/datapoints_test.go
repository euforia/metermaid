package tsdb

import (
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"
)

type dpsTestCase struct {
	dps DataPoints
	// Get
	s uint64
	e uint64

	expected int
	// Encompass
	enc bool
}

var dpsGetTests = []dpsTestCase{
	dpsTestCase{
		DataPoints{},
		10, 20, 0, false,
	},
	dpsTestCase{
		nil,
		10, 20, 0, false,
	},
	dpsTestCase{
		DataPoints{
			DataPoint{10, 1},
		},
		9, 11, 1, false,
	},
	dpsTestCase{
		DataPoints{
			DataPoint{10, 1},
		},
		10, 11, 1, false,
	},
	dpsTestCase{
		DataPoints{
			DataPoint{10, 1},
		},
		11, 12, 0, false,
	},
	dpsTestCase{
		DataPoints{
			DataPoint{10, 1},
			DataPoint{12, 1},
		},
		9, 11, 1, false,
	},
	dpsTestCase{
		DataPoints{
			DataPoint{10, 1},
		},
		9, 11, 1, false,
	},
	dpsTestCase{
		DataPoints{
			DataPoint{10, 1},
			DataPoint{12, 1},
			DataPoint{13, 1},
			DataPoint{17, 1},
			DataPoint{20, 1},
		},
		11, 21, 4, false,
	},
	dpsTestCase{
		DataPoints{
			DataPoint{10, 1},
			DataPoint{12, 1},
			DataPoint{13, 1},
			DataPoint{17, 1},
			DataPoint{20, 1},
		},
		20, 21, 1, false,
	},
}

var dpsDedupTests = []dpsTestCase{
	dpsTestCase{
		dps: DataPoints{
			DataPoint{10, 1},
			DataPoint{12, 1},
			DataPoint{20, 1},
			DataPoint{20, 1},
			DataPoint{20, 1},
		},
		expected: 3,
	},
	dpsTestCase{
		dps: DataPoints{
			DataPoint{10, 1},
			DataPoint{12, 1},
			DataPoint{20, 1},
			DataPoint{12, 1},
			DataPoint{20, 1},
		},
		expected: 3,
	},
}

var dpsPerTests = dpsTestCase{
	dps: DataPoints{
		DataPoint{100, 1},
		DataPoint{900, 1},
		DataPoint{1200, 1},
		DataPoint{1353, 1},
		DataPoint{1534, 1},
		DataPoint{1800, 1},
	},
	expected: 19,
}

var dpsEncTests = []dpsTestCase{
	dpsTestCase{
		dps: DataPoints{
			DataPoint{100, 1},
			DataPoint{900, 1},
			DataPoint{1200, 1},
			DataPoint{1353, 1},
			DataPoint{1534, 1},
			DataPoint{1800, 1},
		},
		s: 100, e: 1300,
		enc: true,
	},
	dpsTestCase{
		dps: DataPoints{
			DataPoint{100, 1},
			DataPoint{900, 1},
			DataPoint{1200, 1},
			DataPoint{1353, 1},
			DataPoint{1534, 1},
			DataPoint{1800, 1},
		},
		s: 50, e: 1300,
		enc: false,
	},
	dpsTestCase{
		dps: DataPoints{
			DataPoint{100, 1},
			DataPoint{900, 1},
			DataPoint{1200, 1},
			DataPoint{1353, 1},
			DataPoint{1534, 1},
			DataPoint{1800, 1},
		},
		s: 120, e: 3000,
		enc: false,
	},
	dpsTestCase{
		dps: DataPoints{
			DataPoint{100, 1},
			DataPoint{900, 1},
			DataPoint{1200, 1},
			DataPoint{1353, 1},
			DataPoint{1534, 1},
			DataPoint{1800, 1},
		},
		s: 50, e: 3000,
		enc: false,
	},
	dpsTestCase{
		s: 50, e: 3000,
		enc: false,
	},
}

func Test_Datapoints_Get(t *testing.T) {
	for _, tc := range dpsGetTests {
		s := tc.dps.Get(tc.s, tc.e)
		assert.Equal(t, tc.expected, s.Len(), "%v", s)
	}
}

func Test_Datapoints_Dedup(t *testing.T) {
	for _, tc := range dpsDedupTests {
		sort.Sort(tc.dps)
		deduped := tc.dps.Dedup()
		assert.Equal(t, tc.expected, deduped.Len())
	}
}

func Test_Datapoints_Scale_Sum(t *testing.T) {
	assert.Equal(t, 5.0, dpsDedupTests[0].dps.Sum())
	dps := dpsDedupTests[0].dps.Scale(.5)
	for _, dp := range dps {
		assert.Equal(t, .5, dp.Value)
	}
	assert.Equal(t, 2.5, dps.Sum())
}
