package main

import (
	"fmt"
	"net/http"
	"sync"
)

type URLResponses struct {
	success map[string]interface{}
	failure map[string]interface{}
}

// map[interface{}]interface{} // THSI IS TYPE UNSAFE

func GenerateResult(url string, response interface{}, err error) map[interface{}]interface{} {
	result := make(map[interface{}]interface{})
	result["request"] = url
	result["response"] = response
	result["error"] = err
	return result
}

func checkGet(url string, results chan<- map[interface{}]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	resp, err := http.Get(url)
	result := GenerateResult(url, resp, err)
	results <- result

}

func GenerateResponse(urls []string, wg *sync.WaitGroup, result_chan chan<- map[interface{}]interface{}) {
	for _, i := range urls {
		wg.Add(1)
		go checkGet(i, result_chan, wg)
	}

}
func PrintReponse(reuslt_chan <-chan map[interface{}]interface{}, wg *sync.WaitGroup) {
	defer wg.Done()
	for value := range reuslt_chan {
		fmt.Printf("URL: %v \n RESPONSE %v \n ERROR %v ", value["request"], value["response"], value["error"])

	}
}

func main() {
	var wg sync.WaitGroup
	result_chan := make(chan map[interface{}]interface{}, 10)
	urls := []string{"https://www.google.com",
		"https://www.facebook.com",
		"https://www.asdadsa.com", // intentionally invalid}
	}
	GenerateResponse(urls, &wg, result_chan)
	wg.Add(1)
	go PrintReponse(result_chan, &wg)
	go func() {
		wg.Wait()          // Wait for all workers (checkGet and PrintResponse)
		close(result_chan) // Safe to close after all sends are done
	}()
	wg.Wait() // Wait for all workers (checkGet and PrintResponse)
	fmt.Printf("ALL DONE \n")

}
We have error in this