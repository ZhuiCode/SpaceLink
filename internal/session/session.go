package session

import (
	"context"
	"crypto/tls"
	"net"
	linkerr "spacelink/error"
	"spacelink/internal/schduler"
	"spacelink/utils"
	"time"

	"github.com/quic-go/quic-go"
)

/*
	type ReceivedPacket struct {
		remoteAddr net.Addr
		data       []byte
		rcvTime    time.Time
		rcvPconn   net.PacketConn
	}
*/
type Session struct {
	sessionCreationTime time.Time
	//lastNetworkActivityTime time.Time

	streams      map[quic.StreamID]quic.Stream
	remoteAddrs4 []net.UDPAddr
	remoteAddrs6 []net.UDPAddr
	conn         quic.Connection

	scheduler *schduler.Scheduler
}

func NewSession(ctx context.Context, config *quic.Config, serAddr string, tlsConf *tls.Config) Session {
	// Initial pathID is 0
	// pathIDs of client-initiated paths are even
	// those of server-initiated paths odd
	var sess Session
	var err error
	sess.remoteAddrs4 = make([]net.UDPAddr, 0)
	sess.remoteAddrs6 = make([]net.UDPAddr, 0)
	sess.conn, err = quic.DialAddr(ctx, serAddr, tlsConf, config)
	if err != nil {
		return sess
	}
	// With the initial path, get the remoteAddr to create paths accordingly
	if sess.conn.RemoteAddr() != nil {
		remAddr, err := net.ResolveUDPAddr("udp", sess.conn.RemoteAddr().String())
		if err != nil {
			utils.DefaultLogger.Errorf("path manager: encountered error while parsing remote addr: %v", remAddr)
		}
	}
	sess.sessionCreationTime = time.Now()
	return sess
}
func (sess *Session) Close() error {
	for id, stm := range sess.streams {
		stm.Close()
		delete(sess.streams, id)
	}
	return nil
}

func (sess *Session) WriteData(dataBuf []byte) (int, error) {
	for _, stm := range sess.streams {
		return stm.Write(dataBuf)
	}
	return 0, linkerr.ErrList[2]
}

func (sess *Session) ReadData(dataBuf []byte) (int, error) {
	for _, stm := range sess.streams {
		return stm.Read(dataBuf)
	}
	return 0, linkerr.ErrList[2]
}
