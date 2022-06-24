package models

type Vote struct {
	Vote int64  `json:"vote"`
	User uint64 `json:"user,string"`

	ID     uint64 `json:"-"`
	UserID uint64 `json:"-"`
	PostID uint64 `json:"-"`
}
