package utils

import (
	"github.com/quic-go/quic-go"
)

/*axPacketBufferSize maximum packet size of any QUIC packet, based on ethernet's max size, minus the IP and UDP headers. IPv6 has a 40 byte header, UDP adds an additional 8 bytes. This is a total overhead of 48 bytes. Ethernet's max packet size is 1500 bytes, 1500 - 48 = 1452.*/
const MaxPacketBufferSize = 1452

type Perspective int

// Perspective determines if we're acting as a server or a client
const (
	PerspectiveServer Perspective = 1
	PerspectiveClient Perspective = 2
)
const ALPN = "perf"

var QConfig = &quic.Config{
	// use massive flow control windows here to make sure that flow control is not the limiting factor
	MaxConnectionReceiveWindow: 1 << 30,
	MaxStreamReceiveWindow:     1 << 30,
}
