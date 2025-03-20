package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
	//"time"
)

type Order struct {
	ID     int32
	Status string
	mu     sync.Mutex
}

var (
	totalUpdates int
	updateMutex  int
)

func main() {

	var wg sync.WaitGroup
	wg.Add(3)

	orders := generateOrders(20)

	// go func() {
	// 	defer wg.Done()
	// 	processOrders(orders)
	// }()

	for i := 0; i < 3; i++ {
		go func() {
			defer wg.Done()
			for _, order := range orders {
				updateorderStatus(order)
			}
		}()
	}
	wg.Wait()

	reportOrderStatus(orders)
	fmt.Println(totalUpdates)

	fmt.Println("All operations completed . Exiting")
}

func processOrders(orders []*Order) {
	for _, order := range orders {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

		fmt.Printf("Processing order %d\n", order.ID)
	}
}

func updateorderStatus(order *Order) {
	
	order.mu.Lock()
	time.Sleep(
		time.Duration(rand.Intn(300)) *
			time.Millisecond,
	)

	status := []string{
		"Processsing", "Shipped", "Delivered",
	}[rand.Intn(3)]
	order.Status = status
	fmt.Printf(
		"Updated order %d status: %s\n",
		order.ID, status,
	)
	order.mu.Unlock()

	order.mu.Lock()
	currentUpdates := totalUpdates
	time.Sleep(5 * time.Millisecond)
	totalUpdates = currentUpdates + 1
	order.mu.Unlock()

}

func reportOrderStatus(orders []*Order) {
	time.Sleep(1 * time.Second)
	fmt.Println("\n--ORDER STATUS REPORT--")
	for _, order := range orders {
		fmt.Printf(
			"Order %d: %s \n",
			order.ID, order.Status,
		)
	}
	fmt.Println("------------------------\n")
}

func generateOrders(count int) []*Order {
	orders := make([]*Order, count) // Allocate slice with count elements
	for i := 0; i < count; i++ {    // Loop condition fixed
		orders[i] = &Order{ // Fixed assignment
			ID:     int32(i) + 1,
			Status: "Pending",
		}
	}
	return orders // Return the populated slice
}
