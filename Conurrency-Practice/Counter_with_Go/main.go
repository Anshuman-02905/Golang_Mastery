// Counter with Goroutines
// Spawn 1000 goroutines that each increment a shared counter.
// Fix the race condition using sync.Mutex.
// 🧠 Concept: Race condition, Mutex

package main

import (
	"fmt"
	"runtime"
	"sync"
	"time"
)

type Counter struct {
	mu      sync.Mutex
	counter int
}

func main() {

	c := Counter{
		counter: 0,
	}

	for i := 0; i < 1000; i++ {

		go func() {
			for {
				c.mu.Lock()
				fmt.Println(c.counter, " ")
				fmt.Println("Number of GoRoutines", runtime.NumGoroutine())
				time.Sleep(1 * time.Second)
				c.counter++
				c.mu.Unlock()
			}

		}()
	}

	select {}
}


// Great question — you're exploring the trade-offs between **concurrent vs sequential execution** at scale. Let's address your two points:

// ---

// ## 🧪 **1. Removing `time.Sleep` with a million goroutines**

// ### ✅ Will output still be consistent?

// **Yes**, if you keep using `sync.Mutex`, the output will still be **consistent and free of race conditions** because:

// * The `mutex` ensures only **one goroutine modifies/prints** at a time.
// * Go’s runtime and scheduler will **queue and execute** each goroutine eventually.

// ### ❌ But you’ll run into **real-world limitations**:

// | Issue                  | Description                                                                                                                    |
// | ---------------------- | ------------------------------------------------------------------------------------------------------------------------------ |
// | **Memory pressure**    | Each goroutine consumes \~2KB of stack memory (grows as needed). 1M goroutines = \~2GB+ memory minimum.                        |
// | **Scheduler overhead** | The Go scheduler handles goroutine scheduling, but with millions, context switching and queueing become slower.                |
// | **CPU bottleneck**     | All goroutines contend for the same mutex. This creates a **sequential bottleneck**, which defeats the purpose of concurrency. |

// ---

// ## 🚫 **2. What if you do the same sequentially (no goroutines)?**

// If you write a loop like this instead:

// ```go
// for i := 0; i < 1000000; i++ {
// 	fmt.Println(i)
// }
// ```

// ### 🔸 Characteristics:

// * **No concurrency** = no race conditions
// * **Fully deterministic** and **ordered**
// * **Much faster per operation** because there’s no context switching or locking

// ### 🔸 But:

// * **You lose parallelism** — on multicore CPUs, it doesn’t scale
// * **Cannot handle I/O waiting** effectively (e.g., real-world workloads that block)

// ---

// ## ⚖️ Summary: Concurrent vs Sequential

// | Feature           | Concurrent (Goroutines + Mutex)  | Sequential                       |
// | ----------------- | -------------------------------- | -------------------------------- |
// | Consistency       | ✅ (with mutex)                   | ✅                                |
// | Speed (with lock) | ❌ slow due to lock contention    | ✅ fast                           |
// | Scalability       | ✅ up to a limit (e.g., \~1M)     | ❌ single-threaded                |
// | Memory use        | ⚠️ high (stack + scheduler)      | ✅ low                            |
// | Best for          | Blocking I/O, scalable workloads | Simple counters, CPU-bound tasks |

// ---

// Would you like to see how to do **batched goroutines** (e.g., only run 1K at a time) to avoid spawning all million at once?
