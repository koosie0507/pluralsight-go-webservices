package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"
	"strings"
)

type Product struct {
	ProductID      int    `json:"productId"`
	Manufacturer   string `json:"manufacturer"`
	SKU            string `json:"sku"`
	UPC            string `json:"upc"`
	PricePerUnit   string `json:"pricePerUnit"`
	QuantityOnHand int    `json:"quantityOnHand"`
	ProductName    string `json:"productName"`
}

var productList []Product
var productNotFound = errors.New("product not found")

func findByID(id int) (*Product, int, error) {
	for i, prod := range productList {
		if prod.ProductID == id {
			return &prod, i, nil
		}
	}
	return &Product{}, -1, productNotFound
}

func getNextID() int {
	max := 0
	for _, p := range productList {
		if max < p.ProductID {
			max = p.ProductID
		}
	}
	return max + 1
}

func productHandler(w http.ResponseWriter, r *http.Request) {
	urlPathSegments := strings.Split(r.URL.Path, "products/")
	productID, err := strconv.Atoi(urlPathSegments[len(urlPathSegments)-1])
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
	}
	product, idx, err := findByID(productID)
	if err == productNotFound {
		w.WriteHeader(http.StatusNotFound)
		return
	} else if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
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
		product = &updatedProduct
		productList[idx] = *product
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
		product.ProductID = getNextID()
		productList = append(productList, product)
		w.WriteHeader(http.StatusCreated)
	}
}

func main() {
	http.HandleFunc("/products", productsHandler)
	http.HandleFunc("/products/", productHandler)
	http.ListenAndServe(":5000", nil)
}
