package main

import (
	"fmt"
	"net"
	"os"

	"github.com/headblockhead/waveshareCloud"
)

const (
	CONN_HOST = "192.168.155.216"
	CONN_PORT = "6868"
	CONN_TYPE = "tcp"
)

func main() {
	// Listen for incoming connections.
	l, err := net.Listen(CONN_TYPE, CONN_HOST+":"+CONN_PORT)
	if err != nil {
		fmt.Printf("Error listening for connections: %v", err)
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting a connection: %v", err)
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	lc := waveshareCloud.NewLoggingConn(conn)
	display := waveshareCloud.NewDisplay(lc, true)
	err := display.Unlock("12345")
	if err != nil {
		fmt.Printf("Error unlocking: %v\n", err)
	}
	// Shutdown the connection.
	err = display.Shutdown()
	if err != nil {
		fmt.Printf("Error shutting down: %v\n", err)
	}
	display.Disconnect()
}
