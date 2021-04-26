package raft

import "sync/atomic"

type Clock struct {
	value int32
}

func (clock *Clock) Set(value int) {
	atomic.StoreInt32(&(clock.value), int32(value))
}

func (clock *Clock) Increment() {
	atomic.AddInt32(&(clock.value), 1)
}

func (clock *Clock) Get() int {
	return int(atomic.LoadInt32(&(clock.value)))
}
