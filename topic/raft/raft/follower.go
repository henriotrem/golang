package raft

func (raft *Raft) destitute(clock int) {

	raft.vote = -1
	raft.state.Set(Follower)
	raft.granted.Reset()
	raft.clock.Set(clock)
}

func (raft *Raft) RequestVote(args *RequestVoteArgs, reply *RequestVoteReply) {

	raft.block()
	defer raft.release()
	defer raft.persist()

	reply.Epoch = Max(args.Epoch, raft.clock.Get())
	reply.VoteGranted = false

	if args.Epoch < raft.clock.Get() {
		return
	}

	if args.Epoch > raft.clock.Get() {
		raft.destitute(args.Epoch)
		go raft.reset()
	}

	if (raft.vote != -1 && raft.vote != args.CandidateID) ||
		(args.LastLogEpoch < raft.logs.lastEpoch() ||
			(args.LastLogEpoch == raft.logs.lastEpoch() && args.LastLogIndex < raft.logs.lastIndex())) {
		return
	}

	raft.vote = args.CandidateID
	reply.VoteGranted = true
}

func (raft *Raft) Message(args *MessageArgs, reply *MessageReply) {

	raft.block()
	defer raft.release()
	defer raft.persist()

	reply.Epoch = Max(args.Epoch, raft.clock.Get())
	reply.Success = false
	reply.ConflictIndex = -1
	reply.ConflictEpoch = -1

	if args.Epoch < raft.clock.Get() {
		return
	}

	if args.Epoch > raft.clock.Get() {
		raft.destitute(args.Epoch)
		go raft.reset()
	}

	go raft.reset()

	if reply.Success, reply.ConflictEpoch, reply.ConflictIndex = raft.checkLogsConflict(args); !reply.Success {
		return
	}

	raft.updateLogs(args)

	if args.LeaderCommit > raft.commitIndex {
		raft.commitIndex = Min(args.LeaderCommit, raft.logs.lastIndex())
		raft.apply()
	}
}

func (raft *Raft) checkLogsConflict(args *MessageArgs) (success bool, conflictEpoch int, conflictIndex int) {

	success = true
	conflictEpoch = -1
	conflictIndex = -1

	if args.PrevLogIndex > raft.logs.lastIndex() {
		success = false
		conflictIndex = raft.logs.lastIndex() + 1

		return
	}

	if epoch := raft.logs.getEpoch(args.PrevLogIndex); epoch != args.PrevLogEpoch {
		success = false
		conflictEpoch = epoch
		conflictIndex, _ = raft.memory.get(epoch)

		return
	}

	return
}

func (raft *Raft) updateLogs(args *MessageArgs) {

	last, first := args.PrevLogIndex+1, 0

	for ; last < raft.logs.lastIndex()+1 && first < len(args.Entries); last, first = last+1, first+1 {
		if raft.logs[last].Epoch != args.Entries[first].Epoch {
			break
		}
	}

	raft.memory.update(raft.clock.Get(), raft.logs[last-1].Epoch, last, args.Entries[first:])
	raft.logs = append(raft.logs[:last], args.Entries[first:]...)
}
