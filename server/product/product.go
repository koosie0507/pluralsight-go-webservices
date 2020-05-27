package product

// Product instances describe products.
type Product struct {
	ProductID      int    `json:"productId"`
	Manufacturer   string `json:"manufacturer"`
	SKU            string `json:"sku"`
	UPC            int    `json:"upc"`
	PricePerUnit   string `json:"pricePerUnit"`
	QuantityOnHand int    `json:"quantityOnHand"`
	ProductName    string `json:"productName"`
}
