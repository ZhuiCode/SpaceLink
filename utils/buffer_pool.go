package utils

import (
	"sync"
)

var bufferPool sync.Pool

func GetPacketBuffer() []byte {
	return bufferPool.Get().([]byte)
}

func PutPacketBuffer(buf []byte) {
	if cap(buf) != int(MaxPacketBufferSize) {
		panic("putPacketBuffer called with packet of wrong size!")
	}
	bufferPool.Put(buf[:0])
}

func Init() {
	bufferPool.New = func() interface{} {
		return make([]byte, 0, MaxPacketBufferSize)
	}
}
