package main

import (
	"fmt"
	"sync"
)

func main() {
	hello := func(wg *sync.WaitGroup, id int) {
		defer wg.Done()
		fmt.Printf("Hello There from %d \n", id)
	}

	const numGreetors = 5
	var wg sync.WaitGroup
	for i := range numGreetors {
		wg.Add(1)
		go func() {
			hello(&wg, i)
		}()
	}
	wg.Wait()
}
