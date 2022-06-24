package models

type Vote struct {
	UserID string `json:"user"`
	Vote   int64  `json:"vote"`
}
