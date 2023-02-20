package routes

import (
	"github.com/gorilla/mux"
)

func RouteInit(r *mux.Router) {
	UserRoutes(r)
	AuthRoutes(r)
	OrderRoutes(r)
	ProductRoutes(r)
	TransactionRoutes(r)
	UserTrcRoutes(r)
}
