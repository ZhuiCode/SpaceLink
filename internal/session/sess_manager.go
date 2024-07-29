package session

import (
	"net"
)

type SessManager struct {
	// Number of session paths, excluding the initial one
	nbPaths uint8
	sess    Session
}

func (sm *SessManager) setup() {
	sm.nbPaths = 0

}

func (pm *SessManager) run() {
	// Close immediately if requested

}

func getIPVersion(ip net.IP) int {
	if ip.To4() != nil {
		return 4
	}
	return 6
}

func (pm *SessManager) createPath(locAddr net.UDPAddr, remAddr net.UDPAddr) error {
	// First check that the path does not exist yet
	return nil
}

func (pm *SessManager) createPaths() error {

	return nil
}

func (pm *SessManager) createPathFromRemote() error {
	/*
		pth := &path{
			pathID: pathID,
			sess:   pm.sess,
			conn:   &conn{pconn: localPconn, currentAddr: remoteAddr},
		}
	*/
	return nil
}

func (sm *SessManager) closeSession() error {
	sm.sess.Close()
	return nil
}
