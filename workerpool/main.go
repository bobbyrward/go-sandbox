package main

import (
	"fmt"
	"sync"
)

type ThingToDo struct {
	ID int
}

type Result struct {
	ID int
}

func worker(workerId int, wg *sync.WaitGroup, tasks chan ThingToDo, results chan Result) {
	for task := range tasks {
		fmt.Printf("Finished work on %d\n", task.ID)
		results <- Result{ID: task.ID}
	}

	fmt.Printf("Worker %d: Done with work\n", workerId)

	wg.Done()

	fmt.Printf("Worker %d: exiting\n", workerId)
}

func queueWork(numberOfJobs int, tasks chan ThingToDo) {
	for i := 0; i < numberOfJobs; i++ {
		fmt.Printf("Queueing %d\n", i)
		tasks <- ThingToDo{ID: i}
	}
	close(tasks)
	fmt.Printf("Done queueing. Exiting\n")
}

func accumulateResults(done chan bool, results chan Result) {
	idx := 0

	for result := range results {
		fmt.Printf("Result %d, Count %d\n", result.ID, idx)
		idx += 1
	}

	fmt.Printf("Result goroutine done. Signaling\n")

	done <- true

	fmt.Printf("Result goroutine exiting\n")
}

func processTasks(numberOfWorkers int, tasks chan ThingToDo, results chan Result) {
	var wg sync.WaitGroup

	for i := 0; i < numberOfWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg, tasks, results)
	}

	fmt.Printf("Waiting on WG\n")
	wg.Wait()

	close(results)

}

func main() {
	numberOfJobs := 100
	numberOfWorkers := 4

	done := make(chan bool)
	tasks := make(chan ThingToDo, numberOfWorkers)
	results := make(chan Result, numberOfWorkers)

	go accumulateResults(done, results)
	go queueWork(numberOfJobs, tasks)

	processTasks(numberOfWorkers, tasks, results)

	fmt.Printf("Waiting on done chan\n")
	<-done

	fmt.Printf("Done\n")
}
