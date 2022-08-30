package main

import (
	"context"
	"log"
	"net"
)

// maxBufferSize specifies the size of the buffers that
// are used to temporarily hold data from the UDP packets
// that we receive.
const maxBufSize = 1024

func server(ctx context.Context, addr string) error {
	doneCh := make(chan error, 1)
	buf := make([]byte, maxBufSize)

	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		log.Panicln(err)
	}

	defer pc.Close()

	//go func() {
	//	for {
	//
	//	}
	//}()
}
