package session

import (
	"fmt"
	"net"
	linkerr "spacelink/error"
	"spacelink/internal/path"
	"spacelink/utils"

	"github.com/quic-go/quic-go"
)

type Session struct {
	// Number of session paths, excluding the initial one
	nbPaths      uint8
	nxtPathID    int64
	perspective  utils.Perspective
	paths        map[quic.StreamID]path.Path
	remoteAddrs4 []net.UDPAddr
	remoteAddrs6 []net.UDPAddr
	conn         quic.Connection
}

func NewSession(conn quic.Connection) Session {
	// Initial PathID is 0
	// PathIDs of client-initiated paths are even
	// those of server-initiated paths odd
	var sess Session
	if sess.perspective == utils.PerspectiveClient {
		sess.nxtPathID = 1
	} else {
		sess.nxtPathID = 2
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
func getIPVersion(ip net.IP) int {
	if ip.To4() != nil {
		return 4
	}
	return 6
}
func (sess *Session) GetPerspective() utils.Perspective {
	return sess.perspective
}
func (sess *Session) AddPathToSess(pth path.Path) error {
	_, ok := sess.paths[pth.PathID]
	if !ok {
		sess.paths[pth.PathID] = pth
	} else {
		utils.DefaultLogger.Errorf("Path ", pth.PathID, "has been found")
		return linkerr.ErrList[0]
	}
	return nil
}
func (sess *Session) DelPathFromSess(pathID quic.StreamID) error {
	path, ok := sess.paths[pathID]
	if !ok {
		return linkerr.ErrList[1]
	} else {
		path.Close()
		delete(sess.paths, pathID)
	}
	return nil
}

func (sess *Session) Close() error {
	for id, path := range sess.paths {
		path.Close()
		delete(sess.paths, id)
	}
	return nil
}

func (sess *Session) CreatePath() error {
	str, err := sess.conn.OpenStream()
	if err != nil {
		return err
	}
	pth := path.NewPath(str)
	sess.paths[str.StreamID()] = pth
	fmt.Println("Path id is ", str.StreamID())
	return nil
}

func (sess *Session) WriteData(dataBuf []byte) (int, error) {
	for _, pth := range sess.paths {
		return pth.Write(dataBuf)
	}
	return 0, linkerr.ErrList[2]
}

func (sess *Session) ReadData(dataBuf []byte) (int, error) {
	for _, pth := range sess.paths {
		return pth.Read(dataBuf)
	}
	return 0, linkerr.ErrList[2]
}
