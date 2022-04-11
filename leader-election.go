package leaderelection

import "fmt"

type LE interface {
	AmILeader() bool
}

type Base string

const (
	NATS Base = "NATS"
)

type LeOpts struct {
	Base Base
	Name string
	Size int
}

func NewLE(opts LeOpts) (LE, error) {
	switch opts.Base {
	case NATS:
		return NewNatsLE(opts.Name, opts.Size)
	default:
		err := fmt.Sprintf("error creating new leaderElector: unsupported base: %v", opts.Base)
		return nil, fmt.Errorf(err)
	}
}
