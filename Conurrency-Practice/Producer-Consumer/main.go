package main

import (
	"fmt"
	"sync"
	"time"
)

func Consumer(id int, ch <-chan int, wg *sync.WaitGroup) {
	defer wg.Done()
	for val := range ch {
		fmt.Printf("Consumer %d workking on data %d \n", id, val)
		time.Sleep(200 * time.Millisecond)
		fmt.Printf("Consumer %d work on %d  is done \n", id, val)
	}
}
func Producer(id, start, end int, ch chan<- int, wg *sync.WaitGroup) {
	defer wg.Done()

	for ; start < end; start++ {
		fmt.Printf("Producer %d is making the produce value %d \n", id, start)
		ch <- start
		time.Sleep(200 * time.Millisecond)
	}

}

func main() {

	ch := make(chan int, 5)

	numOfConsumer := 5
	numOfProducers := 3
	numofjobs := 15

	var wg sync.WaitGroup
	var Pwg sync.WaitGroup
	wg.Add(numOfConsumer)
	for i := 1; i <= numOfConsumer; i++ {
		go Consumer(i, ch, &wg)
	}

	for i := 0; i < numOfProducers; i++ {
		start := i*numofjobs + 1
		end := (i + 1) * numofjobs
		Pwg.Add(1)
		go Producer(i, start, end, ch, &Pwg)
	}

	Pwg.Wait()
	close(ch)
	wg.Wait()

	fmt.Println("All consumer Processing Finished")
}
