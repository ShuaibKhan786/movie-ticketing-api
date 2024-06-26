package models


type UserAdminCredentials struct {
	Role string `json:"role"`
	Email string `json:"email"`
	Password string `json:"password"`
}