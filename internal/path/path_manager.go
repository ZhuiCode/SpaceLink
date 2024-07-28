package path

import (
	"net"

	"github.com/quic-go/quic-go"

	"github.com/quic-go/quic-go/internal/utils"
)

type pathManager struct {
	nxtPathID PathID
	// Number of paths, excluding the initial one
	nbPaths uint8

	remoteAddrs4 []net.UDPAddr
	remoteAddrs6 []net.UDPAddr
}

func (pm *pathManager) setup(conn quic.Connection) {
	// Initial PathID is 0
	// PathIDs of client-initiated paths are even
	// those of server-initiated paths odd
	/*
		if pm.sess.perspective == protocol.PerspectiveClient {
			pm.nxtPathID = 1
		} else {
			pm.nxtPathID = 2
		}*/

	pm.remoteAddrs4 = make([]net.UDPAddr, 0)
	pm.remoteAddrs6 = make([]net.UDPAddr, 0)

	pm.nbPaths = 0

	// With the initial path, get the remoteAddr to create paths accordingly
	if conn.RemoteAddr() != nil {
		remAddr, err := net.ResolveUDPAddr("udp", conn.RemoteAddr().String())
		if err != nil {
			utils.Errorf("path manager: encountered error while parsing remote addr: %v", remAddr)
		}

		if remAddr.IP.To4() != nil {
			pm.remoteAddrs4 = append(pm.remoteAddrs4, *remAddr)
		} else {
			pm.remoteAddrs6 = append(pm.remoteAddrs6, *remAddr)
		}
	}
}

func (pm *pathManager) run() {
	// Close immediately if requested

}

func getIPVersion(ip net.IP) int {
	if ip.To4() != nil {
		return 4
	}
	return 6
}

func (pm *pathManager) createPath(locAddr net.UDPAddr, remAddr net.UDPAddr) error {
	// First check that the path does not exist yet
	return nil
}

func (pm *pathManager) createPaths() error {

	return nil
}

func (pm *pathManager) createPathFromRemote(p *receivedPacket) (*path, error) {

	pth := &path{
		pathID: pathID,
		sess:   pm.sess,
		conn:   &conn{pconn: localPconn, currentAddr: remoteAddr},
	}

	return pth, nil
}

func (pm *pathManager) closePath(pthID PathID) error {

	return nil
}

func (pm *pathManager) closePaths() {
	return
}
