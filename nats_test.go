package leaderelection

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNatsRaftLE(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping TestNatsRaftLE in short mode")
	}
	//start nats server first
	//docker run -p 4222:4222 -ti nats:latest
	size := 3
	name := "test1"
	nodes := []LE{}
	for i := 0; i < 3; i++ {
		node, err := NewNatsLE(name, size)
		assert.NoError(t, err)
		nodes = append(nodes, node)
	}
	countLeader := 0
	for i := range nodes {
		if nodes[i].AmILeader() {
			countLeader++
		}
	}
	assert.Equal(t, 1, countLeader, "there should be only 1 leader")
	// for i := range nodes {
	// 	Log(nodes[i].Node.Id(), nodes[i].Node.State().String())
	// }
}
