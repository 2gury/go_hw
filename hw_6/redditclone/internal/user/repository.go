package user

import "lectures-2022-1/06_databases/99_hw/redditclone/internal/models"

type UserPgRepository interface {
	SelectByUsername(username string) (*models.User, error)
	Insert(usr *models.User) (uint64, error) 
}