package schduler

import (
	"github.com/quic-go/quic-go"
)

type Scheduler struct {
	// XXX Currently round-robin based, inspired from MPTCP scheduler
	quotas map[quic.StreamID]uint
}

func (sch *Scheduler) setup() {
	sch.quotas = make(map[quic.StreamID]uint)
}
