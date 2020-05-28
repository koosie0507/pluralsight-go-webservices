package main

import (
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/koosie0507/pluralsight-go-webservices/server/database"
	"github.com/koosie0507/pluralsight-go-webservices/server/product"
	"github.com/koosie0507/pluralsight-go-webservices/server/receipt"
)

const apiBasePath = "/api"

func main() {
	database.SetupDatabase()
	product.SetupRoutes(apiBasePath)
	receipt.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
