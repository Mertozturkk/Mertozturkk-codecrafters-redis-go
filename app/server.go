package main

import (
	"bufio"
	"fmt"
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
	scanner := bufio.NewScanner(conn)
	for scanner.Scan() {
		requestStr := scanner.Text()
		requestParts := strings.Fields(requestStr)
		if len(requestParts) == 0 {
			continue
		}

		command := requestParts[0]
		parameters := []string{}
		if len(requestParts) > 1 {
			parameters = requestParts[1:]
		}

		request := Request{
			Command: command,
			Args:    parameters,
		}

		response := createResponse(request)

		conn.Write([]byte(response + "\r\n"))
	}

	if err := scanner.Err(); err != nil {
		log.Println("Error reading from connection: ", err)
		return
	}
}

func createResponse(request Request) string {
	switch request.Command {
	case "echo":
		return strings.Join(request.Args, " ")
	case "ping":
		return "PONG"
	default:
		return "-ERR unknown command"
	}

}
