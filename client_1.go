package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	args := os.Args
	if len(args) != 3 {
		fmt.Println("Usage: go run client.go <serverIP> <serverPort>")
		return
	}

	serverIP := args[1]
	serverPort := args[2]

	conn, err := net.Dial("tcp", serverIP+":"+serverPort)
	if err != nil {
		fmt.Println("Error connecting to server:", err.Error())
		return
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Print("Enter message (terminate to exit): ")
		text, _ := reader.ReadString('\n')

		// Trim newline character
		text = strings.TrimSpace(text)

		_, err := conn.Write([]byte(text))
		if err != nil {
			fmt.Println("Error sending message:", err.Error())
			break
		}

		if text == "terminate" {
			break
		}
	}

	conn.Close()
}
