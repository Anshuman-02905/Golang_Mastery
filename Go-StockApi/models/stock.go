package models

// Stock represents the structure of the stock entity.
type Stock struct {
	StockID int64  `json:"stockid"` // Unique identifier for a stock
	Name    string `json:"name"`    // Name of the stock
	Price   string `json:"price"`   // Price of the stock (string type, consider using float64)
	Company string `json:"company"` // Company associated with the stock
}
