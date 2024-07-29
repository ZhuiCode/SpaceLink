package path

import (
	"time"

	"github.com/quic-go/quic-go"
)

const (
	minpathTimer = 10 * time.Millisecond
	// XXX (QDC): To avoid idling...
	maxpathTimer = 1 * time.Second
)

type Path struct {
	PathID        int64 //目前是基于streamID生成的，后续需要进一步考虑调整
	stream        quic.Stream
	readDeadline  time.Time
	writeDeadline time.Time
}

// setup initializes values that are independent of the perspective
func (pth *Path) Setup(stm quic.Stream) {
	pth.PathID = int64(stm.StreamID())
	pth.stream = stm
}

func (pth *Path) Close() {

	pth.PathID = 0
	pth.stream = nil
}

func (p *Path) run() {

}
