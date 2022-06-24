package models

import "time"

type Comment struct {
	ID        uint64    `json:"id,string"`
	Author    *User     `json:"author"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created"`

	UserID uint64 `json:"-"`
	PostID uint64 `json:"-"`
}
