package main

import (
	"fmt"
	"sync"
	"time"
)

type Job struct {
	Id int
}

type JobResult struct {
	Output string
}

func worker(id int, wg *sync.WaitGroup, jobChannel <-chan Job) {
	defer wg.Done()
	for job := range jobChannel {
		doSomething(id, job)
	}
}

func doSomething(workerID int, job Job) JobResult {
	fmt.Printf("Worker #%d Running job #%d\n", workerID, job.Id)
	time.Sleep(time.Second * 1)
	return JobResult{Output: "success"}
}

func main() {
	start := time.Now()
	var jobs []Job

	for i := 0; i < 100; i++ {
		jobs = append(jobs, Job{Id: i})
	}

	const NumberOfWorkers = 10
	var (
		wg         sync.WaitGroup
		jobChannel = make(chan Job)
	)
	wg.Add(NumberOfWorkers)

	// start workers
	for j := 0; j < NumberOfWorkers; j++ {
		go worker(j, &wg, jobChannel)
	}

	// send jobs to worker
	for _, job := range jobs {
		jobChannel <- job
	}

	close(jobChannel)
	wg.Wait()
	fmt.Printf("Took %s\n", time.Since(start))
}
