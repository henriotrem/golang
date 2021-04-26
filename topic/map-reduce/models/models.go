package models

import (
	"os"
	"strconv"
)

//
// Map functions return a slice of KeyValue.
//
type KeyValue struct {
	Key   string
	Value string
}

type MapReduce struct {
	Mapf    func(string, string) []KeyValue
	Reducef func(string, []string) string
}

// for sorting by key.
type ByKey []KeyValue

// for sorting by key.
func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }

type Task struct {
	Id     int
	Time   int64
	Method string
	Args   []string
	State  string
}

//
// RPC definitions.
//
// remember to capitalize all names.
//
type GetTaskArgs struct {
}

type GetTaskReply struct {
	Task Task
}

type UpdateTaskArgs struct {
	Task Task
}

type UpdateTaskReply struct {
}

// Cook up a unique-ish UNIX-domain socket name
// in /var/tmp, for the master.
// Can't use the current directory since
// Athena AFS doesn't support UNIX-domain sockets.
func MasterSock() string {
	s := "/var/tmp/824-mr-"
	s += strconv.Itoa(os.Getuid())
	return s
}
