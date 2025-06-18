package main

import (
	"fmt"
	"sync"
	"time"
)

func Worker(id int, jobs <-chan int, results chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()
	for i := range jobs {
		fmt.Printf("Worker %d is starting the work %d \n", id, i)
		time.Sleep(3 * time.Second) //Act like some work
		fmt.Printf("Worker %d has finished the work %d \n", id, i)
		results <- i * i
	}
	return
}

func main() {

	const numWorker = 3
	const numJobs = 5

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	var wg sync.WaitGroup

	for i := 1; i <= 3; i++ {
		wg.Add(1)
		go func(workerID int) {
			Worker(workerID, jobs, results, &wg)
		}(i) //closure Trap
	}
	for i := 1; i <= numJobs; i++ {
		jobs <- i
	}
	close(jobs)
	wg.Wait()

	close(results)

	for result := range results {
		fmt.Printf("Results: %d", result)
	}

}
