package main

import (
	"fmt"
	"net"
	"os"
	"strings"
	"bufio"
)

func main() {
	args := os.Args
	if len(args) != 4 {
		fmt.Println("Usage: go run client.go <serverIP> <serverPort> <clientType>")
		return
	}

	serverIP := args[1]
	serverPort := args[2]
	clientType := strings.ToUpper(args[3])

	conn, err := net.Dial("tcp", serverIP+":"+serverPort)
	if err != nil {
		fmt.Println("Error connecting to server:", err.Error())
		return
	}

	defer conn.Close()

	if clientType == "PUBLISHER" {
		fmt.Println("Publisher mode")
		_, err := conn.Write([]byte("PUBLISHER\n"))
		if err != nil {
			fmt.Println("Error sending publisher mode:", err.Error())
			return
		}

		reader := bufio.NewReader(os.Stdin)
		for {
			message, _ := reader.ReadString('\n')

			// Trim newline character
			message = strings.TrimSpace(message)

			_, err := conn.Write([]byte(message + "\n"))
			if err != nil {
				fmt.Println("Error sending message:", err.Error())
				break
			}

			if message == "terminate" {
				break
			}
		}
	} else if clientType == "SUBSCRIBER" {
		fmt.Println("Subscriber mode")
		_, err := conn.Write([]byte("SUBSCRIBER\n"))
		if err != nil {
			fmt.Println("Error sending subscriber mode:", err.Error())
			return
		}

		go receiveMessages(conn)

		// Block the main goroutine to allow receiving messages continuously
		select {}
	} else {
		fmt.Println("Invalid client type. Please use 'PUBLISHER' or 'SUBSCRIBER'")
		return
	}
}

func receiveMessages(conn net.Conn) {
	reader := bufio.NewReader(conn)

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		message = strings.TrimSpace(message)
		fmt.Println("Received:", message)
	}

	conn.Close()
}
