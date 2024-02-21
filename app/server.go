package main

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
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

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		// fmt.Println("Received int", n)
		receiveMessage := string(buf[:n])
		log.Printf("Received Data %s", receiveMessage)
		if errors.Is(err, io.EOF) {
			return
		}
		conn.Write([]byte("+PONG\r\n"))
	}

}
