package main

import (
	"flag"
	"github.com/codecrafters-io/redis-starter-go/connections"
	"github.com/codecrafters-io/redis-starter-go/storage"
)

var role = "master"

func main() {

	port := flag.String("port", "6379", "port to listen on")
	replicaofHost := flag.String("replicaof", "", "Replicate to another server")
	flag.Parse()

	if *replicaofHost != "" {
		role = "slave"
	}

	store := storage.NewStore("initialKey", "initialValue")
	manager := connections.NewNetManager(*port, role)
	go manager.Init()
	connections.StartNetManager(manager, store)
}
