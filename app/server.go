package main

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/formatters"
	"github.com/codecrafters-io/redis-starter-go/models"
	"github.com/codecrafters-io/redis-starter-go/storage"
	"net"
	"strings"
)

func main() {
	store := storage.NewStore("initialKey", "initialValue")

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
		go handleConnection(conn, store)
	}
}

func handleConnection(conn net.Conn, store *storage.Store) {
	defer conn.Close()

	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		receiveMessage := string(buf[:n])

		cliData, err := formatters.StringParser(receiveMessage)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cmd := cliData.Command

		switch cmd {
		case models.Echo:
			joinedData := strings.Join(cliData.Data, " ")
			joinedData += "\r\n"
			conn.Write([]byte("+" + joinedData))

		case models.Ping:
			conn.Write([]byte("+PONG\r\n"))

		case models.Set:
			err := store.Set(cliData.Data[0], cliData.Data[1], cliData.Timer)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			conn.Write([]byte("+OK\r\n"))

		case models.Get:
			value, ok := store.Get(cliData.Data[0])
			if ok {
				conn.Write([]byte("+" + value + "\r\n"))
			} else {
				conn.Write([]byte("$-1\r\n"))
			}
		}

	}
}
