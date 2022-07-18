package main

import (
	"fmt"
	"image/jpeg"
	"net"
	"os"

	"github.com/headblockhead/waveshareCloud"
)

const (
	CONN_HOST = "192.168.155.5"
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
	lc := waveshareCloud.NewLoggingConn(conn, false)
	display := waveshareCloud.NewDisplay(lc, false)
	err := display.Unlock("12345")
	if err != nil {
		fmt.Printf("Error unlocking: %v\n", err)
	}
	file, err := os.Open("image.jpg")
	if err != nil {
		fmt.Printf("Error opening image: %v\n", err)
	}
	defer file.Close()
	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Printf("Error decoding image: %v\n", err)
	}
	err = display.SendImage(img)
	if err != nil {
		fmt.Printf("Error sending image: %v\n", err)
	}
	// Shutdown the display.
	err = display.Shutdown()
	if err != nil {
		fmt.Printf("Error shutting down: %v\n", err)
	}
	// Close the connection.
	display.Disconnect()
}
