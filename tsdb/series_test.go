package tsdb

import (
	"testing"

	"github.com/euforia/metermaid/types"
)

func Test_Series(t *testing.T) {
	s := &Series{Name: "foo", Meta: types.Meta{"bar": "chocolate"}}
	t.Error(s)
}
