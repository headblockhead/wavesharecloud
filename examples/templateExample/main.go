package main

import (
	"fmt"
	"net"
	"os"

	"github.com/headblockhead/waveshareCloud"
)

const (
	CONN_HOST = "0.0.0.0"
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
	fmt.Println("New connection from:", conn.RemoteAddr())
	// Setting up the connection to the display.
	lc := waveshareCloud.NewLoggingConn(conn, false)

	// Creating the display. If a password is required to unlock the display, here is where you would enter it.
	// This automatically unlocks the display when created.
	// In this case, the display is not locked, so the password is not required.
	display := waveshareCloud.NewDisplay(lc, "")

	// Your display code goes here.
	var err error // Remove this line when you have your code. It is just here to prevent the compiler from complaining.

	// Shutdown the display.
	err = display.Shutdown()
	if err != nil {
		fmt.Printf("Error shutting down: %v\n", err)
	}
	// Close the connection.
	display.Disconnect()
}
