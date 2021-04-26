package raft

import "fmt"

type Memory struct {
	Epochs  []int
	Hashmap map[int]int
}

func (memory *Memory) add(epoch, index int) {

	lastPosition := len(memory.Epochs) - 1

	if lastPosition >= 0 && memory.Hashmap[memory.Epochs[lastPosition]] == index {
		delete(memory.Hashmap, memory.Epochs[lastPosition])
		memory.Epochs = memory.Epochs[:lastPosition]
	}

	memory.Epochs = append(memory.Epochs, epoch)
	memory.Hashmap[epoch] = index
}

func (memory *Memory) get(epoch int) (int, bool) {
	index, ok := memory.Hashmap[epoch]
	return index, ok
}

func (memory *Memory) update(epoch, conflictEpoch, newIndex int, newEntries []Log) {

	if epoch != conflictEpoch {
		var idx int

		for idx = len(memory.Epochs) - 1; idx >= 0 && memory.Hashmap[memory.Epochs[idx]] >= newIndex; idx-- {
			delete(memory.Hashmap, memory.Epochs[idx])
		}

		memory.Epochs = memory.Epochs[:idx+1]
	}

	for idx, log := range newEntries {
		if log.Epoch != conflictEpoch {
			memory.add(log.Epoch, newIndex+idx)
			conflictEpoch = log.Epoch
		}
	}
}

func (memory *Memory) print() string {
	result := ""
	for _, epoch := range memory.Epochs {
		result += "{ " + fmt.Sprint(epoch) + " " + fmt.Sprint(memory.Hashmap[epoch]) + " } "
	}
	return result
}
