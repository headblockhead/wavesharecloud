package waveshareCloud

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strconv"
	"strings"
)

func NewDisplay(conn net.Conn, locked bool) *Display {
	return &Display{
		Connection: conn,
		unlocked:   !locked,
	}
}

// A Display
type Display struct {
	Connection net.Conn
	unlocked   bool
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
		println("locked")
		err = display.SendCommand("N" + password)
		if err != nil {
			return err
		}
		println("sent password")
		data, err = display.UnsafeReceiveData()
		if err != nil {
			return err
		}
		println("recieved")
		println(data)
		if data == "1" {
			println("unlocked")
			display.unlocked = true
			return nil
		} else {
			return fmt.Errorf("wrong password")
		}
	} else {
		return fmt.Errorf("display is already unlocked")
	}
}
