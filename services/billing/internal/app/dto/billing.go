package dto

type ViewBilling struct {
	Value uint `json:"value"`
}

type UpdateBilling struct {
	Value int    `json:"value"`
	Email string `json:"email"`
}
