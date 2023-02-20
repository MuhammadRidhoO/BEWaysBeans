package routes

import (
	"BEWaysBeans/handlers"
	"BEWaysBeans/pkg/middleware"
	"BEWaysBeans/pkg/mysql"
	"BEWaysBeans/repositories"

	"github.com/gorilla/mux"
)

func UserRoutes(r *mux.Router) {	
	userRepository := repositories.RepositoryUser(mysql.DB)
	h := handlers.HandlerUser(userRepository)

	r.HandleFunc("/users", middleware.Auth(h.FindUsers)).Methods("GET")
	r.HandleFunc("/user/{id}", middleware.Auth(h.GetUser)).Methods("GET")
	r.HandleFunc("/profile", middleware.Auth(h.GetUserLogin)).Methods("GET")
	r.HandleFunc("/user-image/{id}", middleware.Auth(middleware.UploadFile(h.UpdateUser,"image"))).Methods("PATCH")
	r.HandleFunc("/user-password/{id}", middleware.Auth(h.UpdatePasswordUser)).Methods("PATCH")
	r.HandleFunc("/user/{id}", h.DeleteUser).Methods("DELETE")
}
