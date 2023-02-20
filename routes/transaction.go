package routes

import (
	"BEWaysBeans/handlers"
	"BEWaysBeans/pkg/middleware"
	"BEWaysBeans/pkg/mysql"
	"BEWaysBeans/repositories"

	"github.com/gorilla/mux"
)

func TransactionRoutes(r *mux.Router) {
	transactionRepository := repositories.RepositoryTransaction(mysql.DB)
	h := handlers.HandlerTransaction(transactionRepository)

	r.HandleFunc("/transactions", h.FindTransactions).Methods("GET")
	r.HandleFunc("/transaction/{id}", h.GetTransaction).Methods("GET")
	r.HandleFunc("/transaction-user/{id}", h.FindTransactionsByUser).Methods("GET")
	r.HandleFunc("/transaction", middleware.Auth(h.CreateTransaction)).Methods("POST")
	r.HandleFunc("/transaction-cart", middleware.Auth(h.UpdateTransaction)).Methods("PATCH")
	r.HandleFunc("/notification", h.Notification).Methods("POST")
	// r.HandleFunc("/transaction/{id}", h.DeleteTransaction).Methods("DELETE")
}
