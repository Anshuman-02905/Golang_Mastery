package main

import "fmt"

func generator(done <-chan interface{}, integers ...int) <-chan int {
	inStream := make(chan int)
	//	defer close(inStream) //this closes the channel immeditely before the goroutine can send data

	go func() {
		defer close(inStream)

		for _, i := range integers {
			select {
			case <-done:
				return
			case inStream <- i:
			}
		}
	}()
	return inStream
}

func multiply(done <-chan interface{}, multiplier int, intstream <-chan int) <-chan int {
	multiplyStream := make(chan int)
	//defer close(multiplyStream) //this closes the channel immeditely before the goroutine can send data
	go func() {
		defer close(multiplyStream)

		for i := range intstream {
			select {
			case <-done:
				return
			case multiplyStream <- i * multiplier:
			}
		}
	}()
	return multiplyStream
}

func add(done <-chan interface{}, additive int, intstream <-chan int) <-chan int {
	additionStream := make(chan int)
	//defer close(additionStream) //this closes the channel immeditely before the goroutine can send data

	go func() {
		defer close(additionStream)

		for i := range intstream {
			select {
			case <-done:
				return
			case additionStream <- i + additive:
			}
		}
	}()
	return additionStream
}

func substract(done <-chan interface{}, substract int, instream <-chan int) <-chan int {
	SubstractionStream := make(chan int)
	//defer close(SubstractionStream) //this closes the channel immeditely before the goroutine can send data

	go func() {
		defer close(SubstractionStream)
		for i := range instream {
			select {
			case <-done:
				return
			case SubstractionStream <- i - substract:
			}
		}
	}()
	return SubstractionStream
}

func main() {
	var integers = []int{1, 2, 3, 4, 5, 6}

	done := make(chan interface{})
	defer close(done)
	pipeline := multiply(done, 2, substract(done, 1, add(done, 1, generator(done, integers...))))

	for v := range pipeline {
		fmt.Println(v)
	}
}
