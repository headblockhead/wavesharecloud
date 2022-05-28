package main

import (
	"fmt"
	"image"
	"image/color"
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
	err = display.SendImage(generateblack())
	if err != nil {
		fmt.Printf("Error sending image: %v\n", err)
	}
	// Shutdown the connection.
	err = display.Shutdown()
	if err != nil {
		fmt.Printf("Error shutting down: %v\n", err)
	}
	display.Disconnect()
}

func generateblack() (data []byte) {
	data = make([]byte, 400*300/8)
	for i := 0; i < (400 * 300 / 8); i++ {
		data[i] = 0x00
	}
	return data
}

func generateCheckerboard() (data []byte) {
	img := image.NewRGBA(image.Rect(0, 0, 400, 300))
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			img.Set(x, y, color.RGBA{0, 0, 0, 255})
			// if x%2 == 0 {
			// 	img.Set(x, y, color.White)
			// } else {
			// 	img.Set(x, y, color.Black)
			// }
		}
	}
	return convertImageToBits(img)
}

func convertImageToBits(img image.Image) []byte {
	wh := img.Bounds()
	b := make([]byte, (wh.Max.X*wh.Max.Y)/8)
	for y := 0; y < wh.Max.Y; y++ {
		for x := 0; x < wh.Max.X; x++ {
			if img.At(x, y) == color.Black {
				continue
			}
			byteIndex := (y * wh.Max.X / 8) + (x / 8)
			bitIndex := x % 8
			b[byteIndex] |= (1 << bitIndex)
		}
	}
	return b
}
