package spacelink

import (
	"crypto/tls"

	"github.com/quic-go/quic-go"
)

type Client struct {
	hostname string
	config   *quic.Config
	tlsConf  *tls.Config
	streanms map[string]quic.Stream
	conn     quic.Connection
}

/**/
func (*Client) SendDataStream(dataBuffer []byte) (int, error) {
	return len(dataBuffer), nil
}

/**/
