package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/singleflight"
)

// One Request for all duplicate requests

// type Group
// func (g *Group) Do(key string, fn func() (interface{}, error)) (v interface{}, err error, shared bool)
// func (g *Group) DoChan(key string, fn func() (interface{}, error)) <-chan Result
// func (g *Group) Forget(key string)
// type Result

func getData() (interface{}, error) {
	fmt.Println("Getting the data for first_time")
	time.Sleep(500 * time.Millisecond)
	return http.Get("https://google.com")
}
func getDataSingleFlight(g *singleflight.Group, wg *sync.WaitGroup) error {
	defer wg.Done()
	v, err, shared := g.Do("get-google", getData)
	if err != nil {
		return err
	}

	res := v.(*http.Response)
	fmt.Printf("status: %d , shared: %v\n", res.StatusCode, shared)
	return nil

}

func main() {
	var wg sync.WaitGroup
	var g singleflight.Group
	wg.Add(10)
	for i := 0; i < 10; i++ {
		go getDataSingleFlight(&g, &wg)
	}
	wg.Wait()

}

// Sure! Let’s walk through your code **line by line**, explaining how `singleflight.Group` works and how your program ensures **duplicate suppression** using it.

// ---

// ### 📦 Imports

// ```go
// import (
// 	"fmt"
// 	"net/http"
// 	"sync"

// 	"golang.org/x/sync/singleflight"
// )
// ```

// * `fmt`: For printing output.
// * `net/http`: To make HTTP requests.
// * `sync`: For `WaitGroup` to synchronize goroutines.
// * `singleflight`: The magic package that ensures **only one execution for duplicate function calls**.

// ---

// ### 📌 Purpose

// You're simulating **10 goroutines** trying to make the **same HTTP call**, but you want **only one actual HTTP request**, and the others to **reuse the result** — that’s what `singleflight.Group` is for.

// ---

// ### 🧠 `getData` Function

// ```go
// func getData() (interface{}, error) {
// 	return http.Get("https://google.com")
// }
// ```

// This function makes an HTTP GET request to Google and returns the response.
// It matches the required signature for `singleflight.Group.Do()`:

// ```go
// func() (interface{}, error)
// ```

// ---

// ### 🚀 `getDataSingleFlight` Function

// ```go
// func getDataSingleFlight(g *singleflight.Group, wg *sync.WaitGroup) error {
// 	defer wg.Done()
// ```

// * This function is run inside a goroutine.
// * `wg.Done()` is called at the end to signal that this goroutine has completed.

// ---

// ```go
// 	v, err, shared := g.Do("get-google", getData)
// ```

// * `g.Do()` is the core of `singleflight`.
// * `key = "get-google"` → All calls with the **same key** are treated as **duplicates**.
// * `getData` is the actual function to run.
// * If 10 goroutines call `Do("get-google", ...)` at the **same time**, only **one** will run `getData()`, and the result will be **shared** with the others.

// Returns:

// * `v`: The result (an `*http.Response`)
// * `err`: Any error from `getData`
// * `shared`: `true` if the result was shared from another goroutine, not freshly computed.

// ---

// ```go
// 	res := v.(*http.Response)
// 	fmt.Printf("status: %d , shared: %v\n", res.StatusCode, shared)
// 	return nil
// ```

// * Type assert `v` to an `*http.Response` (because we know `getData()` returns that).
// * Print HTTP status and whether this result was **shared**.

// ---

// ### 🧵 `main` Function

// ```go
// func main() {
// 	var wg sync.WaitGroup
// 	var g singleflight.Group
// ```

// * Initialize the `WaitGroup` to track goroutines.
// * Initialize `singleflight.Group` to manage deduplication.

// ---

// ```go
// 	wg.Add(10)
// 	for i := 0; i < 10; i++ {
// 		go getDataSingleFlight(&g, &wg)
// 	}
// 	wg.Wait()
// ```

// * Add 10 "slots" to the waitgroup.
// * Launch 10 goroutines to simulate concurrent calls to `getData`.
// * Wait for all of them to finish.

// ---

// ### ✅ Output

// If all goroutines start nearly at the same time, output might look like:

// ```
// status: 200 , shared: false
// status: 200 , shared: true
// status: 200 , shared: true
// ...
// ```

// Only **one** goroutine actually hits the network (`shared: false`), the rest **reuse** the result (`shared: true`).

// ---

// ### 🧠 How `singleflight` Works Internally (Simplified)

// * It uses a **map of keys** and **a lock**.
// * When `g.Do("key", fn)` is called:

//   * If it's the **first call**, it starts executing `fn()` and stores the ongoing call.
//   * If another goroutine calls with the **same key**, it **waits** for the result of the first.
// * Once the first finishes, all the waiting goroutines **receive the same result**.

// ---

// Let me know if you’d like a diagram or want to experiment with a delay to visualize this better!



type call struct{
	val interface{}
	err error
}

type group struct{
	mu sync.Mutex
	mp map[string]*call
}