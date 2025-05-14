package routes

import (
	"expense-tracker/handlers"

	"github.com/gorilla/mux"
)

func SetupRoutes() *mux.Router {
	router := mux.NewRouter()

	// Define your routes here
	router.HandleFunc("/expenses", handlers.CreateExpense).Methods("POST")
	router.HandleFunc("/expenses/{id}", handlers.GetExpense).Methods("GET")
	router.HandleFunc("/expenses/{id}", handlers.UpdateExpense).Methods("PUT")
	router.HandleFunc("/expenses/{id}", handlers.DeleteExpense).Methods("DELETE")
	router.HandleFunc("/expenses", handlers.GetAllExpenses).Methods("GET")

	return router
}
