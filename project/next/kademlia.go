package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/gorilla/websocket"
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
	Monitor   chan string
	State     int
	Busy      int
	Pings     int
	Bootstrap chan Command
	Routing   [][]string
	Addresses map[string]Address
}

func CreateServer(monitor chan string) *Server {

	return &Server{Input: make(chan Command), Monitor: monitor, Addresses: make(map[string]Address)}
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

func (server *Server) send(destination chan Command, command Command) {
	destination <- command
}

func (server *Server) print(status string) {
	server.Monitor <- "Server " + server.ID + " " + status
}

func (server *Server) start(origin chan Command, action StartAction) {
	server.print("Start")

	if action.Contact == nil {

		server.ID = "#"
		server.State = Online
		server.Busy = rand.Intn(100)
		server.Routing = append(server.Routing, []string{server.ID})
		server.Addresses[server.ID] = Address{Busy: server.Busy, Input: server.Input}
	} else {

		command := CreateCommand(server.Input, BootstrapAction{})

		server.send(action.Contact, command)
	}
}

func (server *Server) bootstrap(origin chan Command, action BootstrapAction) {
	server.print("Bootstrap")

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

	server.send(origin, command)
}

func (server *Server) follow(origin chan Command, action FollowAction) {
	server.print("Follow")

	server.send(action.Contact, CreateCommand(server.Input, action.Action))
}

func (server *Server) serve(origin chan Command, action ServeAction) {
	server.print("Serve")

	server.ID = action.ID
	server.State = Online
	server.Busy = rand.Intn(100)
	server.Routing = action.Routing
	server.Addresses = action.Addresses
}

func (server *Server) update(origin chan Command, action UpdateAction) {
	server.print("Update")

	server.Addresses[action.ID] = Address{Busy: action.Busy, Input: origin}
	server.Pings--
}

func (server *Server) ping() {
	server.print("Ping")

	server.State = Online

	channels := make(map[chan Command]bool)
	channels[server.Input] = true

	for _, address := range server.Addresses {

		if !channels[address.Input] {

			channels[address.Input] = true
			server.Pings++

			go server.send(address.Input, Command{Origin: server.Input, Action: UpdateAction{ID: server.ID, Busy: server.Busy}})
		}
	}
}

func (server *Server) run() {

	var pingCh <-chan time.Time

	for {

		if server.State == Online && server.Pings == 0 {
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
				server.print("Quit")
				return
			}
		case <-pingCh:
			server.ping()
		}
	}
}

type KademliaNetwork struct {
	servers   []*Server
	monitor   chan string
	websocket *websocket.Conn
}

func CreateKademliaNetwork(ws *websocket.Conn) *KademliaNetwork {
	kn := &KademliaNetwork{monitor: make(chan string), websocket: ws}
	go kn.Monitor()
	return kn
}
func (kn *KademliaNetwork) Init() {
	server := CreateServer(kn.monitor)
	kn.servers = append(kn.servers, server)
	go server.run()
	server.Input <- CreateCommand(nil, StartAction{})
}

func (kn *KademliaNetwork) New() {
	server := CreateServer(kn.monitor)
	kn.servers = append(kn.servers, server)
	go server.run()
	server.Input <- CreateCommand(nil, StartAction{Contact: kn.servers[0].Input})
}

func (kn *KademliaNetwork) Print() {
	for _, server := range kn.servers {
		kn.monitor <- fmt.Sprintf("%v", server)
	}
}

func (kn *KademliaNetwork) Quit() {
	for _, server := range kn.servers {
		server.Input <- CreateCommand(nil, QuitAction{})
	}
}

func (kn *KademliaNetwork) Monitor() {
	for {
		if information := <-kn.monitor; kn.websocket != nil {
			kn.websocket.WriteMessage(websocket.TextMessage, []byte(information))
		} else {
			fmt.Println(information)
		}
	}
}
