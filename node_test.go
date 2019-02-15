package metermaid

import "testing"

func Test_Node(t *testing.T) {
	node := NewNode()
	t.Logf("%+v", node)
}
