package main

import (
	"fmt"
	"sync"
)

func main() {

	var val = []string{"abc", "def", "ghi"}
	var wg sync.WaitGroup
	
	for _, v := range val {
		wg.Add(1)
		go func( v2 string) {
			defer wg.Done()
			fmt.Println(v2)
		}(v)
	}
	wg.Wait()
}
