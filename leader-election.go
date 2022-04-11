package leaderelection

import "fmt"

type LE interface {
	AmILeader() bool
	CleanResource() //defer in main
}

type Base string

const (
	NATS Base = "NATS"
)

type LeOpts struct {
	Base                Base
	Name                string
	Size                int
	TimeoutDecideLeader int //default 30
}

func NewLE(opts LeOpts) (LE, error) {
	switch opts.Base {
	case NATS:
		return NewNatsLE(opts)
	default:
		err := fmt.Sprintf("error creating new leaderElector: unsupported base: %v", opts.Base)
		return nil, fmt.Errorf(err)
	}
}
