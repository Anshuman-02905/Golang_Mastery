package main

import (
	"fmt"
	"time"
)

//In a leaky goroutine, the return is either absent or logically unreachable
//(e.g., stuck on a blocking operation like a channel read). This causes the goroutine to never exit,
//even after main ends or the parent context is done.

func leakyGoroutine(ch <-chan int) {
	for val := range ch {
		fmt.Println("Recived", val)
	}
}
func nonLeakyGoroutine(ch <-chan int, done <-chan struct{}) {
	for {
		select {
		case val := <-ch:
			fmt.Println(val)
		case <-done:
			fmt.Println("Goroutine recived Done Channel . Exiting")
			return //Clean exit no leak
		}
	}
}
func main() {
	ch := make(chan int)
	done := make(chan struct{})
	go leakyGoroutine(ch)
	go nonLeakyGoroutine(ch, done)
	ch <- 1

	time.Sleep(1 * time.Second)
	close(done)
	time.Sleep(1 * time.Second)
	fmt.Println("Main exits without closing channel")
}
