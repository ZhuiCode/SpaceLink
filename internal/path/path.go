package path

import (
	"time"

	"github.com/quic-go/quic-go"
)

type PathID int32

const (
	minpathTimer = 10 * time.Millisecond
	// XXX (QDC): To avoid idling...
	maxpathTimer = 1 * time.Second
)

type path struct {
	pathID pathID
	conn   quic.Connection
}

// setup initializes values that are independent of the perspective
func (p *path) setup() {
	return
}

func (p *path) close() error {

	return nil
}

func (p *path) run() {

}
