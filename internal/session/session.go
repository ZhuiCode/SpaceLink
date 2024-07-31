package session

import (
	"fmt"
	"net"
	linkerr "spacelink/error"
	"spacelink/utils"
	"time"

	"github.com/quic-go/quic-go"
)

type ReceivedPacket struct {
	remoteAddr net.Addr
	data       []byte
	rcvTime    time.Time
	rcvPconn   net.PacketConn
}

type path struct {
	pathID quic.StreamID //目前是基于streamID生成的，后续需要进一步考虑调整
	stream quic.Stream
}

type Session struct {
	// Number of session paths, excluding the initial one
	nbpaths   uint8
	nxtpathID int64

	perspective             utils.Perspective
	sessionCreationTime     time.Time
	lastNetworkActivityTime time.Time

	paths        map[quic.StreamID]path
	remoteAddrs4 []net.UDPAddr
	remoteAddrs6 []net.UDPAddr
	conn         quic.Connection
}

func NewSession(conn quic.Connection) Session {
	// Initial pathID is 0
	// pathIDs of client-initiated paths are even
	// those of server-initiated paths odd
	var sess Session
	if sess.perspective == utils.PerspectiveClient {
		sess.nxtpathID = 1
	} else {
		sess.nxtpathID = 2
	}

	sess.remoteAddrs4 = make([]net.UDPAddr, 0)
	sess.remoteAddrs6 = make([]net.UDPAddr, 0)
	// With the initial path, get the remoteAddr to create paths accordingly
	if conn.RemoteAddr() != nil {
		remAddr, err := net.ResolveUDPAddr("udp", conn.RemoteAddr().String())
		if err != nil {
			utils.DefaultLogger.Errorf("path manager: encountered error while parsing remote addr: %v", remAddr)
		}

		if remAddr.IP.To4() != nil {
			sess.remoteAddrs4 = append(sess.remoteAddrs4, *remAddr)
		} else {
			sess.remoteAddrs6 = append(sess.remoteAddrs6, *remAddr)
		}
	}
	sess.conn = conn
	return sess
}

func (sess *Session) GetPerspective() utils.Perspective {
	return sess.perspective
}
func (sess *Session) AddpathToSess(pth path) error {
	_, ok := sess.paths[pth.pathID]
	if !ok {
		sess.paths[pth.pathID] = pth
	} else {
		utils.DefaultLogger.Errorf("path ", pth.pathID, "has been found")
		return linkerr.ErrList[0]
	}
	return nil
}
func (sess *Session) DelpathFromSess(pathID quic.StreamID) error {
	path, ok := sess.paths[pathID]
	if !ok {
		return linkerr.ErrList[1]
	} else {
		path.stream.Close()
		delete(sess.paths, pathID)
	}
	return nil
}

func (sess *Session) Close() error {
	for id, path := range sess.paths {
		path.stream.Close()
		delete(sess.paths, id)
	}
	return nil
}

func (sess *Session) Createpath() error {
	str, err := sess.conn.OpenStream()
	if err != nil {
		return err
	}
	pth := path{stream: str}
	sess.paths[str.StreamID()] = pth
	fmt.Println("path id is ", str.StreamID())
	return nil
}

func (sess *Session) WriteData(dataBuf []byte) (int, error) {
	for _, pth := range sess.paths {
		return pth.stream.Write(dataBuf)
	}
	return 0, linkerr.ErrList[2]
}

func (sess *Session) ReadData(dataBuf []byte) (int, error) {
	for _, pth := range sess.paths {
		return pth.stream.Read(dataBuf)
	}
	return 0, linkerr.ErrList[2]
}
