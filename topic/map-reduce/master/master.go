package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"net/rpc"
	"os"
	"strconv"
	"sync"
	"time"

	"github.com/henriotrem/topic/map-reduce/models"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: mrmaster inputfiles...\n")
		os.Exit(1)
	}

	m := MakeMaster(os.Args[1:], 10)
	for !m.Done() {
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)
}

type Master struct {
	step      string
	tasks     map[string][]models.Task
	remaining int
	mu        sync.Mutex
}

func (m *Master) GetTask(args *models.GetTaskArgs, reply *models.GetTaskReply) error {

	m.mu.Lock()

	for _, task := range m.tasks[m.step] {

		now := time.Now().UnixNano()

		if task.State == "as-yet-unstarted" || (task.State == "in-progress" && now-task.Time > 15000000000 || m.step == "done") {
			m.tasks[m.step][task.Id].State = "in-progress"
			m.tasks[m.step][task.Id].Time = now

			reply.Task = m.tasks[m.step][task.Id]
			break
		}
	}

	m.mu.Unlock()

	return nil
}

func (m *Master) UpdateTask(args *models.UpdateTaskArgs, reply *models.UpdateTaskReply) error {

	m.mu.Lock()

	if m.tasks[args.Task.Method][args.Task.Id].Time == args.Task.Time {

		m.tasks[args.Task.Method][args.Task.Id].State = "completed"
		m.remaining--

		if m.remaining == 0 {
			if m.step == "map" {
				m.step = "reduce"
				m.remaining = len(m.tasks[m.step])
			} else if m.step == "reduce" {
				m.step = "done"
			}
		}
	}

	m.mu.Unlock()

	return nil
}

//
// start a thread that listens for RPCs from worker.go
//
func (m *Master) server() {
	rpc.Register(m)
	rpc.HandleHTTP()
	//l, e := net.Listen("tcp", ":1234")
	sockname := models.MasterSock()
	os.Remove(sockname)
	l, e := net.Listen("unix", sockname)
	if e != nil {
		log.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)
}

//
// main/mrmaster.go calls Done() periodically to find out
// if the entire job has finished.
//
func (m *Master) Done() bool {
	time.Sleep(5 * time.Second)
	return m.step == "done"
}

//
// create a Master.
// main/mrmaster.go calls this function.
// nReduce is the number of reduce tasks to use.
//
func MakeMaster(files []string, nReduce int) *Master {

	m := Master{}

	m.mu.Lock()

	m.tasks = make(map[string][]models.Task)

	m.tasks["map"] = createMapTasks(files, nReduce)
	m.tasks["reduce"] = createReduceTasks(files, nReduce)
	m.tasks["done"] = createDoneTask()

	m.step = "map"
	m.remaining = len(m.tasks[m.step])

	m.mu.Unlock()

	m.server()
	return &m
}

func createMapTasks(files []string, nReduce int) []models.Task {

	tasks := make([]models.Task, len(files))

	for index, filename := range files {
		tasks[index].Id = index
		tasks[index].Method = "map"
		tasks[index].Args = []string{filename, strconv.Itoa(nReduce)}
		tasks[index].State = "as-yet-unstarted"
	}

	return tasks
}

func createReduceTasks(files []string, nReduce int) []models.Task {

	tasks := make([]models.Task, nReduce)

	for index := 0; index < nReduce; index++ {
		tasks[index].Id = index
		tasks[index].Method = "reduce"
		tasks[index].Args = []string{strconv.Itoa(len(files))}
		tasks[index].State = "as-yet-unstarted"
	}

	return tasks
}

func createDoneTask() []models.Task {

	tasks := make([]models.Task, 1)

	tasks[0].Id = 0
	tasks[0].Method = "done"
	tasks[0].State = "as-yet-unstarted"

	return tasks
}
