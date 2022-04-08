package leaderelection

import "fmt"

type LE interface {
	AmILeader() bool
}

type Base string

const (
	NATS Base = "NATS"
)

func NewLE(base Base, name string, size int) (LE, error) {
	switch base {
	case NATS:
		return NewNatsLE(name, size)
	default:
		err := fmt.Sprintf("error creating new leaderElector: unsupported base: %v", base)
		return nil, fmt.Errorf(err)
	}
}
