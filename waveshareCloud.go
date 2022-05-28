package waveshareCloud

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"net"
	"strconv"
	"strings"
)

func NewLoggingConn(conn net.Conn) *LoggingConn {
	return &LoggingConn{
		Connection: conn,
	}
}

type LoggingConn struct {
	Connection net.Conn
}

func (lc *LoggingConn) Read(b []byte) (n int, err error) {
	n, err = lc.Connection.Read(b)
	fmt.Println("< " + string(b))
	return n, err
}

func (lc *LoggingConn) Write(b []byte) (n int, err error) {
	fmt.Println("> " + string(b))
	return lc.Connection.Write(b)
}

func (lc *LoggingConn) Close() error {
	return lc.Connection.Close()
}

func NewDisplay(conn io.ReadWriteCloser, locked bool) *Display {
	return &Display{
		Connection: conn,
		unlocked:   !locked,
		Width:      400,
		Height:     300,
	}
}

// A Display
type Display struct {
	Connection io.ReadWriteCloser
	unlocked   bool
	Width      int
	Height     int
}

// SendCommand sends a formatted command to the display
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

func (display *Display) SendImage(data []byte) (err error) {
	println("start IMAGE SEND")
	if len(data) != (display.Width*display.Height)/8 {
		return fmt.Errorf("data length does not match display size")
	}
	println("SENDING F to pay respects")
	err = display.SendCommand("F")
	if err != nil {
		return err
	}
	_, err = display.UnsafeReceiveData()
	if err != nil {
		return err
	}
	var count uint8
	for i := 0; i < len(data); i += 1024 {
		to := i + 1024
		if to > len(data) {
			to = len(data)
		}
		err = display.SendFrame(uint32(i), count, data[i:to])
		if err != nil {
			return err
		}
		_, err = display.UnsafeReceiveData()
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
	_, err = display.UnsafeReceiveData()
	if err != nil {
		return err
	}
	return nil
}

var closeFrame = []byte{
	0x57,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0,
	0x00,
}

// SendFrame sends formatted image data to the display
func (display *Display) SendCloseFrame() (err error) {
	_, err = display.Connection.Write(closeFrame)
	if err != nil {
		return err
	}
	_, err = display.UnsafeReceiveData()
	if err != nil {
		return err
	}
	return err
}

// SendFrame sends formatted image data to the display
func (display *Display) SendFrame(addr uint32, num uint8, data []byte) (err error) {
	if len(data) > 1024 {
		return fmt.Errorf("data too large, maximum size is 1024")
	}
	frame := new(bytes.Buffer)
	// 0x57
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
	err = binary.Write(frame, binary.BigEndian, uint32(len(data)))
	if err != nil {
		return err
	}
	// 1 byte num
	err = binary.Write(frame, binary.BigEndian, num)
	if err != nil {
		return err
	}
	// data
	err = binary.Write(frame, binary.BigEndian, data)
	if err != nil {
		return err
	}
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
	// verify
	var check byte
	payload := frame.Bytes()
	for i := 1; i < len(payload[1:len(data)]); i++ {
		check = check ^ payload[i]
	}
	err = binary.Write(frame, binary.BigEndian, check)
	if err != nil {
		return err
	}
	display.Connection.Write(frame.Bytes())
	return nil
}

// ReceiveData Receives data from the display and formats it
func (display *Display) ReceiveData(previousCommand string) (data string, err error) {
	// The first set of data is the previous command
	buf := make([]byte, 64)
	_, err = display.Connection.Read(buf)
	if err != nil {
		return "", err
	}
	command := string(buf)
	command = strings.Replace(command, " ", "", -1)
	command = strings.Replace(command, "$", "", -1)
	command = strings.Replace(command, "#", "", -1)
	if !strings.Contains(command, previousCommand) {
		return "", fmt.Errorf("command mismatch: expected %s, got %v", previousCommand, command)
	}
	// The second set of data is the real data
	buf = make([]byte, 64)
	_, err = display.Connection.Read(buf)
	if err != nil {
		return "", err
	}
	data = string(buf)

	data = strings.Replace(data, " ", "", -1)
	data = strings.Replace(data, "$", "", -1)
	data = strings.Replace(data, "#", "", -1)

	return data, nil
}

//UnsafeReceiveData receives data from the display without checking if it is valid and only returns the first set of data it gets
func (display *Display) UnsafeReceiveData() (data string, err error) {
	// The first set of data is the previous command
	buf := make([]byte, 64)
	_, err = display.Connection.Read(buf)
	if err != nil {
		return "", err
	}
	command := string(buf)
	command = strings.Replace(command, " ", "", -1)
	command = strings.Replace(command, "$", "", -1)
	command = strings.Replace(command, "#", "", -1)

	return command, nil
}

func (display *Display) Shutdown() (err error) {
	if display.unlocked {
		display.SendCommand("S")
		return nil
	} else {
		return fmt.Errorf("display is locked")
	}
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
	data, err := display.ReceiveData("G")
	if err != nil {
		return "", err
	}
	return data, nil
}

// GetLocked gets whether the display is locked with a PIN
func (display *Display) GetLocked() (unlocked bool, err error) {
	err = display.SendCommand("C")
	if err != nil {
		return false, err
	}
	data, err := display.ReceiveData("C")
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

// GetLocked gets whether the display is locked with a PIN
func (display *Display) Unlock(password string) (err error) {
	err = display.SendCommand("C")
	if err != nil {
		return err
	}
	data, err := display.ReceiveData("C")
	if err != nil {
		return err
	}
	if strings.Contains(data, "1") {
		err = display.SendCommand("N" + password)
		if err != nil {
			return err
		}
		data, err = display.UnsafeReceiveData()
		data, err = display.UnsafeReceiveData()
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
		return fmt.Errorf("display is already unlocked")
	}
}
