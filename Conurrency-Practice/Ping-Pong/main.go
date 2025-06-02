package main

import (
	"fmt"
	"sync"
	"time"
)

type SafeMux struct {
	mu    sync.Mutex
	state string
}

func PrintPing(Ping chan string, Pong chan string) {
	for {
		fmt.Println(<-Ping)
		time.Sleep(1 * time.Second)

		Pong <- "PONG"
	}
}
func PrintPong(Ping chan string, Pong chan string) {
	for {
		fmt.Println(<-Pong)
		time.Sleep(1 * time.Second)
		Ping <- "PING"

	}
}

func (c *SafeMux) MutexPing() {
	for {
		c.mu.Lock()
		if c.state == "PONG" {
			fmt.Print("PING ")
			c.state = "PING"
		}
		c.mu.Unlock()

		time.Sleep(1 * time.Second)

	}

}
func (c *SafeMux) MutexPong() {

	for {

		c.mu.Lock()
		if c.state == "PING" {
			fmt.Print("PONG ")
			c.state = "PONG"
		}
		c.mu.Unlock()
		time.Sleep(1 * time.Second)

	}

}

func main() {
	Ping := make(chan string)
	Pong := make(chan string)
	go PrintPing(Ping, Pong)
	go PrintPong(Ping, Pong)
	Ping <- "PING"
	select {}

	
	c := SafeMux{
		state: "PONG",
	}

	go c.MutexPing()
	go c.MutexPong()
	select {}

}
