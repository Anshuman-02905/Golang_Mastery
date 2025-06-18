package main

import (
	"fmt"
	"time"
)

func main() {

	requests := make(chan int, 11)

	for i := 0; i <= 10; i++ {
		requests <- i
	}

	close(requests)

	limiter := time.NewTicker(200 * time.Millisecond)
	defer limiter.Stop()

	for req := range requests {
		<-limiter.C
		fmt.Println("Processing request ", req, " at ", time.Now())
	}

}
