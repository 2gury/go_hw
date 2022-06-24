package models

type User struct {
	ID       uint64 `json:"id,string"`
	Username string `json:"username"`
	Password string `json:"-"`
}
