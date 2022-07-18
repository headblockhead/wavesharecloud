package waveshareCloud

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"image"
	"image/color"
	"io"
	"net"
	"strconv"
	"strings"

	"github.com/MaxHalford/halfgone"
	"github.com/disintegration/imaging"
)

// NewLoggingConn returns a new connection that will optionally log all traffic.
func NewLoggingConn(conn net.Conn, debug bool) *LoggingConn {
	return &LoggingConn{
		Connection: conn,
		Debug:      debug,
	}
}

// LoggingConn is a connection that will optionally log all traffic.
type LoggingConn struct {
	Connection net.Conn
	Debug      bool
}

// Read reads raw data from the connection.
func (lc *LoggingConn) Read(b []byte) (n int, err error) {
	n, err = lc.Connection.Read(b)
	if lc.Debug {
		fmt.Println("< " + string(b))
	}
	return n, err
}

// Write writes raw data to the connection.
func (lc *LoggingConn) Write(b []byte) (n int, err error) {
	if lc.Debug {
		fmt.Println("> " + string(b))
	}
	return lc.Connection.Write(b)
}

// Close closes the connection.
func (lc *LoggingConn) Close() error {
	return lc.Connection.Close()
}

// NewDisplay returns a new display, using a witch of 400 and a height of 300.
func NewDisplay(conn io.ReadWriteCloser, locked bool) *Display {
	return &Display{
		Connection: conn,
		unlocked:   !locked,
		Width:      400,
		Height:     300,
	}
}

// Display represents the physical display.
type Display struct {
	Connection io.ReadWriteCloser
	unlocked   bool
	Width      int
	Height     int
}

// SendCommand sends a command to the display.
func (display *Display) SendCommand(command string) (err error) {
	var check uint32
	for i := 0; i < len(command); i++ {
		check = check ^ uint32(command[i])
	}

	_, err = display.Connection.Write([]byte(";" + command + "/" + string(rune(check))))
	if err != nil {
		return err
	}
	return nil
}

// SendImageBytes displays an array of bytes on the screen.
func (display *Display) SendImageBytes(data []byte) (err error) {
	if !display.unlocked {
		return fmt.Errorf("display is locked")
	}
	if len(data) != (display.Width*display.Height)/8 {
		return fmt.Errorf("data length does not match display size")
	}
	err = display.SendCommand("F")
	if err != nil {
		return err
	}
	err = display.ReadBlindlyAndIgnore()
	if err != nil {
		return err
	}
	var count uint8
	for i := 0; i < len(data); i += 1024 {
		to := i + 1024
		if to > len(data) {
			to = len(data)
		}
		err = display.SendFrame(uint32(i), uint8((uint32(i)%4096)/1024), data[i:to])
		if err != nil {
			return err
		}
		err = display.ReadBlindlyAndIgnore()
		if err != nil {
			return err
		}
		count++
	}
	err = display.SendCloseFrame()
	if err != nil {
		return err
	}
	err = display.SendCommand("D")
	if err != nil {
		return err
	}
	err = display.ReadBlindlyAndIgnore()
	if err != nil {
		return err
	}
	return nil
}

// CloseFrame is the last frame to be sent to the display when drawing an image. It is completely empty.
var closeFrame = []byte{
	0x57,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0x00,
}

// SendCloseFrame sends the last frame to the display
func (display *Display) SendCloseFrame() (err error) {
	if !display.unlocked {
		return fmt.Errorf("display is locked")
	}
	_, err = display.Connection.Write(closeFrame)
	if err != nil {
		return err
	}
	// We do not need to know the checksum: it does not matter to displaying the image here.
	err = display.ReadBlindlyAndIgnore()
	if err != nil {
		return err
	}
	return nil
}

// SendImage converts an image into bytes, then sends it to the display. If the image is not 400x300, it will resize it or crop it, depending on the scale argument. This requires the display to be unlocked.
func (display *Display) SendImage(img image.Image, scale bool) (err error) {
	if !display.unlocked {
		return fmt.Errorf("display is locked")
	}
	var resizedImage image.Image
	if scale {
		resizedImage, err = resizeImage(img, display.Width, display.Height)
		if err != nil {
			return err
		}
	} else {
		resizedImage, err = cropImage(img, image.Rect(0, 0, 400, 300))
		if err != nil {
			return err
		}
	}
	imageBytes, err := loadImage(resizedImage)
	if err != nil {
		return err
	}
	return display.SendImageBytes(imageBytes)
}

