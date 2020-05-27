package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/koosie0507/pluralsight-go-webservices/server/product"
)

func middlewareHandler(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("before handler; middleware start:", r.URL.Path, "Method:", r.Method)
		start := time.Now()
		handler.ServeHTTP(w, r)
		fmt.Printf("middleware finished; %s\n", time.Since(start))
	})
}

const apiBasePath = "/api"

func main() {
	product.SetupRoutes(apiBasePath)
	http.ListenAndServe(":5000", nil)
}
