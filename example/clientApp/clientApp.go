package main

import (
	"fmt"
	spacelink "spacelink/api"

	"github.com/quic-go/quic-go"
)

func main() {
	config := quic.Config{
		MaxConnectionReceiveWindow: 1 << 30,
		MaxStreamReceiveWindow:     1 << 30,
	}
	serAddr := "127.0.0.1:6061"
	client, err := spacelink.NewClient(&config, serAddr)
	if err == nil {
		dataBuffer := make([]byte, 1024*16)
		client.SendData(dataBuffer)
	} else {
		fmt.Println("NewClient failed", err.Error())
	}
}
