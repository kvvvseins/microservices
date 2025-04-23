package dto

type RegistrationRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
