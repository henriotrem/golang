package main

import (
	"fmt"
	"time"
)

type Action int

const (
	Start     Action = 0
	Bootstrap Action = 1
	Follow    Action = 2
	Serve     Action = 3
	Add       Action = 4
	Get       Action = 5
	Quit      Action = 6
)

type Command struct {
	Origin chan Command
	Action Action
	Value  interface{}
}

func CreateCommand(origin chan Command, action Action) Command {
	return Command{Origin: origin, Action: action}
}

type Server struct {
	ID        string
	Input     chan Command
	Quit      chan bool
	Bootstrap chan Command
	Routing   [][]string
	Servers   map[string]chan Command
}

func CreateServer() *Server {

	return &Server{Input: make(chan Command), Quit: make(chan bool), Servers: make(map[string]chan Command)}
}

func (server *Server) start(command Command) {

}

func (server *Server) bootstrap(command Command) {

}

func (server *Server) run() {

	for {
		command := <-server.Input

		switch command.Action {
		case Start:
			fmt.Println("Start")
			server.start(command)
		case Bootstrap:
			fmt.Println("Bootstrap")
			server.bootstrap(command)
		case Follow:
			fmt.Println("Follow")
		case Serve:
			fmt.Println("Serve")
		case Add:
			fmt.Println("Add")
		case Get:
			fmt.Println("Get")
		case Quit:
			fmt.Println("Quit")
		}
	}
}

func main() {

	server1 := CreateServer()
	go server1.bootstrap()

	server2 := CreateServer()
	go server2.bootstrap()

	time.Sleep(3 * time.Second)

	server1.Quit <- true
	server2.Quit <- true
}

// if command.Args.bootstrap == nil {

// 	server.ID = "#"
// 	server.Routing = append(server.Routing, []string{server.ID})
// 	server.Servers[server.ID] = server.Input
// } else {

// 	command := CreateCommand(server.Input, Join)
// 	server.Bootstrap <- command
// }
