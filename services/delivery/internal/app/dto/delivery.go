package dto

import (
	"github.com/google/uuid"
)

type ViewDelivery struct {
	UserId           uuid.UUID `json:"user_id"`
	OrderId          uuid.UUID `json:"order_id"`
	PlannedDateStart string    `json:"planned_date_start"`
	PlannedDateEnd   string    `json:"planned_date_end"`
}

type CreateDelivery struct {
	OrderId          uuid.UUID `json:"order_id"`
	PlannedDateStart string    `json:"planned_date_start"`
	PlannedDateEnd   string    `json:"planned_date_end"`
}