// SendFrame sends a frame of image data to the display. This requires the display to be in data mode.
func (display *Display) SendFrame(addr uint32, num uint8, data []byte) (err error) {
	if !display.unlocked {
		return fmt.Errorf("display is locked")
	}
	if len(data) > 1024 {
		return fmt.Errorf("data too large, maximum size is 1024")
	}
	frame := new(bytes.Buffer)
	// This is the header of the frame. It includes:
	// 0x57 - data identifier
	err = binary.Write(frame, binary.BigEndian, uint8(0x57))
	if err != nil {
		return err
	}
	// 4 byte addr
	err = binary.Write(frame, binary.BigEndian, addr)
	if err != nil {
		return err
	}
	// 4 bytes len
	err = binary.Write(frame, binary.BigEndian, uint32(1024))
	if err != nil {
		return err
	}
	// 1 byte num
	err = binary.Write(frame, binary.BigEndian, num)
	if err != nil {
		return err
	}
	// This is the data of the frame. It is after the header.
	err = binary.Write(frame, binary.BigEndian, data)
	if err != nil {
		return err
	}
	// If there are bytes missing in a frame, fill them with 0xFF.
	// This should be expected on the final frame, as the avalible pixels on the screen do not divide into 1024 prefectly.
	if remaining := 1024 - len(data); remaining > 0 {
		remainder := make([]byte, remaining)
		for i := 0; i < len(remainder); i++ {
			remainder[i] = 0xFF
		}
		err = binary.Write(frame, binary.BigEndian, remainder)
		if err != nil {
			return err
		}
	}
	// Calculate the checksum of the frame
	var check byte
	payload := frame.Bytes()
	// The first byte is the 0x57 identifier, so we skip it.
	for i := 1; i < len(payload[1:])+1; i++ {
		// CheckSum8 Xor
		check ^= payload[i]
	}
	// The checksum byte is the last byte of the frame. It is stored as BigEndian.
	err = binary.Write(frame, binary.BigEndian, check)
	if err != nil {
		return err
	}
	// Write out the full frame to the display.
	display.Connection.Write(frame.Bytes())
	return nil
}

// ReceiveCommandOutput receives a previously sent command's output from the display.
func (display *Display) ReceiveCommandOutput(sentCommand string) (data string, err error) {
	// The first set of data is the previous command repeated back.
	buf := make([]byte, 64)
	_, err = display.Connection.Read(buf)
	if err != nil {
		return "", err
	}
	command := string(buf)
	command = formatTransmissionString(command)
	if !strings.Contains(command, sentCommand) {
		return "", fmt.Errorf("command mismatch: expected %s, got %v", sentCommand, command)
	}
	// The second set of data is the output from calling the command.
	buf = make([]byte, 64)
	_, err = display.Connection.Read(buf)
	if err != nil {
		return "", err
	}
	data = string(buf)

	data = formatTransmissionString(data)

	return data, nil
}

// ReadBlindlyAndIgnore mindlessly reads data from the display and does not do anything wth it.
func (display *Display) ReadBlindlyAndIgnore() (err error) {
	buf := make([]byte, 64)
	_, err = display.Connection.Read(buf)
	if err != nil {
		return err
	}
	return nil
}

// ReadBlindly reads any avalible data from the display and returns it (after formatting).
func (display *Display) ReadBlindly() (data string, err error) {
	buf := make([]byte, 64)
	_, err = display.Connection.Read(buf)
	if err != nil {
		return "", err
	}
	command := string(buf)
	command = formatTransmissionString(command)

	return command, nil
}

// Restart sends a command to the display to restart it. This requires the device to be unlocked.
func (display *Display) Restart() (err error) {
	if !display.unlocked {
		return fmt.Errorf("display is locked")
	}
	display.SendCommand("R")
	return nil
}

// GetBatteryLevel sends a command to the display to get its battery level. This requires the device to be unlocked.
func (display *Display) GetBatteryLevel() (err error) {
	if !display.unlocked {
		return fmt.Errorf("display is locked")
	}
	display.SendCommand("b")
	display.ReceiveCommandOutput("b")
	return nil
}

