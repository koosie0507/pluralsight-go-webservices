package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"

	"github.com/koosie0507/pluralsight-go-webservices/server/middleware"
)

const productsBasePath = "products"

//SetupRoutes is a utility function for setting up the products API
func SetupRoutes(apiBasePath string) {
	productListHandler := http.HandlerFunc(productsHandler)
	productItemHandler := http.HandlerFunc(productHandler)

	http.Handle(
		fmt.Sprintf("%s/%s", apiBasePath, productsBasePath),
		middleware.Log(middleware.JSON(middleware.CORS(productListHandler))),
	)
	http.Handle(
		fmt.Sprintf("%s/%s/", apiBasePath, productsBasePath),
		middleware.Log(middleware.JSON(middleware.CORS(productItemHandler))),
	)
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	product, err := getProduct(productID)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	switch r.Method {
	case http.MethodGet:
		productJSON, err := json.Marshal(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(productJSON)
	case http.MethodPut:
		var updatedProduct Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		err = json.Unmarshal(bodyBytes, &updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if updatedProduct.ProductID != product.ProductID {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		count, err := updateProduct(updatedProduct)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else if count == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		w.WriteHeader(http.StatusAccepted)
		return
	case http.MethodDelete:
		count, err := removeProduct(productID)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
		} else if count == 0 {
			w.WriteHeader(http.StatusNotFound)
		}
		w.WriteHeader(http.StatusOK)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productList, err := getProductList()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		productsJSON, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Write(productsJSON)
	case http.MethodPost:
		var product Product
		bodyBytes, err := ioutil.ReadAll(r.Body)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		err = json.Unmarshal(bodyBytes, &product)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
		}
		if product.ProductID != 0 {
			w.WriteHeader(http.StatusBadRequest)
		}
		_, err = insertProduct(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
	case http.MethodOptions:
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}
