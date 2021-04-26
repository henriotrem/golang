package raft

func (raft *Raft) institute() {

	raft.state.Set(Leader)

	for idx := range raft.peers {
		raft.nextIndex[idx] = raft.logs.lastIndex() + 1
		raft.matchIndex[idx] = 0
	}
}

func (raft *Raft) newHeartbeat() {

	for idx := range raft.peers {
		if idx != raft.identifier {

			args := &MessageArgs{
				Epoch:        raft.clock.Get(),
				LeaderID:     raft.identifier,
				PrevLogIndex: raft.nextIndex[idx] - 1,
				PrevLogEpoch: raft.logs[raft.nextIndex[idx]-1].Epoch,
				LeaderCommit: raft.commitIndex,
				Entries:      raft.logs[raft.nextIndex[idx]:],
			}

			reply := &MessageReply{}

			go raft.sendMessage(idx, args, reply)
		}
	}
}

func (raft *Raft) sendMessage(identifier int, args *MessageArgs, reply *MessageReply) {
	ok := raft.peers[identifier].Call("Raft.Message", args, reply)

	raft.block()
	defer raft.release()
	defer raft.persist()

	if !ok || raft.state.Get() != Leader || args.Epoch != raft.clock.Get() || reply.Epoch < raft.clock.Get() {
		return
	}

	if reply.Epoch > raft.clock.Get() {
		raft.destitute(reply.Epoch)
		go raft.reset()
		return
	}

	raft.matchIndex[identifier], raft.nextIndex[identifier] = raft.updatePeers(identifier, args, reply)

	if ok, newCommitIndex := raft.newCommitIndex(); ok {

		raft.commitIndex = newCommitIndex
		raft.apply()
	}
}

func (raft *Raft) updatePeers(identifier int, args *MessageArgs, reply *MessageReply) (matchIndex, nextIndex int) {

	if reply.Success {

		matchIndex = Max(raft.matchIndex[identifier], args.PrevLogIndex+len(args.Entries))
		nextIndex = matchIndex + 1
	} else if reply.ConflictEpoch < 0 {

		nextIndex = reply.ConflictIndex
		matchIndex = nextIndex - 1
	} else {

		if index, ok := raft.memory.get(reply.ConflictEpoch); ok {
			nextIndex = index
		} else {
			nextIndex = reply.ConflictIndex
		}
		matchIndex = nextIndex - 1
	}

	return
}

func (raft *Raft) newCommitIndex() (success bool, newCommitIndex int) {

	success = false
	newCommitIndex = -1

	sortedIndex := CopyAndSort(raft.matchIndex)
	epochIndex, ok := raft.memory.get(raft.clock.Get())

	for i := Min(raft.logs.lastIndex(), sortedIndex[len(raft.peers)/2+1]); ok && i > raft.commitIndex && i >= epochIndex; i-- {
		if raft.logs[i].Epoch == raft.clock.Get() {
			success = true
			newCommitIndex = i

			return
		}
	}

	return
}
