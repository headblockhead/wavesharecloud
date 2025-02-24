package main

import (
	"fmt"
	"image"
	"image/jpeg"
	"net"
	"os"

	"github.com/headblockhead/wavesharecloud"
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
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	defer l.Close()
	
	for {
		// Wait for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting a connection: %v", err)
			os.Exit(1)
		}
		// Handle the connection in a new goroutine.
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	fmt.Println("New connection from:", conn.RemoteAddr())
	lc := wavesharecloud.NewLoggingConn(conn, false)

	// If a password is required to unlock the display, it should be provided as the second argument, and the display will be unlocked.
	// In this case, the display is not locked, so the password is not required.
	display := wavesharecloud.NewDisplay(lc, "")

	// Open and JPEG decode the test image into an 'image.Image'.
	testPatternImage, err := openImage("testpattern.jpg")
	if err != nil {
		fmt.Printf("Error opening test pattern image: %v\n", err)
	}

	// Draw the testpattern image to the display.
	err = display.SendImage(testPatternImage)
	if err != nil {
		fmt.Printf("Error sending testpattern image: %v\n", err)
	}

	err := display.Shutdown()
	if err != nil {
		fmt.Printf("Error shutting down: %v\n", err)
	}
	display.Disconnect()
}

func openImage(path string) (img image.Image, err error) {
	imageFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer imageFile.Close()
	img, err = jpeg.Decode(imageFile)
	if err != nil {
		return nil, err
	}
	return img, nil
}
