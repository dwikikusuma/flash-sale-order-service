package entity

type Order struct {
	ID              int64          `json:"id"`
	UserID          int64          `json:"user_id"`
	ProductRequests []OrderRequest `json:"product_requests"` // List of products in the order
	Quantity        int            `json:"quantity"`
	TotalPrice      float64        `json:"total_price"`
	Status          string         `json:"status"` // e.g., "pending", "completed", "cancelled"
	HashValue       string         `json:"hash_value"`
}

type OrderRequest struct {
	ProductID  int64   `json:"product_id"`
	Quantity   int64   `json:"quantity"`
	MarkUp     float64 `json:"markup"`      // Percentage markup on the product price
	Discount   float64 `json:"discount"`    // Percentage discount on the product price
	FinalPrice float64 `json:"final_price"` // Final price after applying markup and discount
	OrderID    int64   `json:"order_id"`
	HashValue  string  `json:"hash_value"`
}

type AvailabilityChannel struct {
	ProductID int64
	Available bool
	Error     error
}
