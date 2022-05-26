package waveshareCloud

import (
	"fmt"
	"net"
	"os"
	"strings"
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
		fmt.Println("Error listening:", err.Error())
		os.Exit(1)
	}
	// Close the listener when the application closes.
	defer l.Close()
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	for {
		// Listen for an incoming connection.
		conn, err := l.Accept()
		fmt.Println("Accepted connection")
		if err != nil {
			fmt.Println("Error accepting: ", err.Error())
			os.Exit(1)
		}
		// Handle connections in a new goroutine.
		go handleRequest(conn)
	}
}

func send(command string, conn net.Conn) {
	conn.Write([]byte(";" + command + "/" + command))
	fmt.Println("Sent:", command)
}

func recieve(conn net.Conn) (command string, data string) {
	buf := make([]byte, 64)
	_, err := conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	command = string(buf)

	buf = make([]byte, 64)
	_, err = conn.Read(buf)
	if err != nil {
		fmt.Println("Error reading:", err.Error())
	}
	data = string(buf)

	command = strings.Replace(command, " ", "", -1)
	command = strings.Replace(command, "#", "", -1)
	command = strings.Replace(command, "$", "", -1)

	data = strings.Replace(data, " ", "", -1)
	data = strings.Replace(data, "#", "", -1)
	data = strings.Replace(data, "$", "", -1)

	fmt.Println("Received:", data, "From:", command)
	return command, data
}

func handleRequest(conn net.Conn) {
	fmt.Println("handle connection")
	// Get device ID
	send("G", conn)
	command, data := recieve(conn)
	print(command, data)
	// Shutdown the connection.
	send("S", conn)
	conn.Close()
}
