package main

import (
	"encoding/json"
	"math/rand"
	"time"
)

const (
	Offline = 0
	Online  = 1
)

type Server struct {
	Address Address `json:"address"`

	NodeLookup    map[string]*Node    `json:"-"`
	AddressLookup map[string]*Address `json:"address_lookup"`

	Input   chan Command `json:"-"`
	Network chan Command `json:"-"`
	Monitor chan string  `json:"-"`
}

func CreateServer(network chan Command, monitor chan string) *Server {

	server := &Server{
		NodeLookup:    make(map[string]*Node),
		AddressLookup: make(map[string]*Address),
		Input:         make(chan Command),
		Network:       network,
		Monitor:       monitor}

	return server
}

func (server *Server) busiest() string {

	var max = server.Address.Busy
	var busiest string

	for _, address := range server.AddressLookup {

		if server.Address.IP != address.IP && address.Busy > max {
			max = address.Busy
			busiest = address.IP
		}
	}

	return busiest
}

func (server *Server) send(command Command) {
	server.Network <- command
}

func (server *Server) print() {
	srv, _ := json.Marshal(server)
	tree, _ := json.Marshal(server.NodeLookup["#"])
	server.Monitor <- "STATE\t" + server.Address.IP + "\t" + string(srv) + "\t" + string(tree)
}

func (server *Server) start(origin, destination string, action StartAction) {

	server.Address = Address{IP: destination, Leaves: make(map[string]bool)}
	server.AddressLookup[server.Address.IP] = &server.Address

	if action.Contact == "" {

		root := &Node{Key: "#", Count: 0}

		server.Address.Leaves[root.Key] = true
		server.NodeLookup[root.Key] = root

		server.Address.State = Online
		server.Address.Busy = rand.Intn(10000)
	} else {

		command := CreateCommand(server.Address.IP, action.Contact, BootstrapAction{})
		go server.send(command)
	}
}

func (server *Server) bootstrap(origin, destination string, action BootstrapAction) {

	var command Command

	if busyIP := server.busiest(); busyIP == "" {

		server.AddressLookup[origin] = &Address{IP: origin, Leaves: make(map[string]bool)}

		var firstHalf []*Node
		var secondHalf []*Node

		length := len(server.Address.Leaves) / 2
		for key := range server.Address.Leaves {
			leaf := server.NodeLookup[key]
			if length > 0 {
				secondHalf = append(secondHalf, leaf)
				length--
			} else {
				leaf.Split()
				half := len(leaf.Children) / 2
				firstHalf, secondHalf = append(firstHalf, leaf.Children[:half]...), append(secondHalf, leaf.Children[half:]...)
			}
			delete(server.Address.Leaves, key)
			if length == 0 {
				break
			}
		}

		for _, leaf := range firstHalf {
			server.NodeLookup[leaf.Key] = leaf
			server.Address.Leaves[leaf.Key] = true
		}

		for _, leaf := range secondHalf {
			server.NodeLookup[leaf.Key] = leaf
			server.AddressLookup[origin].Leaves[leaf.Key] = true
		}

		var nodes []Node

		for _, node := range server.NodeLookup {
			nodes = append(nodes, *node)
		}

		var addresses = make(map[string]Address)

		for key, address := range server.AddressLookup {
			addresses[key] = *address
		}

		server.Address.Busy = server.Address.Busy / 2

		command = CreateCommand(server.Address.IP, origin, ServeAction{Nodes: nodes, AddressLookup: addresses})
	} else {

		command = CreateCommand(server.Address.IP, origin, FollowAction{Contact: busyIP, Action: action})
	}

	go server.send(command)
}

func (server *Server) follow(origin, destination string, action FollowAction) {

	go server.send(CreateCommand(server.Address.IP, action.Contact, action.Action))
}

func (server *Server) serve(origin, destination string, action ServeAction) {

	for ip, address := range action.AddressLookup {

		if _, ok := server.AddressLookup[ip]; !ok {
			server.AddressLookup[ip] = &Address{IP: address.IP, Leaves: map[string]bool{}}
		}
		for leaf := range address.Leaves {
			server.AddressLookup[ip].Leaves[leaf] = true
		}
	}

	for _, node := range action.Nodes {
		server.NodeLookup[node.Key] = &Node{Key: node.Key, Count: node.Count}
	}

	for _, node := range action.Nodes {
		current := server.NodeLookup[node.Key]
		if parent, ok := server.NodeLookup[node.GetParentKey()]; ok {
			current.Parent = parent
			parent.Children = append(parent.Children, current)
		}
	}

	server.Address.State = Online
	server.Address.Busy = rand.Intn(10000)
}

func (server *Server) get(origin, destination string, action GetAction) {

	go server.send(CreateCommand(server.Address.IP, origin, UpdateAction{Busy: server.Address.Busy}))
}

func (server *Server) update(origin, destination string, action UpdateAction) {

	server.AddressLookup[origin].Busy = action.Busy
}

func (server *Server) ping() {

	for ip, address := range server.AddressLookup {
		if server.Address.IP != ip {
			go server.send(Command{Origin: server.Address.IP, Destination: ip, Action: GetAction{Nodes: address.Leaves}})
		}
	}
}

func (server *Server) run() {

	var pingCh = time.After(2 * time.Second)

	for {

		select {
		case command := <-server.Input:
			switch action := command.Action.(type) {
			case StartAction:
				server.start(command.Origin, command.Destination, action)
			case BootstrapAction:
				server.bootstrap(command.Origin, command.Destination, action)
			case FollowAction:
				server.follow(command.Origin, command.Destination, action)
			case ServeAction:
				server.serve(command.Origin, command.Destination, action)
			case GetAction:
				server.get(command.Origin, command.Destination, action)
			case UpdateAction:
				server.update(command.Origin, command.Destination, action)
			case QuitAction:
				return
			}
		case <-pingCh:
			server.ping()
			pingCh = time.After(2 * time.Second)
		}

		server.print()
	}
}
