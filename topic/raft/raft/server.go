package raft

import (
	"bytes"
	"fmt"
	"math/rand"
	"time"

	"github.com/henriotrem/topic/raft/labgob"
	"github.com/henriotrem/topic/raft/labrpc"
)

type ApplyMsg struct {
	CommandValid bool
	Command      interface{}
	CommandIndex int
}

type Raft struct {
	identifier int
	peers      []*labrpc.ClientEnd
	persister  *Persister
	state      State
	clock      Clock
	logs       Logs
	memory     Memory
	dead       Flag

	applyCh   chan ApplyMsg
	blockCh   chan bool
	resetCh   chan bool
	releaseCh chan bool
	quitCh    chan bool

	vote    int
	granted Counter

	nextIndex  []int
	matchIndex []int

	commitIndex int
	lastApplied int

	trace bool
}

func (raft *Raft) print() {
	if raft.trace {
		fmt.Println("Server", raft.identifier, "-", raft.state.Print(), raft.clock.Get(), raft.commitIndex, "Logs", len(raft.logs), raft.logs.print(), "Memory", raft.memory.print())
	}
}

func (raft *Raft) apply() {

	for i := raft.lastApplied + 1; i <= raft.commitIndex; i++ {
		raft.applyCh <- ApplyMsg{
			CommandValid: true,
			Command:      raft.logs[i].Command,
			CommandIndex: i,
		}
		raft.lastApplied = i
	}
}

func (raft *Raft) persist() {

	w := new(bytes.Buffer)
	e := labgob.NewEncoder(w)
	if e.Encode(raft.clock.Get()) != nil ||
		e.Encode(raft.vote) != nil ||
		e.Encode(raft.logs) != nil ||
		e.Encode(raft.memory) != nil {
		panic("failed to encode raft persistent state")
	}
	data := w.Bytes()
	raft.persister.SaveRaftState(data)
}

func (raft *Raft) readPersist(data []byte) {
	if data == nil || len(data) < 1 {
		return
	}
	r := bytes.NewBuffer(data)
	d := labgob.NewDecoder(r)
	if d.Decode(&raft.clock.value) != nil ||
		d.Decode(&raft.vote) != nil ||
		d.Decode(&raft.logs) != nil ||
		d.Decode(&raft.memory) != nil {
		panic("failed to decode raft persistent state")
	}
}

func (raft *Raft) newTimeouts() (<-chan time.Time, <-chan time.Time) {
	if raft.state.Get() != Leader {
		return time.After(time.Duration(360+rand.Intn(240)) * time.Millisecond), time.After(60 * time.Second)
	} else {
		return time.After(60 * time.Second), time.After(60 * time.Millisecond)
	}
}

func (raft *Raft) block() {
	raft.blockCh <- true
}

func (raft *Raft) wait() {
	raft.releaseCh <- true
}

func (raft *Raft) release() {
	<-raft.releaseCh
}

func (raft *Raft) reset() {
	raft.resetCh <- true
	<-raft.releaseCh
}

func (raft *Raft) run() {

	var candidateCh, heartbeatCh <-chan time.Time
	var reset = true

	for !raft.dead.Get() {

		if reset {
			candidateCh, heartbeatCh = raft.newTimeouts()
		}

		reset = true

		select {
		case <-raft.blockCh:
			reset = false
			raft.wait()
		case <-raft.resetCh:
			raft.wait()
		case <-candidateCh:
			raft.newElection()
		case <-heartbeatCh:
			raft.newHeartbeat()
		}

		if reset {
			raft.print()
		}
	}
}
