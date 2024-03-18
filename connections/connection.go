package connections

import (
	"fmt"
	"github.com/codecrafters-io/redis-starter-go/formatters"
	"github.com/codecrafters-io/redis-starter-go/models"
	"github.com/codecrafters-io/redis-starter-go/storage"
	"net"
	"strings"
)

type Client struct {
	manager    *NetManager
	connection net.Conn
	send       chan []byte
	message    chan []byte
}

func NewClient(manager *NetManager, connection net.Conn) *Client {
	return &Client{
		manager:    manager,
		connection: connection,
		send:       make(chan []byte),
		message:    make(chan []byte),
	}
}

type NetManager struct {
	listener   *net.Listener
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

func (manager *NetManager) Init() {

	for {
		select {
		case conn := <-manager.register:
			manager.clients[conn] = true
		case conn := <-manager.unregister:
			if _, ok := manager.clients[conn]; ok {
				close(conn.send)
				delete(manager.clients, conn)
			}
		case message := <-manager.broadcast:
			for conn := range manager.clients {
				select {
				case conn.send <- message:
				default:
					close(conn.send)
					delete(manager.clients, conn)
				}
			}
		}

	}
}

func NewNetManager(port string) *NetManager {
	manager := &NetManager{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}

	listener, err := net.Listen("tcp", "localhost:"+port)

	if err != nil {
		fmt.Println("Error", err)
		return nil
	}

	manager.listener = &listener

	return manager
}

func StartNetManager(manager *NetManager, store *storage.Store) {

	defer (*manager.listener).Close()
	fmt.Println("Server listening on ", (*manager.listener).Addr().String())

	for {
		fmt.Println("Waiting for connection")
		conn, err := (*manager.listener).Accept()
		fmt.Println("Connection accepted")
		if err != nil {
			fmt.Println("Error", err)
			continue
		}
		fmt.Println("Creating new client")
		client := NewClient(manager, conn)
		fmt.Println("Registering client")
		manager.register <- client
		fmt.Println("Handling connection")
		go HandleConnection(client, store)
	}
}

func HandleConnection(client *Client, store *storage.Store) {
	defer client.connection.Close()

	for {
		buf := make([]byte, 2048)
		n, err := client.connection.Read(buf)
		if err != nil {
			return
		}
		receiveMessage := string(buf[:n])
		fmt.Println("Received Data", receiveMessage)
		cliData, err := formatters.StringParser(receiveMessage)
		if err != nil {
			fmt.Println(err)
			continue
		}
		cmd := cliData.Command
		fmt.Println("Command", cmd)

		switch cmd {
		case models.INFO:
			fmt.Println("INFO")
			client.connection.Write([]byte("$11\r\nrole:master\r\n"))
		case models.Echo:
			joinedData := strings.Join(cliData.Data, " ")
			joinedData += "\r\n"
			client.connection.Write([]byte("+" + joinedData))

		case models.Ping:
			client.connection.Write([]byte("+PONG\r\n"))

		case models.Set:
			err := store.Set(cliData.Data[0], cliData.Data[1], cliData.Timer)
			if err != nil {
				fmt.Println("Error:", err)
				return
			}
			client.connection.Write([]byte("+OK\r\n"))
		case models.Get:
			value, ok := store.Get(cliData.Data[0])
			if ok {
				client.connection.Write([]byte("+" + value + "\r\n"))
			} else {
				client.connection.Write([]byte("$-1\r\n"))
			}
		}
	}
}
