package models

type Comment struct {
	Author    User   `json:"author"`
	Body      string `json:"body"`
	CreatedAt string `json:"created"`
	ID        string `json:"id"`
}
