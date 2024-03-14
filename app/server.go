package main

import (
	"flag"
	"github.com/codecrafters-io/redis-starter-go/connections"
	"github.com/codecrafters-io/redis-starter-go/storage"
)

func main() {
	store := storage.NewStore("initialKey", "initialValue")
	// use flag for port
	port := flag.String("port", "6379", "port to listen on")
	flag.Parse()

	manager := connections.NewNetManager(*port)
	go manager.Init()
	connections.StartNetManager(manager, store)
}
