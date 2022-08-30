package main

import (
	"context"
	"fmt"
	"net"
	"time"
)

// maxBufferSize specifies the size of the buffers that
// are used to temporarily hold data from the UDP packets
// that we receive.
const maxBufSize = 1024
const timeout = 15 * time.Second

func server(ctx context.Context, addr string) error {
	doneCh := make(chan error, 1)
	buf := make([]byte, maxBufSize)

	pc, err := net.ListenPacket("udp", addr)
	if err != nil {
		return err
	}

	defer pc.Close()

	go func() {
		for {
			// Read from client.
			n, a, err := pc.ReadFrom(buf)
			if err != nil {
				doneCh <- err
				return
			}

			fmt.Printf("packet-received: bytes=%d from=%s\n", n, a.String())

			err = pc.SetWriteDeadline(time.Now().Add(timeout))
			if err != nil {
				doneCh <- err
				return
			}

			// Write the packet's content back to the client.
			n, err = pc.WriteTo(buf[:n], a)
			if err != nil {
				doneCh <- err
				return
			}

			fmt.Printf("packet-written: bytes=%d to=%s\n", n, a.String())
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("\tServer says: bye!")
		err = ctx.Err()
	case <-doneCh:
		return err
	}

	return nil
}

func client(ctx context.Context, addr string) error {
	doneCh := make(chan error, 1)

	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}

	defer conn.Close()

	//go func() {
	//
	//}()
}
