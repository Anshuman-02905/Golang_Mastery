package main

/*
import (
	"fmt"
	"time"
)


func printMessage() {
	fmt.Println("Hello from Goroutine")
}

func main() {
	go printMessage()
	fmt.Println("Hello from main function")
	// Wait for the Goroutine to finish
	time.Sleep(time.Second)
}*/

// import (
// 	"fmt"
// 	"time"
// )

// func printMessage(message chan string) {
// 	time.Sleep(time.Second * 2)
// 	message <- "Hello from Goroutine"
// }

// func main() {
// 	message := make(chan string)
// 	go printMessage(message)
// 	fmt.Println("Hello from main function")
// 	fmt.Println(<-message)
// }

import (
	"fmt"
	"time"
)

func worker(id int, jobs <-chan int, results chan<- int) {
	for j := range jobs {
		fmt.Println("Worker", id, "processing job", j)
		time.Sleep(time.Second)
		results <- j * 2
	}
}

func main() {
	const numJobs = 5
	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	for w := 1; w <= 3; w++ {
		go worker(w, jobs, results)
	}

	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	for a := 1; a <= numJobs; a++ {
		fmt.Println("Result:", <-results)
	}
}
