package main

import (
	"fmt"
	"net"
	"strconv"
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

const (
	Array = '*'
	Bulk  = '$'
	echo  = "echo"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()
	for {
		buf := make([]byte, 2048)
		n, err := conn.Read(buf)
		if err != nil {
			return
		}
		receiveMessage := string(buf[:n])
		arr := StringParser(receiveMessage)
		cmd, data := ArrayReader(arr)

		switch cmd {
		case echo:
			joinedData := strings.Join(data, "\r\n")
			conn.Write([]byte(joinedData))
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
		}
	}
}

func StringParser(s string) []string {
	var arr []string
	parseredData := strings.Split(s, "\r\n")
	for i := 0; i < len(parseredData); i++ {
		arr = append(arr, parseredData[i])
	}
	return arr
}

func ArrayReader(arr []string) (string, []string) {
	var data []string
	cmd := arr[2]
	for i := 3; i < len(arr); i++ {
		element := arr[i]
		if element == "" {
			continue
		}
		type_ := element[0]
		if strconv.Itoa(int(type_)) == "" {
			continue
		}
		if type_ == Bulk {
			continue
		}
		data = append(data, element)
	}
	return cmd, data
}
