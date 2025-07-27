// package main

// import (
// 	"fmt"
// 	"time"
// )

// //Wait for the first of multiple cancellation signals (i.e., "whichever done channel closes first").

// //OrDone combines a data channel and a done channel
// //it stops delivering values if done is closed

// func orDone(done <-chan struct{}, ch <-chan int) <-chan int {
// 	out := make(chan int)

// 	go func() {
// 		defer close(out)
// 		for {
// 			select {
// 			case <-done:
// 				return
// 			case val, ok := <-ch:
// 				if !ok {
// 					return
// 				}
// 				select {
// 				case out <- val:
// 				case <-done:
// 				}

// 			}
// 		}
// 	}()
// 	return out
// }

// func producer(ch chan<- int) {
// 	for i := 1; i <= 10; i++ {
// 		ch <- i
// 		time.Sleep(500 * time.Millisecond)
// 	}
// 	close(ch)
// }

// func consumer(done <-chan struct{}, data <-chan int) {
// 	for val := range orDone(done, data) {
// 		fmt.Println("Recieve:", val)
// 		time.Sleep(1 * time.Second)
// 	}
// 	fmt.Println("Consumer stopped")
// }

// func or(channels ...<-chan struct{}) <-chan struct{} {
// 	switch len(channels) {
// 	case 0:
// 		return nil
// 	case 1:
// 		return channels[0]
// 	case 2:
// 		out := make(chan struct{})
// 		go func() {
// 			defer close(out)
// 			select {
// 			case <-channels[0]:
// 			case <-channels[1]:
// 			}
// 		}()
// 		return out
// 	default:
// 		first := channels[0]
// 		rest := or(channels[1:]...)
// 		return or(first, rest)
// 	}
// }

// func main() {
// 	// data := make(chan int)
// 	// done := make(chan struct{})
// 	// go producer(data)
// 	// go consumer(done, data)
// 	// time.Sleep(3 * time.Second)
// 	// close(done)

// 	// time.Sleep(1 * time.Second)
// 	// fmt.Println("Main Exited")
// }

package main

import (
	"fmt"
	"time"
)

func or(channels ...<-chan any) <-chan any {
	switch len(channels) {
	case 0:
		return nil
	case 1:
		return channels[0]
	case 2:
		orDone := make(chan any)
		go func() {
			defer close(orDone)
			select {
			case <-channels[0]:
			case <-channels[1]:

			}
		}()
		return orDone
	default:
		orDone := make(chan any)
		go func() {
			defer close(orDone)
			select {
			case <-channels[0]:
			case <-channels[1]:
			case <-or(append(channels[2:], orDone)...):
			}
		}()
		return orDone

	}
}

func worker(id int, orDone <-chan any) {
	fmt.Printf("Worker %d : started \n", id)
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-orDone:
			fmt.Printf("Worder %d : Received cancel signal exiting \n", id)
			return
		case <-tick.C:
			fmt.Printf("Worker %d : working...\n", id)
		}
	}
}

func signal(after time.Duration, name string) <-chan any {
	c := make(chan any)
	go func() {
		time.Sleep(after)
		fmt.Printf(">>> Signal %s fired after %v \n", name, after)
		close(c)
	}()
	return c
}

func main() {
	timeoutSignal := signal(4*time.Second, "timeout")
	manualSignal := signal(5*time.Second, "Manual Cancel")
	shutDownHook := signal(6*time.Second, "shutdown")
	timeoutSignal1 := signal(7*time.Second, "timeout1")
	manualSignal1 := signal(8*time.Second, "Manual Cancel1")
	shutDownHook1 := signal(9*time.Second, "shutdown1")
	orDone := or(timeoutSignal, manualSignal, shutDownHook,timeoutSignal1, manualSignal1, shutDownHook1)
	go worker(1, orDone)
	go worker(2, orDone)
	<-orDone
	fmt.Println("MAIN received cancellation signal , all workers hsould exit")
	time.Sleep(1 * time.Second)
}
