package fl

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

const testDelimitedStr = "foo,bar,bars,,flip"

type testOp struct {
	in string

	eOp    string
	eField string
}

var testOps = []testOp{
	testOp{in: "foo", eField: "foo", eOp: NoOp},
	testOp{in: "foo,bar,baz", eField: "foo,bar,baz", eOp: NoOp},
	testOp{in: ":foo", eField: ":foo", eOp: NoOp},
	testOp{in: "xx:foo", eField: "xx:foo", eOp: NoOp},
	testOp{in: "gt:foo", eField: "foo", eOp: OpGreater},
	testOp{in: "lt:foo", eField: "foo", eOp: OpLess},
	testOp{in: "le:foo", eField: "foo", eOp: OpLessEqual},
	testOp{in: "!foo", eField: "foo", eOp: OpNotEqual},
	testOp{in: "!foo", eField: "foo", eOp: OpNotEqual},
}

func Test_parseDelimited(t *testing.T) {
	out := parseDelimited(testDelimitedStr, ",")
	assert.Equal(t, 4, len(out))
}

func Test_parseOp(t *testing.T) {
	for _, top := range testOps {
		op, field := parseOp(top.in)
		assert.Equal(t, top.eOp, op)
		assert.Equal(t, top.eField, field)
	}
}

func Test_ParseValue(t *testing.T) {
	for _, top := range testOps {
		op, fields := ParseValue(top.in)
		assert.Equal(t, top.eOp, op)
		assert.Equal(t, top.eField, strings.Join(fields, ListDelimiter))
	}
}
