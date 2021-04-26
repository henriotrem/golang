package raft

type RequestVoteArgs struct {
	Epoch        int
	CandidateID  int
	LastLogIndex int
	LastLogEpoch int
}

type RequestVoteReply struct {
	Epoch       int
	VoteGranted bool
}

type MessageArgs struct {
	Epoch        int
	LeaderID     int
	PrevLogIndex int
	PrevLogEpoch int
	LeaderCommit int

	Entries []Log
}

type MessageReply struct {
	Epoch         int
	Success       bool
	ConflictIndex int
	ConflictEpoch int
}
