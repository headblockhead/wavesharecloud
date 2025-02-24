package main

import (
	"fmt"
	"net"
	"os"

	"github.com/headblockhead/wavesharecloud"
)

const (
	CONN_HOST = "0.0.0.0"
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
	fmt.Println("Listening on " + CONN_HOST + ":" + CONN_PORT)
	defer l.Close()
	
	for {
		// Wait for an incoming connection.
		conn, err := l.Accept()
		if err != nil {
			fmt.Printf("Error accepting a connection: %v", err)
			os.Exit(1)
		}
		// Handle the connection in a new goroutine.
		go handleRequest(conn)
	}
}

func handleRequest(conn net.Conn) {
	fmt.Println("New connection from:", conn.RemoteAddr())
	lc := wavesharecloud.NewLoggingConn(conn, false)

	// If a password is required to unlock the display, it should be provided as the second argument, and the display will be unlocked.
	// In this case, the display is not locked, so the password is not required.
	display := wavesharecloud.NewDisplay(lc, "")

	battery, err := display.GetBatteryLevel()
	if err != nil {
		fmt.Printf("Error getting battery level: %v\n", err)
	}
	fmt.Println("Battery level:", battery)

	err := display.Shutdown()
	if err != nil {
		fmt.Printf("Error shutting down: %v\n", err)
	}
	display.Disconnect()
}
