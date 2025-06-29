package dto

import "github.com/google/uuid"

type Reserve struct {
	OrderId  uuid.UUID        `json:"order_id"`
	Products []ReserveProduct `json:"products"`
}

type ReserveProduct struct {
	Guid     uuid.UUID `json:"guid"`
	Quantity int       `json:"quantity"`
}
