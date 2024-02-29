package network

import (
	"flag"
	"fmt"
	"net"
	"os"
)

func main1() {
	var addr string
	flag.StringVar(&addr, "e", ":4040", "service address endpoint")
	flag.Parse()

	// create local addr for socket
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println(err)
	}

	// announce service using ListenTCP
	listener, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	defer listener.Close()
	fmt.Println("listening at (tcp)", laddr.String())

	for i := 0; i < 2; i++ {
		go handleAccept(listener)
	}

}

func handleAccept(listener *net.TCPListener) {
	// req/response loop
	for {
		// use TCPListener to block and wait for TCP
		// connection request using AcceptTCP which creates a TCPConn
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println("failed to accept conn:", err)
			conn.Close()
			continue
		}

		fmt.Println("connected to: ", conn.RemoteAddr())

		go handleConnection(conn)
	}
}

func handleConnection(conn *net.TCPConn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	w, err := conn.Write(buf[:n])
	if err != nil {
		fmt.Println("failed to write to client: ", err)
		return
	}

	if w != n {
		fmt.Println("Warning: not all data sent to client")
		return
	}
}
