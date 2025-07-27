package main

import (
	"fmt"
	"sync"
)

func even(idx, i int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("EVEN AT index %d IS %d\n", idx, i)
}
func odd(idx, i int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Printf("ODD AT index %d IS %d\n", idx, i)
}

//EACH GOROUTINE IS WORKING ON ITS OWN DATA
//You're capturing idx and val by value (i.e., passed as arguments to the goroutines).
//There’s no shared state being mutated, and no risk of race conditions on arr, idx, or val.
//Each goroutine works on its own arguments only, making it a textbook example of ad hoc confinement.

//Yes — perfect! 👏
//Your understanding is spot-on:
//Lexical confinement means we confine a variable to a goroutine, and communicate with the outside world only through channels.

func StartWokrker() <-chan int {
	ch := make(chan int)
	go func(chan<- int) {
		count := 0 //count can only be accessed by this gorouting
		for {
			ch <- count
			count++
		}
	}(ch)
	return ch
}

func main() {
	var arr []int = []int{1, 2, 3, 4}
	var wg sync.WaitGroup
	wg.Add(4)
	for idx, val := range arr {
		if val%2 == 0 {
			go even(idx, val, &wg)
		} else {
			go odd(idx, val, &wg)
		}
	}
	wg.Wait()
	worker := StartWokrker()
	fmt.Println(<-worker)
	fmt.Println(<-worker)
	fmt.Println(<-worker)

}
