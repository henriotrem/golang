package raft

import "sync/atomic"

type Counter struct {
	value int32
}

func (counter *Counter) Reset() {
	atomic.StoreInt32(&(counter.value), 0)
}

func (counter *Counter) Increment() {
	atomic.AddInt32(&(counter.value), 1)
}

func (counter *Counter) Get() int {
	return int(atomic.LoadInt32(&(counter.value)))
}
