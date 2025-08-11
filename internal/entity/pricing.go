package entity

type Pricing struct {
	ProductID  int64   `json:"product_id"`  // ID of the product
	MarkUp     float64 `json:"markup"`      // Percentage markup on the product price
	Discount   float64 `json:"discount"`    // Percentage discount on the product price
	FinalPrice float64 `json:"final_price"` // Final price after applying markup and discount
}
