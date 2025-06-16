package dto

type ViewOrder struct {
	Number string `json:"number"`
	Price  uint   `json:"price"`
}

type CreateOrder struct {
	Email    string       `json:"email"`
	Products []ProductDto `json:"products"`
}

type ProductDto struct {
	Guid  string `json:"guid"`
	Price int    `json:"price"` // Как появится сервис склада цена уберется
}
