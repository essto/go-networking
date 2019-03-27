package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

// This program implements a simple echo client over tcp
// It sends a text content to the server and displays
// the response on the screen
//
// Usage: echoc -e <host-addr-endpoint> <text content>
func main() {
	var addr string
	// set up default flag
	flag.StringVar(&addr, "e", "localhost:4040", "service address endpoint")
	flag.Parse()
	text := flag.Arg(0)

	// call ResolveTCPAddr to create address that is bound to this value (passed by flag "e")
	raddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// next calll DialTCP to create a connection to the remote address
	// !!! there is no need to specify the local address
	// return the conn of typ TCPconn
	conn, err := net.DialTCP("tcp", nil, raddr)
	if err != nil {
		fmt.Println("failed to connect to server:", err)
		os.Exit(1)
	}
	defer conn.Close()

	// send text to server
	// takes as value text we get from the command line
	// text converts to slice of bytes an send it to the host
	_, err = conn.Write([]byte(text))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// read response
	// set up a buf of size
	buf := make([]byte, 1024)
	// n - bytes we read from the remote host
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println("failed reading response:", err)
		os.Exit(1)
	}
	fmt.Println(string(buf[:n]))
}

// go run echos.go "hello"

// go run echos.go
// hello
