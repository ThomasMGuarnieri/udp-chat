package main

import (
	"bufio"
	"flag"
	"fmt"
	"net"
	"os"
)

const (
	bufSize = 1024
	ip      = "127.0.0.1"
	port    = 4444
)

var isServer = flag.Bool("s", false, "whether it should be run as a server")

func server() error {
	p := make([]byte, bufSize)
	addr := net.UDPAddr{
		Port: port,
		IP:   net.ParseIP(ip),
	}

	conn, err := net.ListenUDP("udp", &addr)
	if err != nil {
		return err
	}
	for {
		_, remoteAddr, err := conn.ReadFromUDP(p)
		fmt.Printf("Read a message from %v %s \n", remoteAddr, p)
		if err != nil {
			fmt.Printf("Some error  %v", err)
			continue
		}
		go sendResponse(conn, remoteAddr)
	}
}

func client() error {
	p := make([]byte, bufSize)

	conn, err := net.Dial("udp", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		fmt.Printf("Some error %v", err)
		return nil
	}

	defer conn.Close()

	// scan the stdin for client messages
	scn := bufio.NewScanner(os.Stdin)
	for scn.Scan() {
		_, err = fmt.Fprintf(conn, "%s", scn.Bytes())
		if err != nil {
			return err
		}

		// Read the answer from server
		n, err := bufio.NewReader(conn).Read(p)
		if err != nil {
			return err
		}

		fmt.Printf("%s\n", p[:n])
	}

	return nil
}

//func broadcast(b []byte, as []net.Addr, conn net.PacketConn) error {
//	for _, a := range as {
//		fmt.Println("Broadcast: ", a)
//		_, err := conn.WriteTo(b, a)
//		if err != nil {
//			return err
//		}
//	}
//
//	return nil
//}

func sendResponse(conn *net.UDPConn, addr *net.UDPAddr) {
	_, err := conn.WriteToUDP([]byte("From server: Hello I got your message "), addr)
	if err != nil {
		fmt.Printf("Couldn't send response %v", err)
	}
}

func main() {
	flag.Parse()

	if *isServer {
		fmt.Println("running as a server on ")
		err := server()
		if err != nil {
			fmt.Println(err)
		}
	}

	fmt.Println("sending to ")
	err := client()
	if err != nil {
		fmt.Println(err)
	}

}
