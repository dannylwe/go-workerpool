package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

type Job struct {
	Id int
}

type JobResult struct {
	Output string
}

var (
	workerPool = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "Worker_pool",
		Help: "Get the Total of the current active workers",
	})
)

func init() {
	prometheus.MustRegister(workerPool)
	workerPool.Set(0)
}

func worker(id int, wg *sync.WaitGroup, jobChannel <-chan Job, resutlChannel chan JobResult) {
	defer wg.Done()
	for job := range jobChannel {
		resutlChannel <- doSomething(id, job)
	}
}

func doSomething(workerID int, job Job) JobResult {
	fmt.Printf("Worker #%d Running job #%d\n", workerID, job.Id)
	time.Sleep(time.Millisecond * 500)
	workerPool.Inc()
	return JobResult{Output: "success"}
}

func main() {
	start := time.Now()
	var jobs []Job

	// create jobs
	for i := 0; i < 100; i++ {
		jobs = append(jobs, Job{Id: i})
	}

	const NumberOfWorkers = 10
	var (
		wg               sync.WaitGroup
		jobChannel       = make(chan Job)
		jobResultChannel = make(chan JobResult, len(jobs))
	)
	wg.Add(NumberOfWorkers)

	// start workers
	for j := 0; j < NumberOfWorkers; j++ {

		go worker(j, &wg, jobChannel, jobResultChannel)
	}

	// send jobs to worker
	for _, job := range jobs {
		jobChannel <- job
	}

	close(jobChannel)
	wg.Wait()
	close(jobResultChannel)

	var jobResults []JobResult

	for result := range jobResultChannel {
		// workerPool.Dec()
		jobResults = append(jobResults, result)
	}
	fmt.Println(jobResults)
	fmt.Printf("Took %s\n", time.Since(start))

	http.Handle("/metrics", promhttp.Handler())
	log.Printf("writing on PORT 8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
