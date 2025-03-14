package routes

import (
	"stock/middleware"

	"github.com/gorilla/mux"
)

//GET One stock
//GET ALL stock
//Create STock
//Update Stock
//Delete Stock

func SetupRoutes() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/api/stock/{id}", middleware.GetStock).Methods("GET", "OPTIONS")
	r.HandleFunc("/api/stock/", middleware.GetAllStock).Methods("GET", "OPTIOS")
	r.HandleFunc("/api/newstock", middleware.CreateStock).Methods("POST", "OPTIONS")
	r.HandleFunc("/api/stock/{id}", middleware.UpdateStock).Methods("PUT", "OPTIONS")
	r.HandleFunc("/api/deletestock/{id}", middleware.DeleteStock).Methods("DELETE", "OPTIONS")
	return r
}
