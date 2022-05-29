package main

import (
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"math/rand"
	"net"
	"os"

	"github.com/MaxHalford/halfgone"
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
	_, err = loadImage("/home/headb/Downloads/400x300.jpg")
	if err != nil {
		fmt.Printf("Error loading image: %v\n", err)
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
		// if (i % 3) == 0 {
		// 	data[i] = 0x81
		// } else if (i % 2) == 0 {
		// 	data[i] = 0x7E
		// } else {
		// 	data[i] = 0xEF
		// }
		data[i] = byte(rand.Intn(0xFF))
	}
	return data
}

func loadImage(path string) (data []byte, err error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	img, err := jpeg.Decode(file)
	if err != nil {
		return nil, err
	}
	gray := halfgone.FloydSteinbergDitherer{}.Apply(halfgone.ImageToGray(img))

	w, _ := os.Create("/home/headb/Downloads/400x300_dithered.jpg")
	jpeg.Encode(w, gray, &jpeg.Options{Quality: 100})

	bytes := convertImageToBits(gray)
	convertedBack := convertBitsToImage(bytes, img.Bounds())
	w2, _ := os.Create("/home/headb/Downloads/400x300_dithered_converted_back.jpg")
	jpeg.Encode(w2, convertedBack, &jpeg.Options{Quality: 100})

	return bytes, nil
}

func generateCheckerboard() (data []byte) {
	img := image.NewRGBA(image.Rect(0, 0, 400, 300))
	for y := 0; y < img.Bounds().Max.Y; y++ {
		for x := 0; x < img.Bounds().Max.X; x++ {
			if x%2 == 0 && y%2 == 0 {
				img.Set(x, y, color.White)
			} else {
				img.Set(x, y, color.Black)
			}
		}
	}
	bits := convertImageToBits(img)
	fmt.Printf("%b\n", bits[0])
	return bits
}

func convertImageToBits(img image.Image) []byte {
	wh := img.Bounds()
	b := make([]byte, (wh.Max.X*wh.Max.Y)/8)
	for y := 0; y < wh.Max.Y; y++ {
		for x := 0; x < wh.Max.X; x++ {
			if r, g, b, _ := img.At(x, y).RGBA(); r == 0 && g == 0 && b == 0 {
				continue
			}
			byteIndex := (y * wh.Max.X / 8) + (x / 8)
			bitIndex := x % 8
			b[byteIndex] |= (1 << bitIndex)
		}
	}
	return b
}

func convertBitsToImage(b []byte, bounds image.Rectangle) (img *image.Gray) {
	w, h := bounds.Dx(), bounds.Dy()
	img = image.NewGray(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			byteIndex := (y * w / 8) + (x / 8)
			bitIndex := x % 8
			if hasBit(b[byteIndex], uint(bitIndex)) {
				img.Set(x, y, color.White)
				continue
			}
			img.Set(x, y, color.Black)
		}
	}
	return img
}

func hasBit(n byte, pos uint) bool {
	val := n & (1 << pos)
	return (val > 0)
}
