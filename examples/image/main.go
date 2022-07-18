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

	// Open the test pattern image and decode it.
	testPatternFile, err := os.Open("testpattern.jpg")
	if err != nil {
		fmt.Printf("Error opening testpattern image: %v\n", err)
	}
	defer testPatternFile.Close()
	testPatternImage, err := jpeg.Decode(testPatternFile)
	if err != nil {
		fmt.Printf("Error decoding testpattern image: %v\n", err)
	}

	// Open the flower image and decode it.
	flowerFile, err := os.Open("flowers.jpg")
	if err != nil {
		fmt.Printf("Error opening flowers image: %v\n", err)
	}
	defer flowerFile.Close()
	flowerImage, err := jpeg.Decode(flowerFile)
	if err != nil {
		fmt.Printf("Error decoding flowers image: %v\n", err)
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
