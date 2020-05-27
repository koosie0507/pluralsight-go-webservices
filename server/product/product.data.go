package product

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"sync"
)

var productMap = struct {
	sync.RWMutex
	m map[int]Product
}{m: make(map[int]Product)}

func init() {
	fmt.Println("loading products ...")
	products, err := loadProductMap()
	if err != nil {
		log.Fatal(err)
	}
	productMap.m = products
	fmt.Printf("%d products loaded ...\n", len(productMap.m))
}

func loadProductMap() (map[int]Product, error) {
	fileName := "products.json"
	file, err := ioutil.ReadFile(fileName)
	if err != nil {
		return nil, fmt.Errorf("unable to read from '%s': %s", fileName, err)
	}
	productList := make([]Product, 0)
	err = json.Unmarshal([]byte(file), &productList)
	if err != nil {
		return nil, fmt.Errorf("unable to parse JSON from file '%s': %s", fileName, err)
	}
	m := make(map[int]Product)
	for _, product := range productList {
		m[product.ProductID] = product
	}
	return m, nil
}

func getProduct(id int) *Product {
	productMap.RLock()
	defer productMap.RUnlock()
	if product, ok := productMap.m[id]; ok {
		return &product
	}
	return nil
}

func removeProduct(id int) {
	productMap.Lock()
	defer productMap.Unlock()
	delete(productMap.m, id)
}

func getProductList() []Product {
	productMap.RLock()
	products := make([]Product, 0, len(productMap.m))
	for _, value := range productMap.m {
		products = append(products, value)
	}
	productMap.RUnlock()
	return products
}

func getNextProductID() int {
	maxID := 0
	for _, product := range productMap.m {
		if maxID < product.ProductID {
			maxID = product.ProductID
		}
	}
	return maxID + 1
}

func putProduct(product Product) (int, error) {
	newID := -1
	if product.ProductID > 0 {
		oldProduct := getProduct(product.ProductID)
		if oldProduct == nil {
			return 0, fmt.Errorf("product id '%d' doesn't exist", product.ProductID)
		}
		newID = product.ProductID
	} else {
		newID = getNextProductID()
		product.ProductID = newID
	}
	productMap.Lock()
	productMap.m[newID] = product
	productMap.Unlock()
	return newID, nil
}
