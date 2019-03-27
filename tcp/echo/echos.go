package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

// This program implements a simple echo server over TCP.
// When the server receives a request, it returns its content
// immediately.
//
// Usage:
// echos -e <host:address>
func main() {
	var addr string
	// bind to port :4040 an all local interface
	flag.StringVar(&addr, "e", ":4040", "service address endpoint")
	flag.Parse()

	// create local address for socket
	laddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// #1 announce service using ListenTCP
	// which a TCPListener
	// this call create a listener for client
	l, err := net.ListenTCP("tcp", laddr)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	//http://rosettacode.org/wiki/Exceptions#Go
	//http://rosettacode.org/wiki/Exceptions#Python
	defer l.Close()
	fmt.Println("listening at (tcp)", laddr.String())

	// #2 req/response loop
	for {
		// AcceptTCP on that listener
		// use TCPListener to block and waiting for request from the client
		// connection request using AcceptTCP which creates a TCPConn
		// we have conn value which we can use to subsequence interaction to the client
		// conn is type TCPConn
		conn, err := l.AcceptTCP()
		if err != nil {
			fmt.Println("failed to accept conn:", err)
			conn.Close()
			continue
		}
		fmt.Println("connected to: ", conn.RemoteAddr())

		// pass coonn to the new goroutine
		// goroutine allows the listener loop to contiue and come backs to l.AcceptTCP()
		// blok and waiting
		// goroutine https://gobyexample.com/goroutines
		// but the same time... #3
		go handleConnection(conn)
	}
}

// #3 ... handleConnection will kick off
// handleConnection reads request from connection with conn.Read()
// then write response using conn.Write()
// then the connection is closed
func handleConnection(conn *net.TCPConn) {
	// the first make sure when handleConnection is down
	// this connection will close properly
	defer conn.Close() // clean up when done

	// set up a size buffer
	buf := make([]byte, 1024)

	// #4 use conn.Read to block and wait on values from the client
	// read something form the client err or the n (number of bytes)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}

	// #5 user conn.Write method send back to the client the content we read
	// from 0 to n bytes
	w, err := conn.Write(buf[:n])
	if err != nil {
		fmt.Println("failed to write to client:", err)
		return
	}
	// silly check to make sure the number of bytes was written we read form the client
	if w != n { // was all data sent
		fmt.Println("warning: not all data sent to client")
		return
	}
}

// go run echos.go -h
// go run echos.go

// nc localhost 4040
