package models

import "time"

const (
	Array = '*'
	Bulk  = '$'
	Echo  = "echo"
	Set   = "set"
	Get   = "get"
	Px    = "px"
	Ping  = "ping"
	INFO  = "info"
)

type CommandList map[string]bool

var Commands = CommandList{
	"echo": true,
	"set":  true,
	"get":  true,
	"ping": true,
	"info": true,
}

var SubCommands = CommandList{
	"px": true,
}

type CliData struct {
	Command string
	Data    []string
	Timer   time.Duration
}

func NewCliData(command string, data []string, timer time.Duration) CliData {
	return CliData{
		Command: command,
		Data:    data,
		Timer:   timer,
	}
}

type ReceiveData struct {
	Type byte
	Data CliData
}

// *2 \r\n $4 \r\n echo\r\n $3\r\n hey \r\n px \r\n 1000 \r\n
