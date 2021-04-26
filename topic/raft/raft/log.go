package raft

import "fmt"

type Log struct {
	Epoch   int
	Command interface{}
}

type Logs []Log

func (logs Logs) lastIndex() int {
	return len(logs) - 1
}

func (logs Logs) lastEpoch() int {
	if logs.lastIndex() == -1 {
		return -1
	}
	return logs[logs.lastIndex()].Epoch
}

func (logs Logs) getEpoch(index int) int {
	return logs[index].Epoch
}

func (logs Logs) print() string {
	result := ""
	for _, log := range logs {
		bytes := fmt.Sprintf("%v", log.Command)
		length := Min(len(bytes), 10)
		result += ", {" + fmt.Sprint(log.Epoch) + " " + string(bytes[:length]) + "}"
	}
	return result
}
