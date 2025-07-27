package main

import (
	"fmt"
	"time"
)

// repeat emits the same value infinitely into a channel until `done` is closed.

func repeat(done <-chan interface{}, values ...interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)
		for {
			for _, v := range values {
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}

// take takes `n` values from the input channel and then stops.
func take(done <-chan interface{}, n int, in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})

	go func() {
		defer close(out)
		for i := 0; i < n; i++ {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				out <- v
			}
		}
	}()

	return out
}

// sleep adds a delay between reading from the input and writing to the output\

func sleep(done <-chan interface{}, duration time.Duration, in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{})
	go func() {
		defer close(out)

		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				time.Sleep(duration)
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}

	}()
	return out
}
func buffer(done <-chan interface{}, size int, in <-chan interface{}) <-chan interface{} {
	out := make(chan interface{}, size)
	go func() {
		defer close(out)
		for {
			select {
			case <-done:
				return
			case v, ok := <-in:
				if !ok {
					return
				}
				select {
				case <-done:
					return
				case out <- v:
				}
			}
		}
	}()
	return out
}

func main() {
	done := make(chan interface{})
	defer close(done)

	zeros := take(done, 3, repeat(done, 0))
	short := sleep(done, 1*time.Second, zeros)
	buffer := buffer(done, 2, short)
	long := sleep(done, 4*time.Second, buffer)
	pipeline := long
	for v := range pipeline {
		fmt.Println("Received:", v)
	}
}
