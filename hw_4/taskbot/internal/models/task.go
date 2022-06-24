package models

type Task struct {
	ID uint64
	Description string
	CreatedBy *User
	AssignTo *User
}