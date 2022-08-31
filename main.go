package main

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// maxBufferSize specifies the size of the buffers that
// are used to temporarily hold data from the UDP packets
// that we receive.
const (
	maxBufSize = 1024
	timeout    = 15 * time.Second
	host       = "127.0.0.1"
	port       = 1337
)

var isServer = flag.Bool("s", false, "whether it should be run as a server")

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
	buf := make([]byte, maxBufSize)

	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		return err
	}

	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		return err
	}

	defer conn.Close()

	go func() {
		for {
			r := bufio.NewReader(os.Stdin)
			str, err := r.ReadString('\n')
			n, err := conn.Write([]byte(str))
			if err != nil {
				doneCh <- err
				return
			}

			fmt.Printf("packet-written: bytes=%d\n", n)

			deadline := time.Now().Add(timeout)
			err = conn.SetReadDeadline(deadline)
			if err != nil {
				doneCh <- err
				return
			}

			nRead, addr, err := conn.ReadFrom(buf)
			if err != nil {
				doneCh <- err
				return
			}

			fmt.Printf("packet-received: bytes=%d from=%s\n", nRead, addr.String())
		}
	}()

	select {
	case <-ctx.Done():
		fmt.Println("cancelled")
		err = ctx.Err()
	case err = <-doneCh:
		return err
	}

	return err
}

func main() {
	flag.Parse()

	var (
		err     error
		address = fmt.Sprintf("%s:%d", host, port)
	)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		sigChan := make(chan os.Signal)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		cancel()
	}()

	if *isServer {
		fmt.Println("running as a server on " + address)
		err = server(ctx, address)
		if err != nil && err != context.Canceled {
			panic(err)
		}
		return
	}

	fmt.Println("sending to " + address)
	err = client(ctx, address)
	if err != nil && err != context.Canceled {
		panic(err)
	}
}
