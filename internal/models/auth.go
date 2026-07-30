package models

type RegisterRequest struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Passowrd string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Passowrd string `json:"password" binding:"required"`
}

type LoginResponce struct {
	Token string `json:"token"`
}
