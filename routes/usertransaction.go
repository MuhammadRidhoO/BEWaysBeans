package routes

import (
	handlers "BEWaysBeans/handlers"
	"BEWaysBeans/pkg/mysql"
	"BEWaysBeans/repositories"

	"github.com/gorilla/mux"
)

func UserTrcRoutes(r *mux.Router) {
	usertrcRepository := repositories.RepositoryUserTrc(mysql.DB)
	h := handlers.HandlerUsertrc(usertrcRepository)

	r.HandleFunc("/user-transactions/{id}", h.FindUserTrc).Methods("GET")

}
