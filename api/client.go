package spacelink

import (
	"crypto/tls"

	"github.com/quic-go/quic-go"
)

type client struct {
	hostname string
	config   *quic.Config
	tlsConf  *tls.Config
	streanms map[string]quic.Stream
	conn     quic.Connection
}

/**/
func (*client) SendDataStream(dataBuffer []byte) (int, error) {
	return len(dataBuffer), nil
}

/**/
