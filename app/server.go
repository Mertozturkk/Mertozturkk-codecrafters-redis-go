package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

func main() {

	l, err := net.Listen("tcp", "localhost:6379")
	if err != nil {
		fmt.Println("Error", err)
		return
	}

	defer l.Close()
	fmt.Println("Server listening on ", l.Addr().String())

	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error", err)
			continue
		}
		go handleConnection(conn)
	}
}

type Request struct {
	Command string
	Args    []string
}

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		receiveMessage := string(buf[:n])

		log.Printf("Received Data %s", receiveMessage)
		if errors.Is(err, io.EOF) {
			return
		}
		formattedData := strings.Split(string(receiveMessage), "\r\n")[0]

		splittedMessage := strings.Split(formattedData, " ")
		command := splittedMessage[0]
		parameters := []string{}
		if len(splittedMessage) > 1 {
			parameters = splittedMessage[1:]
		}
		//
		request := Request{
			Command: command,
			Args:    parameters,
		}
		switch request.Command {
		case "echo":
			data := strings.Join(request.Args, "\r\n") + "\r\n"
			conn.Write([]byte(data))
		case "ping":
			conn.Write([]byte("PONG\r\n"))
		}

	}
}
