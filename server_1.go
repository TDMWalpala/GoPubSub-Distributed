package main

import (
	"fmt"
	"net"
	"os"
)

func handleConnection(conn net.Conn) {
	for {
		buffer := make([]byte, 1024)
		bytesRead, err := conn.Read(buffer)
		if err != nil {
			fmt.Println("Error reading:", err.Error())
			break
		}

		clientMsg := string(buffer[:bytesRead])
		fmt.Println("Received:", clientMsg)

		if clientMsg == "terminate" {
			break
		}
	}

	conn.Close()
}

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage: go run server.go <port>")
		return
	}

	port := args[1]
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Println("Error listening:", err.Error())
		return
	}

	fmt.Println("Server listening on port", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Error accepting connection:", err.Error())
			break
		}

		go handleConnection(conn)
	}

	listener.Close()
}
