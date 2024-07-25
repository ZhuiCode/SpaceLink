package buffer

import (
	"sync"

	"github.com/quic-go/quic-go/internal/protocol"
)

var bufferPool sync.Pool

func GetPacketBuffer() []byte {
	return bufferPool.Get().([]byte)
}

func PutPacketBuffer(buf []byte) {
	if cap(buf) != int(protocol.MaxPacketBufferSize) {
		panic("putPacketBuffer called with packet of wrong size!")
	}
	bufferPool.Put(buf[:0])
}

func Init() {
	bufferPool.New = func() interface{} {
		return make([]byte, 0, protocol.MaxPacketBufferSize)
	}
}
