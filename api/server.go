package spacelink

import (
	"context"
	"crypto/ed25519"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"io"
	"log"
	"math/big"
	"net"
	"spacelink/internal/session"
	"spacelink/utils"
	"time"

	"github.com/quic-go/quic-go"
)

/*
type packetHandler interface {
	quic.Connection
	handlePacket(*session.ReceivedPacket)
	GetVersion() quic.VersionNumber
	run() error
	closeRemote(error)
}*/

type Server struct {
	hostname string
	config   *quic.Config
	tlsConf  *tls.Config
}

func generateSelfSignedTLSConfig() (*tls.Config, error) {
	pubKey, privKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, err
	}

	notBefore := time.Now()
	notAfter := notBefore.Add(365 * 24 * time.Hour)

	template := &x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			CommonName: "localhost",
		},
		NotBefore:             notBefore,
		NotAfter:              notAfter,
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
		IsCA:                  false,
		IPAddresses:           []net.IP{net.ParseIP("127.0.0.1")},
	}

	derBytes, err := x509.CreateCertificate(rand.Reader, template, template, pubKey, privKey)
	if err != nil {
		return nil, err
	}

	certPEM := pem.EncodeToMemory(&pem.Block{Type: "CERTIFICATE", Bytes: derBytes})
	b, err := x509.MarshalPKCS8PrivateKey(privKey)
	if err != nil {
		return nil, err
	}
	keyPEM := pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: b})

	cert, err := tls.X509KeyPair(certPEM, keyPEM)
	if err != nil {
		return nil, err
	}

	return &tls.Config{
		Certificates: []tls.Certificate{cert},
	}, nil
}
func NewServer(hostname string) {
	serConfig := quic.Config{
		MaxConnectionReceiveWindow: 1 << 30,
		MaxStreamReceiveWindow:     1 << 30,
	}
	ser := Server{hostname: "localhost",
		config: &serConfig,
	}
	tlsConf, err := generateSelfSignedTLSConfig()
	if err != nil {
		ser.tlsConf = tlsConf
	}
	session.NewSession()
}

func (s *Server) RunServListen() error {
	tlsConf, err := generateSelfSignedTLSConfig()
	if err != nil {
		log.Fatal(err)
	}
	tlsConf.NextProtos = []string{utils.ALPN}
	//tlsConf.KeyLogWriter = keyLogFile

	conf := s.config.Clone()
	ln, err := quic.ListenAddr(s.hostname, tlsConf, conf)
	if err != nil {
		return err
	}
	log.Println("Listening on", ln.Addr())
	defer ln.Close()
	for {
		conn, err := ln.Accept(context.Background())
		if err != nil {
			return fmt.Errorf("accept errored: %w", err)
		}
		go s.HandleServConn(conn)
	}
}

func (s *Server) HandleServConn(conn quic.Connection) {
	for {
		stm, err := conn.AcceptStream(context.Background())
		if err != nil {
			return
		}
		log.Println("AcceptStream ", stm.StreamID())
		go s.handleServerStream(stm)
	}
}

func (s *Server) handleServerStream(stm quic.Stream) {

	var dataBuffer []byte
	// receive data until the client sends a FIN
	for {
		if _, err := stm.Read(dataBuffer); err != nil {
			if err == io.EOF {
				break
			}
		}
	}
}
