package main

import (
	"fmt"
	"net/http"
	"BEWaysBeans/database"
	"BEWaysBeans/pkg/mysql"
	"BEWaysBeans/routes"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/joho/godotenv"
)

func main() {
	firstDate := time.Date(2022, 4, 13, 1, 0, 0, 0, time.UTC)
	secondDate := time.Date(2021, 2, 12, 5, 0, 0, 0, time.UTC)
	difference := firstDate.Sub(secondDate)
          

	// Init godotenv here ...
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	// initial DB
	mysql.DatabaseInit()

	// run migration
	database.RunMigration()

	r := mux.NewRouter()

	routes.RouteInit(r.PathPrefix("/api/v1").Subrouter())

	// biar folder uploads bisa diakses
	r.PathPrefix("/uploads").Handler(http.StripPrefix("/uploads/", http.FileServer(http.Dir("./uploads"))))

	// Initialization "uploads" folder to public here ...

	var AllowedHeaders = handlers.AllowedHeaders([]string{"X-Requested-With", "Content-Type", "Authorization"})
	var AllowedMethods = handlers.AllowedMethods([]string{"GET", "POST", "PUT", "HEAD", "OPTIONS", "PATCH", "DELETE"})
	var AllowedOrigins = handlers.AllowedOrigins([]string{"*"})

	var port = "5000"
	fmt.Printf("Years: %d\n", int64(difference.Hours()/24/365))
	fmt.Println("server running localhost: " + port)
	http.ListenAndServe("localhost:"+port, handlers.CORS(AllowedHeaders, AllowedMethods, AllowedOrigins)(r))
}
