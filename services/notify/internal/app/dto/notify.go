package dto

type ViewNotify struct {
	Email   string `json:"email"`
	Message string `json:"message"`
}

type CreateNotify struct {
	Email string                 `json:"email"`
	Type  string                 `json:"type"`
	Data  map[string]interface{} `json:"data"`
}
