package models

type Session struct {
	Iat  int64 `json:"iat"`
	Exp  int64 `json:"exp"`
	User *User `json:"user"`
}
