package waveshareCloud

import (
	"fmt"
	"net"
	"strings"
)

func Send(command string, conn net.Conn) {
	conn.Write([]byte(";" + command + "/" + command))
	fmt.Println("Sent:", command)
}

func Recieve(conn net.Conn) (err error, command string, data string) {
	buf := make([]byte, 64)
	_, err = conn.Read(buf)
	if err != nil {
		return err, "", ""
	}
	command = string(buf)

	buf = make([]byte, 64)
	_, err = conn.Read(buf)
	if err != nil {
		return err, "", ""
	}
	data = string(buf)

	command = strings.Replace(command, " ", "", -1)
	command = strings.Replace(command, "#", "", -1)
	command = strings.Replace(command, "$", "", -1)

	data = strings.Replace(data, " ", "", -1)
	data = strings.Replace(data, "#", "", -1)
	data = strings.Replace(data, "$", "", -1)

	return nil, command, data
}
