package raft

func (raft *Raft) propose() {

	raft.state.Set(Candidate)
	raft.clock.Increment()
	raft.memory.add(raft.clock.Get(), raft.logs.lastIndex()+1)
	raft.granted.Reset()
	raft.granted.Increment()
	raft.vote = raft.identifier
}

func (raft *Raft) newElection() {

	raft.propose()
	raft.persist()

	for idx := range raft.peers {
		if idx != raft.identifier {

			args := &RequestVoteArgs{
				Epoch:        raft.clock.Get(),
				CandidateID:  raft.identifier,
				LastLogIndex: raft.logs.lastIndex(),
				LastLogEpoch: raft.logs.lastEpoch(),
			}

			reply := &RequestVoteReply{}

			go raft.sendRequestVote(idx, args, reply)
		}
	}
}

func (raft *Raft) sendRequestVote(identifier int, args *RequestVoteArgs, reply *RequestVoteReply) {
	ok := raft.peers[identifier].Call("Raft.RequestVote", args, reply)

	raft.block()
	defer raft.release()

	if !ok || raft.state.Get() != Candidate || args.Epoch != raft.clock.Get() || reply.Epoch < raft.clock.Get() {
		return
	}

	if reply.Epoch > raft.clock.Get() {
		raft.destitute(reply.Epoch)
		raft.persist()
		go raft.reset()
		return
	}

	if reply.VoteGranted {
		raft.granted.Increment()

		if raft.granted.Get() == len(raft.peers)/2+1 {
			raft.institute()
			go raft.reset()
			return
		}
	}
}
