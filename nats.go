/**
elect based on raft algorithm
*/
package leaderelection

import (
	"os"
	"strings"
	"time"

	log "github.com/kanguki/log"
	"github.com/nats-io/graft"
	"github.com/nats-io/nats.go"
)

type NatsLE struct {
	Node *graft.Node
	Size int    //number of nodes in cluster
	Name string //name of the election
}

func NewNatsLE(name string, size int) (*NatsLE, error) {
	ci := graft.ClusterInfo{Name: name, Size: size}
	do := nats.GetDefaultOptions()
	if quorum := os.Getenv(NATS_QUORUM); quorum != "" {
		do.Servers = strings.Split(quorum, ",")
	}
	rpc, err := graft.NewNatsRpc(&do)
	if err != nil {
		log.Log("error creating new NatsRaftLE: %v", err)
		return nil, err
	}
	errChan := make(chan error)
	stateChangeChan := make(chan graft.StateChange)
	handler := graft.NewChanHandler(stateChangeChan, errChan)

	node, err := graft.New(ci, handler, rpc, "/tmp/graft.log")
	if err != nil {
		log.Log("error creating new NatsRaftLE: %v", err)
		return nil, err
	}
	go func() {
		for {
			select {
			case sc := <-stateChangeChan:
				log.Debug("node %v's state changed from %v to %v", node.Id(), sc.From, sc.To)
			case err := <-errChan:
				log.Debug("node %v received error %v", node.Id(), err)
			}
		}
	}()

	return &NatsLE{Node: node, Size: size, Name: name}, nil
}

func (n *NatsLE) AmILeader() bool {
	noLeader := func() bool {
		return n.Node.State() != graft.LEADER && n.Node.Leader() == ""
	}
	waitCount := 0
	for noLeader() && waitCount < 30 { //if it takes too
		time.Sleep(1 * time.Second)
		waitCount++
		log.Debug("%v: 1 sec passed by without leader in cluster", n.Node.Id())
		//TODO: integrate notification
	}
	return n.Node.State() == graft.LEADER
}

func (n *NatsLE) CleanResource() {
	n.Node.Close()
}
