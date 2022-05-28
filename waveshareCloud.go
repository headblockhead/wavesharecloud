package waveshareCloud

import (
	"encoding/binary"
	"fmt"
	"net"
	"strings"
)

func NewDisplay(conn net.Conn) *Display {
	return &Display{
		Connection: conn,
	}
}

// A Display
type Display struct {
	Connection net.Conn
}

// SendCommand sends a formatted command to the display
func (display *Display) SendCommand(command string) (err error) {
	_, err = display.Connection.Write([]byte(";" + command + "/" + command))
	if err != nil {
		return err
	}
	return nil
}

// SendData sends formatted image data to the display
func (display *Display) SendData(data []byte) (err error) {
	// 0x57+4Byte addr+4Byte len+1Byte num+len Byte data +Verify
	prefix := uint32(0x57)
	var check uint32
	for i := 0; i < len(data); i++ {
		check = check ^ uint32(data[i])
	}
	err = binary.Write(display.Connection, binary.LittleEndian, prefix)
	if err != nil {
		return err
	}
	err = binary.Write(display.Connection, binary.LittleEndian, data)
	if err != nil {
		return err
	}
	err = binary.Write(display.Connection, binary.LittleEndian, check)
	if err != nil {
		return err
	}
	return nil
}

// ReceiveData Receives data from the display and formats it
func (display *Display) ReceiveData(previousCommand string) (command string, data string, err error) {
	// The first set of data is the previous command
	buf := make([]byte, 64)
	_, err = display.Connection.Read(buf)
	if err != nil {
		return "", "", err
	}
	command = string(buf)
	command = strings.Replace(command, " ", "", -1)
	command = strings.Replace(command, "$", "", -1)
	command = strings.Replace(command, "#", "", -1)
	if !strings.Contains(command, previousCommand) {
		return command, "", fmt.Errorf("command mismatch: expected %s, got %v", previousCommand, command)
	}

	// The second set of data is the real data
	buf = make([]byte, 64)
	_, err = display.Connection.Read(buf)
	if err != nil {
		return "", "", err
	}
	data = string(buf)

	data = strings.Replace(data, " ", "", -1)
	data = strings.Replace(data, "$", "", -1)
	data = strings.Replace(data, "#", "", -1)

	return command, data, nil
}

// Disconnect closes the connection to the display
func (display *Display) Disconnect() {
	display.Connection.Close()
}
