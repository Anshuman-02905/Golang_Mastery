package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Order struct {
	ID     int32
	Status string
}

func main() {
	orders := generateOrders(20)
	processOrders(orders)
	updateorderStatus(orders)
	reportOrderStatus(orders)

	fmt.Println("All operations completed . Exiting")
}

func processOrders(orders []*Order) {
	for _, order := range orders {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)

		fmt.Printf("Processing order %d\n", order.ID)
	}
}

func updateorderStatus(orders []*Order) {
	for _, order := range orders {
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		status := []string{
			"Processsing", "Shipped", "Delivered",
		}[rand.Intn(3)]
		order.Status = status
		fmt.Printf(
			"Updated order %d status: %s\n",
			order.ID, status,
		)

	}
}

func reportOrderStatus(orders []*Order) {
	for i := 0; i < 5; i++ {
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
