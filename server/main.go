package main

import (
	"net/http"

	"github.com/koosie0507/pluralsight-go-webservices/server/product"
)

const apiBasePath = "/api"

func main() {
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
