package dto

type Passwd struct {
	UserName  string `json:"userName" validate:"required"`
	OldPasswd string `json:"oldPasswd" validate:"required"`
	NewPasswd string `json:"newPasswd" validate:"required"`
}
