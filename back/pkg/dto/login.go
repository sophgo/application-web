package dto

type LoginRequest struct {
	UserName string `json:"userName" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type LogoutRequest struct {
	Token string `json:"token" binding:"required"`
}

type LoginResponse struct {
	Token string `json:"token"`
}
