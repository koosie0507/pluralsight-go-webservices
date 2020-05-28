package product

import (
	"fmt"

	"github.com/koosie0507/pluralsight-go-webservices/server/database"
)

func getProduct(id int) (*Product, error) {
	row := database.DbConnection.QueryRow(`
SELECT productId, manufacturer, sku, upc, pricePerUnit, quantityOnHand, productName
FROM products
WHERE productId = ?`, id)
	if row == nil {
		return nil, fmt.Errorf("failed to fetch product [%d] from the db", id)
	}
	var product Product
	err := row.Scan(
		&product.ProductID,
		&product.Manufacturer,
		&product.SKU,
		&product.UPC,
		&product.PricePerUnit,
		&product.QuantityOnHand,
		&product.ProductName)
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func removeProduct(id int) (int64, error) {
	result, err := database.DbConnection.Exec("DELETE FROM products WHERE productID=?", id)
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func getProductList() ([]Product, error) {
	results, err := database.DbConnection.Query(`
SELECT productId, manufacturer, sku, upc, pricePerUnit, quantityOnHand, productName
FROM products`)
	if err != nil {
		return nil, err
	}
	defer results.Close()
	products := make([]Product, 0)
	for results.Next() {
		var product Product
		results.Scan(
			&product.ProductID,
			&product.Manufacturer,
			&product.SKU,
			&product.UPC,
			&product.PricePerUnit,
			&product.QuantityOnHand,
			&product.ProductName)
		products = append(products, product)
	}
	return products, nil
}

func updateProduct(product Product) (int64, error) {
	result, err := database.DbConnection.Exec(`
UPDATE products
SET manufacturer=?, sku=?, upc=?, pricePerUnit=?, quantityOnHand=?, productName=?
WHERE productID=?`,
		product.Manufacturer,
		product.SKU,
		product.UPC,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName,
		product.ProductID)
	if err != nil {
		return 0, err
	}
	count, err := result.RowsAffected()
	if err != nil {
		return 0, err
	}
	return count, nil
}

func insertProduct(product Product) (int64, error) {
	result, err := database.DbConnection.Exec(`
INSERT INTO products (manufacturer, sku, upc, pricePerUnit, quantityOnHand, productName)
VALUES (?, ?, ?, ?, ?, ?)`,
		product.Manufacturer,
		product.SKU,
		product.UPC,
		product.PricePerUnit,
		product.QuantityOnHand,
		product.ProductName)
	if err != nil {
		return 0, err
	}
	insertID, err := result.LastInsertId()
	if err != nil {
		return 0, err
	}
	return insertID, nil
}
