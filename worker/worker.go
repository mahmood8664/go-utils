package main

import (
	"fmt"
	"sync"
)

const NumberOfWorkers = 6

///////JOB//////////
type Job struct {
	Name string
	Days int
}

func (r Job) print() {
	weeks := r.Days / 7
	days := r.Days % 7
	fmt.Printf("%s has worked %d weeks and %d days in the company\r\n", r.Name, weeks, days)
}

////////JOB///////////

////////DISPATCHER////////
type Dispatcher struct {
	workers []Worker
	wg      *sync.WaitGroup
}

func NewDispatcher() *Dispatcher {
	workers := make([]Worker, NumberOfWorkers)
	wg := sync.WaitGroup{}
	for i := 0; i < NumberOfWorkers; i++ {
		workers[i] = Worker{
			JobCh:    make(chan Job),
			Count:    0,
			Finished: make(chan bool),
			wg:       &wg,
		}
		workers[i].Start()
		wg.Add(1)
	}
	return &Dispatcher{
		workers: workers,
		wg:      &wg,
	}
}

func (r *Dispatcher) print() {
	fmt.Print("\r\nInfo:\r\n")
	fmt.Printf("Workers Count: %d\r\n", len(r.workers))
	for i, worker := range r.workers {
		fmt.Printf("Worker#%d -> %d elements processed\r\n", i, worker.Count)
	}
}

func (r *Dispatcher) Dispatch(input map[string]int) {
	go func() {
		workerNumber := 0
		for name, days := range input {
			r.workers[workerNumber].JobCh <- Job{
				Name: name,
				Days: days,
			}
			workerNumber = (workerNumber + 1) % NumberOfWorkers
		}
		for _, worker := range r.workers {
			worker.Finished <- true
		}
	}()
	r.wg.Wait()
}

///////DISPATCHER///////

///////WORKER///////
type Worker struct {
	JobCh    chan Job
	Count    int
	Finished chan bool
	wg       *sync.WaitGroup
}

func (r *Worker) Start() {
	go func() {
		for {
			select {
			case job := <-r.JobCh:
				job.print()
				r.Count += 1
			case <-r.Finished:
				r.wg.Done()
				return
			}
		}
	}()
}

///////WORKER///////

/////////MAIN//////////
func main() {
	jobMap := map[string]int{"AAAA": 30, "BBBB": 475, "CCCC": 1022, "DDDD": 99, "EEEE": 34, "FFFF": 234, "GGGG": 245, "HHHH": 56765, "IIII": 456, "JJJJ": 56}
	dispatcher := NewDispatcher()
	dispatcher.Dispatch(jobMap)
	dispatcher.print()
}
