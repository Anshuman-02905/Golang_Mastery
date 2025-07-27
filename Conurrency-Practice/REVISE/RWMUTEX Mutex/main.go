package main

import (
	"fmt"
	"sync"
	"time"
)

// RWMutex Mutex

var num = 10
var mu sync.Mutex
var rwmu sync.Mutex

func Change() {
	mu.Lock()
	time.Sleep(3 * time.Second)
	num += 1
	mu.Unlock()

}

func Read() {
	rwmu.Lock()
	time.Sleep(1 * time.Second)
	fmt.Println(num)
	rwmu.Unlock()

}

func main() {
	
	for i :=range(100){

	}

}
