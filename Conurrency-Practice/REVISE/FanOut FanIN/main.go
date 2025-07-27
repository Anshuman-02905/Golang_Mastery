package main

import (
	"fmt"
	"sync"
	"time"
)

func GetUrls(done <-chan interface{}) <-chan string {
	urlStream := make(chan string, 50)
	count := 1
	go func(count int) {
		defer close(urlStream)
		for {
			time.Sleep(2 * time.Millisecond)
			select {
			case <-done:
				return
			case urlStream <- fmt.Sprintf("URL %d", count):
			}
			count++

		}
	}(count)
	return urlStream
}

func Resize(done <-chan interface{}, id int, images <-chan string) <-chan string {
	ResizedStream := make(chan string, 50)
	go func() {
		defer close(ResizedStream)
		for {
			select {
			case <-done:
				return
			case img, ok := <-images:
				if !ok {
					return
				}
				time.Sleep(3 * time.Millisecond)
				ResizedStream <- fmt.Sprintf("Workder %d: Resized %s", id, img)
			}

		}

	}()
	return ResizedStream
}
func fanIn(done <-chan interface{}, channels ...<-chan string) <-chan string {
	var wg sync.WaitGroup
	multiplexedStream := make(chan string)
	output := func(c <-chan string) {
		defer wg.Done()
		for val := range c {
			select {
			case <-done:
				return
			case multiplexedStream <- val:
			}
		}
	}

	wg.Add(len(channels))
	for _, ch := range channels {
		go output(ch)
	}

	go func() {
		wg.Wait()
		close(multiplexedStream)
	}()

	return multiplexedStream
}

func main() {
	done := make(chan interface{})
	defer close(done)

	urls := GetUrls(done)

	const numWorkers = 3
	var worker []<-chan string
	for i := 1; i <= numWorkers; i++ {
		worker = append(worker, Resize(done, i, urls))
	}
	Resized := fanIn(done, worker...)
	for i := 0; i < 10; i++ {
		fmt.Println(<-Resized)
	}
	time.Sleep(5*time.Second)
}
