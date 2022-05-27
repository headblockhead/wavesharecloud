package waveshareCloud

import (
	"errors"
	"net"
	"strings"
)

// A Display
type Display struct {
	Connection net.Conn
}

// Send a formatted command to the display
func (display Display) Send(command string) (err error) {
	_, err = display.Connection.Write([]byte(";" + command + "/" + command))
	if err != nil {
		return err
	}
	return nil
}

// Receive data from the display and format it
func (display Display) Receive(previousCommand string) (command string, data string, err error) {
	// The first set of data is the previous command
	buf := make([]byte, 64)
	conn := display.Connection
	_, err = conn.Read(buf)
	if err != nil {
		return "", "", err
	}
	command = string(buf)
	command = strings.Replace(command, " ", "", -1)
	command = strings.TrimPrefix(command, "$")
	command = strings.TrimSuffix(command, "#")

	if command != previousCommand {
		return command, "", errors.New("command mismatch")
	}

	// The second set of data is the real data
	buf = make([]byte, 64)
	_, err = conn.Read(buf)
	if err != nil {
		return "", "", err
	}
	data = string(buf)

	data = strings.Replace(data, " ", "", -1)
	data = strings.TrimPrefix(data, "$")
	data = strings.TrimSuffix(data, "#")

	return command, data, nil
}
