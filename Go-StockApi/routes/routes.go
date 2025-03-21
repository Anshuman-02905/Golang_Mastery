package routes

import (
	"stock/middleware"

	"github.com/gorilla/mux"
)

// SetupRoutes initializes and returns a new router with predefined API endpoints.
/*
   Returns:
   - *mux.Router: A configured router instance with all stock-related routes.
*/

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()

	r.HandleFunc("/", middleware.GetHelloStock).Methods("GET", "OPTIONS")
	// Route to fetch a single stock by ID
	r.HandleFunc("/api/stock/{id}", middleware.GetStock).Methods("GET", "OPTIONS")

	// Route to fetch all stocks
	r.HandleFunc("/api/stock/", middleware.GetAllStock).Methods("GET", "OPTIONS")

	// Route to create a new stock entry
	r.HandleFunc("/api/newstock", middleware.CreateStock).Methods("POST", "OPTIONS")

	// Route to update an existing stock entry by ID
	r.HandleFunc("/api/stock/{id}", middleware.UpdateStock).Methods("PUT", "OPTIONS")

	// Route to delete a stock entry by ID
	r.HandleFunc("/api/deletestock/{id}", middleware.DeleteStock).Methods("DELETE", "OPTIONS")

	return r
}
