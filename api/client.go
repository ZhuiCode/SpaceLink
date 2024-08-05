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
	ctx      context.Context
	cancel   context.CancelFunc
}

/**/

func NewClient(config *quic.Config, serAddr string) (Client, error) {
	var clt Client
	clt.hostname = serAddr
	clt.ctx, clt.cancel = context.WithCancel(context.Background())
	clt.tlsConf = &tls.Config{
		InsecureSkipVerify: true,
		NextProtos:         []string{utils.ALPN},
	}

	clt.sess = session.NewSession(clt.ctx, config, serAddr, clt.tlsConf)
	return clt, nil
}

func (clt *Client) SendData(dataBuffer []byte) (int, error) {
	return clt.sess.WriteData(dataBuffer)
}
