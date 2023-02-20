package routes

import (
	"BEWaysBeans/handlers"
	"BEWaysBeans/pkg/middleware"
	"BEWaysBeans/pkg/mysql"
	"BEWaysBeans/repositories"

	"github.com/gorilla/mux"
)

func ProductRoutes(r *mux.Router) {
	productRepository := repositories.RepositoryProduct(mysql.DB)
	h := handlers.HandlerProduct(productRepository)

	r.HandleFunc("/products",h.FindProducts).Methods("GET")
	r.HandleFunc("/filterproducts",h.FilterProducts).Methods("GET")
	r.HandleFunc("/product/{id}",h.GetProduct).Methods("GET")
	r.HandleFunc("/product", middleware.Auth(middleware.UploadFile(h.CreateProduct, "image_product"))).Methods("POST")
	r.HandleFunc("/product/{id}", middleware.Auth(middleware.UploadFile(h.UpdateProduct, "image_product"))).Methods("PATCH")
}
