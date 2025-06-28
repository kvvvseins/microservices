package dto

type Reserve struct {
	OrderId  string           `json:"order_id"`
	Products []ReserveProduct `json:"products"`
}

type ReserveProduct struct {
	Guid     string `json:"guid"`
	Quantity int    `json:"quantity"`
}
