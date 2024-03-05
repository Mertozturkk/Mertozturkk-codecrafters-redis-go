package main

import (
	"errors"
	"fmt"
	"net"
	"strconv"
	"strings"
	"sync"
)

type Store struct {
	KeyValue map[string]string
	Lock     sync.Mutex
}

func NewStore(initialKey, initialValue string) *Store {
	return &Store{
		KeyValue: map[string]string{initialKey: initialValue},
	}
}

func (s *Store) Set(key string, value string) error {
	s.Lock.Lock()
	defer s.Lock.Unlock()

	if len(key) == 0 {
		return errors.New("key cannot be empty")
	}

	s.KeyValue[key] = value
	return nil
}

func (s *Store) Get(key string) (string, bool) {
	s.Lock.Lock()
	defer s.Lock.Unlock()
	val, ok := s.KeyValue[key]
	return val, ok
}

var KeyValue = make(map[string]string)

const (
	Array = '*'
	Bulk  = '$'
	echo  = "echo"
	set   = "set"
	get   = "get"
)

func handleConnection(conn net.Conn, store *Store) {
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
			joinedData += "\r\n"
			conn.Write([]byte("+" + joinedData))
		case "ping":
			conn.Write([]byte("+PONG\r\n"))
		case set:
			err := store.Set(data[0], data[1])
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			conn.Write([]byte("+OK\r\n"))
		case get:
			value, ok := store.Get(data[0])
			if ok {
				conn.Write([]byte("+" + value + "\r\n"))
			} else {
				conn.Write([]byte("+$-1\r\n"))
			}
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

func main() {
	store := NewStore("initialKey", "initialValue")

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
