package dto

import "github.com/google/uuid"

type ViewOrder struct {
	Number string `json:"number"`
	Price  uint   `json:"price"`
}

type CreateOrder struct {
	Email    string       `json:"email"`
	Products []ProductDto `json:"products"`
}

type ProductDto struct {
	Guid     uuid.UUID `json:"guid"`
	Quantity int       `json:"quantity"`
}

type Price struct {
	Guid  uuid.UUID `json:"guid"`
	Price uint      `json:"price"`
}