// Shutdown sends a command to the display to shut it down. This requires the device to be unlocked.
func (display *Display) Shutdown() (err error) {
	if !display.unlocked {
		return fmt.Errorf("display is locked")
	}
	display.SendCommand("S")
	return nil
}

// Disconnect closes the connection to the display
func (display *Display) Disconnect() {
	display.Connection.Close()
}

// GetID gets the ID of the display
func (display *Display) GetID() (ID string, err error) {
	err = display.SendCommand("G")
	if err != nil {
		return "", err
	}
	data, err := display.ReceiveCommandOutput("G")
	if err != nil {
		return "", err
	}
	return data, nil
}

// GetLocked gets whether the display is locked with a password
func (display *Display) GetLocked() (unlocked bool, err error) {
	err = display.SendCommand("C")
	if err != nil {
		return false, err
	}
	data, err := display.ReceiveCommandOutput("C")
	if err != nil {
		return false, err
	}
	bytedata := []byte(data)
	bytedata = bytes.Trim(bytedata, "\x00")
	unlocked, err = strconv.ParseBool(string(bytedata))
	if err != nil {
		return false, err
	}
	return unlocked, nil
}

// Unlock unlocks the display with a password. It will error if the display is already unlocked.
func (display *Display) Unlock(password string) (err error) {
	err = display.SendCommand("C")
	if err != nil {
		return err
	}
	data, err := display.ReceiveCommandOutput("C")
	if err != nil {
		return err
	}
	if strings.Contains(data, "1") {
		err = display.SendCommand("N" + password)
		if err != nil {
			return err
		}
		err = display.ReadBlindlyAndIgnore()
		if err != nil {
			return err
		}
		data, err = display.ReadBlindly()
		if err != nil {
			return err
		}
		if strings.Contains(data, "1") {
			display.unlocked = true
			return nil
		} else {
			return fmt.Errorf("wrong password")
		}
	} else {
		display.unlocked = true
		return fmt.Errorf("display is already unlocked")
	}
}

// convertImageToBits converts a black and white image to a list of bytes.
func convertImageToBits(img image.Image) []byte {
	wh := img.Bounds()
	b := make([]byte, (wh.Max.X*wh.Max.Y)/8)
	for y := 0; y < wh.Max.Y; y++ {
		for x := 0; x < wh.Max.X; x++ {
			if r, g, b, _ := img.At(x, y).RGBA(); r == 0 && g == 0 && b == 0 {
				continue
			}
			byteIndex := (y * wh.Max.X / 8) + (x / 8)
			bitIndex := 7 - (x % 8)
			b[byteIndex] |= (1 << bitIndex)
		}
	}
	return b
}

// convertBitsToImage converts a list of bytes (and the size of the image from the bytes) to a black and white image.
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

// loadImage takes an image and returns it as a list of black and white, dithered bytes.
func loadImage(img image.Image) (data []byte, err error) {
	gray := halfgone.FloydSteinbergDitherer{}.Apply(halfgone.ImageToGray(img))
	bytes := convertImageToBits(gray)
	return bytes, nil
}

// formatTransmissionString extracts the data from the input and output by removing the leading and trailing characters
func formatTransmissionString(toFormat string) (formatted string) {
	toFormat = strings.Replace(toFormat, " ", "", -1)
	toFormat = strings.Replace(toFormat, "$", "", -1)
	toFormat = strings.Replace(toFormat, "#", "", -1)
	return toFormat
}

// resizeImage resizes an image to the given width and height
func resizeImage(img image.Image, width, height int) (resized image.Image, err error) {
	resized = imaging.Resize(img, width, height, imaging.Lanczos)
	return resized, nil
}

// cropImage takes an image and crops it to the specified rectangle.
func cropImage(img image.Image, crop image.Rectangle) (croppedImage image.Image, err error) {
	type subImager interface {
		SubImage(r image.Rectangle) image.Image
	}

	// // img is an Image interface. This checks if the underlying value has a
	// // method called SubImage. If it does, then we can use SubImage to crop the
	// // image.
	simg, ok := img.(subImager)
	if !ok {
		return nil, fmt.Errorf("image does not support cropping")
	}

	return simg.SubImage(crop), nil
}
