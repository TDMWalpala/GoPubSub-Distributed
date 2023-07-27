package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

type Client struct {
	conn       net.Conn
	isPublisher bool
}

type Message struct {
	senderID string
	text     string
}

var (
	clients     = make(map[net.Conn]Client)
	clientsLock sync.Mutex
	subscribers []net.Conn
	subLock     sync.Mutex
	messageCh   = make(chan Message)
)

func handleConnection(conn net.Conn) {
	defer conn.Close()

	client := Client{
		conn:       conn,
		isPublisher: false,
	}

	clientsLock.Lock()
	clients[conn] = client
	clientsLock.Unlock()

	reader := bufio.NewReader(conn)

	clientType, err := reader.ReadString('\n')
	if err != nil {
		fmt.Println("Error reading client type:", err.Error())
		return
	}

	clientType = strings.TrimSpace(clientType)

	if clientType == "PUBLISHER" {
		client.isPublisher = true
	} else if clientType != "SUBSCRIBER" {
		fmt.Println("Invalid client type. Please use 'PUBLISHER' or 'SUBSCRIBER'")
		return
	}

	if !client.isPublisher {
		subLock.Lock()
		subscribers = append(subscribers, conn)
		subLock.Unlock()
	}

	for {
		message, err := reader.ReadString('\n')
		if err != nil {
			break
		}

		message = strings.TrimSpace(message)

		if client.isPublisher {
			handlePublisherMessage(conn, message)
		} else {
			handleSubscriberMessage(message)
		}
	}

	clientsLock.Lock()
	delete(clients, conn)
	clientsLock.Unlock()

	if !client.isPublisher {
		subLock.Lock()
		for i, sub := range subscribers {
			if sub == conn {
				subscribers = append(subscribers[:i], subscribers[i+1:]...)
				break
			}
		}
		subLock.Unlock()
	}
}

func handlePublisherMessage(publisherConn net.Conn, message string) {
	msg := Message{
		senderID: publisherConn.RemoteAddr().String(),
		text:     message,
	}

	messageCh <- msg
}

func handleSubscriberMessage(message string) {
	fmt.Println("Received message from publisher:", message)

	subLock.Lock()
	for _, sub := range subscribers {
		_, err := sub.Write([]byte(message + "\n"))
		if err != nil {
			fmt.Println("Error sending message to subscriber:", err.Error())
		}
	}
	subLock.Unlock()
}

func broadcastMessages() {
	for {
		msg := <-messageCh

		subLock.Lock()
		for _, sub := range subscribers {
			_, err := sub.Write([]byte(msg.text + "\n"))
			if err != nil {
				fmt.Println("Error sending message to subscriber:", err.Error())
			}
		}
		subLock.Unlock()
	}
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

	go broadcastMessages()

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
