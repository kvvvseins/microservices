package dto

type ViewUser struct {
	Email string `json:"email"`
}

type CreateUser struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type UpdateUser struct {
	Email string `json:"email"`
}
