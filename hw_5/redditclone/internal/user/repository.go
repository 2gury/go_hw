package user

import "lectures-2022-1/05_web_app/99_hw/redditclone/internal/models"

type UserRepository interface {
	IsUsernameExist(username string) bool
	InsertUser(user models.User) string
	CheckPassword(user models.User) (string, error)
}