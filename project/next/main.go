package main

import (
	"fmt"
	"math/rand"
	"time"
)

const (
	Offline = 0
	Online  = 1
)

type StartAction struct {
	Contact chan Command
}

type BootstrapAction struct {
}

type FollowAction struct {
	Contact chan Command
	Action  interface{}
}

type ServeAction struct {
	ID        string
	Routing   [][]string
	Addresses map[string]Address
}

type UpdateAction struct {
	ID   string
	Busy int
}

type QuitAction struct {
}

type Command struct {
	Origin chan Command
	Action interface{}
}

type Address struct {
	Busy  int
	Input chan Command
}

func CreateCommand(origin chan Command, action interface{}) Command {
	return Command{Origin: origin, Action: action}
}

type Server struct {
	ID        string
	Input     chan Command
	State     int
	Busy      int
	Quit      chan bool
	Bootstrap chan Command
	Routing   [][]string
	Addresses map[string]Address
}

func CreateServer() *Server {

	return &Server{Input: make(chan Command), Quit: make(chan bool), Addresses: make(map[string]Address)}
}

func (server *Server) busiest() chan Command {

	var max = server.Busy
	var busiest chan Command

	for _, address := range server.Addresses {

		if server.Input != address.Input && address.Busy > max {
			max = address.Busy
			busiest = address.Input
		}
	}

	return busiest
}

func (server *Server) start(origin chan Command, action StartAction) {
	fmt.Println("Start", server.ID)

	if action.Contact == nil {

		server.ID = "#"
		server.State = Online
		server.Busy = rand.Intn(100)
		server.Routing = append(server.Routing, []string{server.ID})
		server.Addresses[server.ID] = Address{Busy: server.Busy, Input: server.Input}
	} else {

		command := CreateCommand(server.Input, BootstrapAction{})
		action.Contact <- command
	}
}

func (server *Server) bootstrap(origin chan Command, action BootstrapAction) {
	fmt.Println("Bootstrap", server.ID)

	var command Command

	if busy := server.busiest(); busy == nil {

		newServerID := server.ID + "2"
		server.ID = server.ID + "1"

		server.Addresses[newServerID] = Address{Busy: 0, Input: origin}
		server.Addresses[server.ID] = Address{Busy: server.Busy, Input: server.Input}

		server.Routing = append(server.Routing, []string{server.ID, newServerID})

		command = CreateCommand(server.Input, ServeAction{ID: newServerID, Routing: server.Routing, Addresses: server.Addresses})
	} else {

		command = CreateCommand(server.Input, FollowAction{Contact: busy, Action: action})
	}

	origin <- command
}

func (server *Server) follow(origin chan Command, action FollowAction) {
	fmt.Println("Follow", server.ID)

	action.Contact <- CreateCommand(server.Input, action.Action)
}

func (server *Server) serve(origin chan Command, action ServeAction) {
	fmt.Println("Serve", server.ID)

	server.ID = action.ID
	server.State = Online
	server.Busy = rand.Intn(100)
	server.Routing = action.Routing
	server.Addresses = action.Addresses
}

func (server *Server) update(origin chan Command, action UpdateAction) {
	fmt.Println("Update", server.ID)

	server.Addresses[action.ID] = Address{Busy: action.Busy, Input: origin}
}

func (server *Server) ping() {
	fmt.Println("Ping", server.ID)

	channels := make(map[chan Command]bool)
	channels[server.Input] = true

	for _, address := range server.Addresses {

		if !channels[address.Input] {

			channels[address.Input] = true
			address.Input <- Command{Origin: server.Input, Action: UpdateAction{ID: server.ID, Busy: server.Busy}}
		}
	}
}

func (server *Server) run() {

	for {

		var pingCh <-chan time.Time

		if server.State == Online {
			pingCh = time.After(2 * time.Second)
		}

		select {
		case command := <-server.Input:
			switch action := command.Action.(type) {
			case StartAction:
				server.start(command.Origin, action)
			case BootstrapAction:
				server.bootstrap(command.Origin, action)
			case FollowAction:
				server.follow(command.Origin, action)
			case ServeAction:
				server.serve(command.Origin, action)
			case UpdateAction:
				server.update(command.Origin, action)
			case QuitAction:
				return
			}
		case <-pingCh:
			server.ping()
		}

	}
}

func main() {

	var servers []*Server

	for i := 0; i < 10; i++ {

		server := CreateServer()
		go server.run()

		if i == 0 {
			server.Input <- CreateCommand(nil, StartAction{})
		} else {
			server.Input <- CreateCommand(nil, StartAction{Contact: servers[0].Input})

		}

		servers = append(servers, server)
		time.Sleep(3 * time.Second)
	}

	for _, server := range servers {

		server.Quit <- true
		fmt.Println(server)
	}
}
