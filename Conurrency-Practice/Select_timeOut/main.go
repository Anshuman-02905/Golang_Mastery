package main

import (
	"fmt"
	"time"
)

func main() {

	c := make(chan string)
	go func() {
		time.Sleep(3 * time.Second)
		c <- "Hello"

	}()

	select {
	case msg := <-c:
		fmt.Println(msg)
	case <-time.After(4 * time.Second):
		fmt.Println("TIME OUT")
	}
}
