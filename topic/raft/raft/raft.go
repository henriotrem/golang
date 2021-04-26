package raft

import (
	"github.com/henriotrem/topic/raft/labrpc"
)

func (raft *Raft) GetState() (int, bool) {

	raft.block()
	defer raft.release()

	return raft.clock.Get(), raft.state.Get() == Leader
}

func (raft *Raft) Start(command interface{}) (int, int, bool) {

	raft.block()
	defer raft.release()

	lastIndex := raft.logs.lastIndex()
	epoch := raft.clock.Get()

	if raft.state.Get() != Leader {
		return -1, epoch, false
	}

	defer raft.persist()

	raft.logs = append(raft.logs, Log{
		Epoch:   epoch,
		Command: command,
	})

	return lastIndex + 1, epoch, true
}

func (raft *Raft) Kill() {

	raft.block()
	defer raft.release()

	raft.dead.Set(true)
}

func Make(peers []*labrpc.ClientEnd, identifier int,
	persister *Persister, applyCh chan ApplyMsg) *Raft {

	raft := &Raft{
		identifier: identifier,
		peers:      peers,
		persister:  persister,
		logs:       []Log{{Epoch: 0}},
		memory:     Memory{Hashmap: make(map[int]int)},
		applyCh:    applyCh,
		blockCh:    make(chan bool),
		resetCh:    make(chan bool),
		releaseCh:  make(chan bool),
		quitCh:     make(chan bool),
		vote:       -1,
		nextIndex:  make([]int, len(peers)),
		matchIndex: make([]int, len(peers)),
		trace:      false,
	}

	go raft.run()

	// initialize from state persisted before a crash
	raft.readPersist(persister.ReadRaftState())

	return raft
}
