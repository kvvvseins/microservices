package dto

type ViewOrder struct {
	Number string `json:"number"`
}

type CreateOrder struct {
	Products []ProductDto `json:"products"`
}

type ProductDto struct {
	Guid  string `json:"guid"`
	Price int    `json:"price"` // Как появится сервис склада цена уберется
}
