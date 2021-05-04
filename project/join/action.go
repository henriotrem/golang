package main

type StartAction struct {
	Contact string
}

type BootstrapAction struct {
}

type FollowAction struct {
	Contact string
	Action  interface{}
}

type ServeAction struct {
	AddressLookup map[string]Address `json:"address_lookup"`
	Nodes         []Node             `json:"-"`
}

type GetAction struct {
	Nodes map[string]bool `json:"nodes"`
}

type UpdateAction struct {
	Busy int `json:"busy"`
}

type QuitAction struct {
}
