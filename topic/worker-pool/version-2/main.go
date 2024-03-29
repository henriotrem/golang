package main

import (
	"fmt"
	"runtime"
	"strconv"
	"sync"
	"time"
)

type Job interface {
	Process()
}

type Worker struct {
	done          *sync.WaitGroup
	readyPool     chan chan Job
	assignedQueue chan Job
	quit          chan bool
}

func NewWorker(readyPool chan chan Job, done *sync.WaitGroup) *Worker {
	return &Worker{
		done:          done,
		readyPool:     readyPool,
		assignedQueue: make(chan Job),
		quit:          make(chan bool),
	}
}

func (w *Worker) Start() {
	w.done.Add(1)
	go func() {
		for {
			w.readyPool <- w.assignedQueue
			select {
			case job := <-w.assignedQueue:
				job.Process()
			case <-w.quit:
				w.done.Done()
				return
			}
		}
	}()
}

func (w *Worker) Stop() {
	w.quit <- true
}

type JobQueue struct {
	internalQueue     chan Job
	readyPool         chan chan Job
	workers           []*Worker
	dispatcherStopped *sync.WaitGroup
	workersStopped    *sync.WaitGroup
	quit              chan bool
}

func NewJobQueue(maxWorkers int) *JobQueue {
	workersStopped := &sync.WaitGroup{}
	readyPool := make(chan chan Job, maxWorkers)
	workers := make([]*Worker, maxWorkers)
	for i := 0; i < maxWorkers; i++ {
		workers[i] = NewWorker(readyPool, workersStopped)
	}
	return &JobQueue{
		internalQueue:     make(chan Job),
		readyPool:         readyPool,
		workers:           workers,
		dispatcherStopped: &sync.WaitGroup{},
		workersStopped:    workersStopped,
		quit:              make(chan bool),
	}
}

func (q *JobQueue) Start() {
	for i := 0; i < len(q.workers); i++ {
		q.workers[i].Start()
	}

	go q.dispatch()
}

func (q *JobQueue) Stop() {
	q.quit <- true
	q.dispatcherStopped.Wait()
}

func (q *JobQueue) dispatch() {
	q.dispatcherStopped.Add(1)
	go func() {
		for {
			select {
			case job := <-q.internalQueue:
				workerChannel := <-q.readyPool
				workerChannel <- job
			case <-q.quit:
				for i := 0; i < len(q.workers); i++ {
					q.workers[i].Stop()
				}
				q.workersStopped.Wait()
				q.dispatcherStopped.Done()
				return
			}
		}
	}()
}

func (q *JobQueue) Submit(job Job) {
	q.internalQueue <- job
}

type TestJob struct {
	ID string
}

func (t *TestJob) Process() {
	fmt.Println("Processing job number : ", t.ID)
	time.Sleep(1 * time.Second)
}

func main() {
	jobQueue := NewJobQueue(1 * runtime.NumCPU())
	jobQueue.Start()

	for i := 0; i < 4*runtime.NumCPU(); i++ {
		jobQueue.Submit(&TestJob{strconv.Itoa(i)})
	}

	jobQueue.Stop()
	fmt.Println("Finished")
}
