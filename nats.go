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
	Node                *graft.Node
	Size                int    //number of nodes in cluster
	Name                string //name of the election
	TimeoutDecideLeader int    //default 30
}

func NewNatsLE(opts LeOpts) (*NatsLE, error) {
	ci := graft.ClusterInfo{Name: opts.Name, Size: opts.Size}
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
	timeout := 30
	if opts.TimeoutDecideLeader != 0 {
		timeout = opts.TimeoutDecideLeader
	}
	return &NatsLE{Node: node, Size: opts.Size, Name: opts.Name, TimeoutDecideLeader: timeout}, nil
}

func (n *NatsLE) AmILeader() bool {
	noLeader := func() bool {
		return n.Node.State() != graft.LEADER && n.Node.Leader() == ""
	}
	stop := make(chan interface{}, 1)
	go func() {
		for {
			select {
			case <-time.After(time.Duration(n.TimeoutDecideLeader) * time.Second):
				stop <- struct{}{}
			}
		}
	}()
L:
	for noLeader() {
		select {
		case <-stop:
			log.Debug("%v: %v sec passed by without leader in cluster", n.Node.Id(), n.TimeoutDecideLeader)
			//TODO: integrate notification
			break L
		case <-time.After(time.Second):
			log.Debug("%v: 1 sec passed by without leader in cluster", n.Node.Id())
		}
	}
	return n.Node.State() == graft.LEADER
}

func (n *NatsLE) CleanResource() {
	n.Node.Close()
}
