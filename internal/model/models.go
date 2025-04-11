package model

// Product представляет продукт в системе
// swagger:model Product
type Product struct {
	// Название продукта
	// required: true
	// example: Product Name
	Name string `json:"name"`

	// Описание продукта
	// example: Description of the product
	Description string `json:"description"`

	// Цена продукта
	// required: true
	// example: 19.99
	Price float64 `json:"price"`
}
