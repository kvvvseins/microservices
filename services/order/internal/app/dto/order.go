package dto

import "github.com/google/uuid"

type ViewOrder struct {
	Number string `json:"number"`
	Price  uint   `json:"price"`
}

type CreateOrder struct {
	Email    string       `json:"email"`
	Products []ProductDto `json:"products"`
	Delivery DeliveryDto  `json:"delivery"`
}

type DeliveryDto struct {
	PlannedDateStart string `json:"planned_date_start"`
	PlannedDateEnd   string `json:"planned_date_end"`
}

type ProductDto struct {
	Guid     uuid.UUID `json:"guid"`
	Quantity int       `json:"quantity"`
}

type Price struct {
	Guid  uuid.UUID `json:"guid"`
	Price uint      `json:"price"`
}
