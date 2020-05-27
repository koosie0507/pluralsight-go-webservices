package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

const productsBasePath = "products"

func SetupRoutes(apiBasePath string) {
	productListHandler := http.HandlerFunc(productsHandler)
	productItemHandler := http.HandlerFunc(productHandler)

	http.Handle(fmt.Sprintf("%s/%s", apiBasePath, productsBasePath), productListHandler)
	http.Handle(fmt.Sprintf("%s/%s/", apiBasePath, productsBasePath), productItemHandler)
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	product := getProduct(productID)
	if product == nil {
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
		w.Header().Set("Content-Type", "application/json")
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
		putProduct(updatedProduct)
		w.WriteHeader(http.StatusAccepted)
		return
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}
}
func productsHandler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		productList := getProductList()
		productsJSON, err := json.Marshal(productList)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "application/json")
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
		newID, err := putProduct(product)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusCreated)
		w.Write([]byte(strconv.Itoa(newID)))
	}
}
