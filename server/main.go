package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/koosie0507/pluralsight-go-webservices/server/database"
	"github.com/koosie0507/pluralsight-go-webservices/server/product"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
