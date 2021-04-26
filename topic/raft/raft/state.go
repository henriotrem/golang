package raft

import "sync/atomic"

type State struct {
	value int32
}

const (
	Follower  = 0
	Candidate = 1
	Leader    = 2
)

func (state *State) Get() int {
	return int(atomic.LoadInt32(&state.value))
}
func (state *State) Set(value int) {
	atomic.StoreInt32(&state.value, int32(value))
}

func (state *State) Print() string {
	var label string
	switch state.Get() {
	case 0:
		label = "Follower"
	case 1:
		label = "Candidate"
	case 2:
		label = "Leader"
	}
	return label
}
