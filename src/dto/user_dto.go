package dto

type SignupData struct {
	Name     string `json:"name" bindings:"required"`
	Email    string `json:"email" bindings:"required"`
	Password string `json:"password" bindings:"required"`
}

type LoginData struct {
	Email    string `json:"email" bindings:"required"`
	Password string `json:"password" bindings:"required"`
}
