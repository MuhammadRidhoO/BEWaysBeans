package routes

import (
	"BEWaysBeans/handlers"
	"BEWaysBeans/pkg/middleware"
	"BEWaysBeans/pkg/mysql"
	"BEWaysBeans/repositories"

	"github.com/gorilla/mux"
)

func OrderRoutes(r *mux.Router) {
	orderRepository := repositories.RepositoryOrder(mysql.DB)
	h := handlers.HandlerOrder(orderRepository)

	// find orders
	r.HandleFunc("/orders", middleware.Auth(h.FindOrders)).Methods("GET")

	// get 1 order
	r.HandleFunc("/order/{id}", middleware.Auth(h.GetOrder)).Methods("GET")

	// add order
	r.HandleFunc("/order", middleware.Auth(h.Create_Order_Product)).Methods("POST")

	// update order
	r.HandleFunc("/order/{id}", middleware.Auth(h.UpdateOrder)).Methods("PATCH")

	// delete order
	r.HandleFunc("/order/{id}", middleware.Auth(h.DeleteOrder)).Methods("DELETE")

	// delete all order
	r.HandleFunc("/orderallorder", middleware.Auth(h.DeleteAll)).Methods("DELETE")
}
