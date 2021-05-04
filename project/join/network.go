package main

import (
	"encoding/json"
	"fmt"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Address struct {
	IP      string          `json:"IP"`
	State   int             `json:"state"`
	Busy    int             `json:"busy"`
	Busiest int             `json:"busiest"`
	Leaves  map[string]bool `json:"leaves"`
}

type Command struct {
	Origin      string
	Destination string
	Action      interface{}
}

func CreateCommand(origin, destination string, action interface{}) Command {
	return Command{Origin: origin, Destination: destination, Action: action}
}

type Network struct {
	servers    []*Server
	input      chan Command
	lookup     map[string]chan Command
	connection *websocket.Conn
	monitoring chan string
}

func CreateNetwork(connection *websocket.Conn) *Network {
	network := &Network{input: make(chan Command), lookup: make(map[string]chan Command), connection: connection, monitoring: make(chan string)}

	go network.run()
	go network.monitor()

	return network
}
func (network *Network) Init() {
	server := CreateServer(network.input, network.monitoring)
	network.servers = append(network.servers, server)

	go server.run()

	newIp := uuid.NewString()
	network.lookup[newIp] = server.Input
	network.lookup[newIp] <- CreateCommand("", newIp, StartAction{})
}

func (network *Network) New() {
	server := CreateServer(network.input, network.monitoring)
	network.servers = append(network.servers, server)

	go server.run()

	newIp := uuid.NewString()
	network.lookup[newIp] = server.Input
	network.lookup[newIp] <- CreateCommand("", newIp, StartAction{Contact: network.servers[0].Address.IP})
}

func (network *Network) Print() {
	for _, server := range network.servers {
		network.monitoring <- fmt.Sprintf("%v", server)
	}
}

func (network *Network) Quit() {
	for _, server := range network.servers {
		server.Input <- CreateCommand("", server.Address.IP, QuitAction{})
	}
}

func (network *Network) run() {
	for {
		command := <-network.input
		data, _ := json.Marshal(command.Action)
		network.monitoring <- fmt.Sprintf("NETWORK\t%v\t%v\t%T\t%v", command.Origin, command.Destination, command.Action, string(data))
		network.lookup[command.Destination] <- command
	}
}

func (network *Network) monitor() {
	for {
		if information := <-network.monitoring; network.connection != nil {
			network.connection.WriteMessage(websocket.TextMessage, []byte(information))
		} else {
			fmt.Println(information)
		}
	}
}
