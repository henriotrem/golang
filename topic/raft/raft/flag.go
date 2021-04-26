package raft

import "sync/atomic"

type Flag struct{ flag int32 }

func (b *Flag) Set(value bool) {
	var i int32 = 0
	if value {
		i = 1
	}
	atomic.StoreInt32(&(b.flag), int32(i))
}

func (b *Flag) Get() bool {
	return atomic.LoadInt32(&(b.flag)) != 0
}
