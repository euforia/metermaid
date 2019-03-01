package node

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_Node(t *testing.T) {
	node, _ := New()
	assert.NotEmpty(t, node.CPUShares)
	assert.NotEmpty(t, node.Memory)
}
