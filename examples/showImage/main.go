package main

import (
	"fmt"
	"image"
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

	// Creating the display. If a password is required to unlock the display, here is where you would enter it.
	// This automatically unlocks the display when created.
	// In this case, the display is not locked, so the password is not required.
	display := waveshareCloud.NewDisplay(lc, "")

	// Open the timages and decode them.
	testPatternImage, err := openImage("testpattern.jpg")
	if err != nil {
		fmt.Printf("Error opening test pattern image: %v\n", err)
	}
	flowerImage, err := openImage("flowers.jpg")
	if err != nil {
		fmt.Printf("Error opening flower image: %v\n", err)
	}

	// Drawing the testpattern image to the display.
	err = display.SendImage(testPatternImage)
	if err != nil {
		fmt.Printf("Error sending testpattern image: %v\n", err)
	}

	// Wait 3 seconds before sending the next image. This gives the user time to see the image.
	time.Sleep(time.Second * 3)

	// Drawing the flower image to the display - Cropped from the top left
	err = display.SendImage(flowerImage)
	if err != nil {
		fmt.Printf("Error sending cropped flower image: %v\n", err)
	}

	// Wait 3 seconds before sending the next image. This gives the user time to see the image.
	time.Sleep(time.Second * 3)

	// Drawing the flower image to the display - Scaled to fill the screen
	err = display.SendImageScaled(flowerImage)
	if err != nil {
		fmt.Printf("Error sending scaled flower image: %v\n", err)
	}

	// Shutdown the display.
	err = display.Shutdown()
	if err != nil {
		fmt.Printf("Error shutting down: %v\n", err)
	}
	// Close the connection.
	display.Disconnect()
}

func openImage(path string) (img image.Image, err error) {
	// Open the test pattern image and decode it.
	testPatternFile, err := os.Open("path")
	if err != nil {
		return nil, err
	}
	defer testPatternFile.Close()
	testPatternImage, err := jpeg.Decode(testPatternFile)
	if err != nil {
		return nil, err
	}
	return testPatternImage, nil
}
