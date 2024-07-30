package spacelink

import (
	"context"
	"crypto/tls"

	"spacelink/internal/session"
	"spacelink/utils"

	"github.com/quic-go/quic-go"
)

type Client struct {
	hostname string
	config   *quic.Config
	tlsConf  *tls.Config
	sess     session.Session
}

/**/
func (*Client) SendDataStream(dataBuffer []byte) (int, error) {
	return len(dataBuffer), nil
}

/**/

func NewClient(config *quic.Config, serAddr string) (Client, error) {
	var clt Client
	clt.hostname = serAddr
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	clt.tlsConf = &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{utils.ALPN},
	}
	conn, err := quic.DialAddr(ctx, serAddr, clt.tlsConf, config)
	if err != nil {
		return clt, err
	}
	clt.sess = session.NewSession(conn)
	return clt, nil
}

func (clt *Client) SendData(dataBuffer []byte) (int, error) {
	return clt.sess.WriteData(dataBuffer)
}
