/**
* declare all available env for this repo
 */
package leaderelection

const (
	//log
	LOG_PATH = "LOG_PATH" //string. log path. default log to stdout

	//le
	//nats
	NATS_QUORUM = "NATS_QUORUM" //comma separated. default: nats://127.0.0.1:4222. reference: https://github.com/nats-io/go-nats/blob/master/example_test.go
)
