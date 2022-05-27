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
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		fmt.Println("Accepted connection")
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	fmt.Println("handle connection")
	// Get device ID
	display := waveshareCloud.Display{
		connection: conn,
	}
	display.Send("G")
	command, data, err := display.Receive()
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	fmt.Println("Received:", data, "From:", command)
	// Shutdown the connection.
	display.Send("S")
	conn.Close()
}
