package main

import (
	"fmt"
	"image/jpeg"
	"net"
	"os"
	"time"

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
	fmt.Println("New connection from:", conn.RemoteAddr())
	// Setting up the connection to the display.
	lc := waveshareCloud.NewLoggingConn(conn, false)
	// Creating the representation of the display. It is not locked at the moment, so this boolean is false.
	display := waveshareCloud.NewDisplay(lc, false)

	file, err := os.Open("image.jpg")
	if err != nil {
		fmt.Printf("Error opening image: %v\n", err)
	}
	defer file.Close()
	img, err := jpeg.Decode(file)
	if err != nil {
		fmt.Printf("Error decoding image: %v\n", err)
	}

	// We do not want to scale the image to the display size, so this boolean is false.
	// This will crop the image to the display size from the top left corner.
	err = display.SendImage(img, false)
	if err != nil {
		fmt.Printf("Error sending image: %v\n", err)
	}

	// Wait 3 seconds before sending the next image. This gives the user time to see the image.
	time.Sleep(time.Second * 3)

	// This time, we do want to scale the image to the display size, so this boolean is true.
	// This will scale the image by squashing it to the display size.
	err = display.SendImage(img, true)
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
