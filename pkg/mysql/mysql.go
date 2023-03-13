package mysql

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	// "gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func DatabaseInit() {

	var PGHOST = os.Getenv("PGHOST")
	var PGUSER = os.Getenv("PGUSER")
	var PGPASSWORD = os.Getenv("PGPASSWORD")
	var PGDATABASE = os.Getenv("PGDATABASE")
	var PGPORT = os.Getenv("PGPORT")
	var err error
	// dsn := "root:@tcp(127.0.0.1:3306)/waysbeans?charset=utf8mb4&parseTime=True&loc=Local"
	// DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", PGHOST, PGUSER, PGPASSWORD, PGDATABASE, PGPORT)
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	fmt.Println("Connected to Database")
}

// client/src/config/api.js

// package mysql

// import (
// 	"fmt"

// 	"gorm.io/driver/mysql"
// 	"gorm.io/gorm"
// )

// var DB *gorm.DB

// func DatabaseInit() {
// 	var err error
// 	dsn := "root:@tcp(127.0.0.1:3306)/waysbeans?charset=utf8mb4&parseTime=True&loc=Local"
// 	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})

// 	if err != nil {
// 		panic(err)
// 	}

// 	fmt.Println("Connected to Database")
// }

// client/src/config/api.js
