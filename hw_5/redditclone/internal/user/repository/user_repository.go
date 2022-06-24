package repository

import (
	"fmt"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/models"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/user"
	"sync"
)

type UserRepository struct {
	users  map[string]models.User
	lastID uint64
	mx     *sync.Mutex
}

func NewUserRepository() user.UserRepository {
	return &UserRepository{
		users:  map[string]models.User{},
		lastID: 0,
		mx:     &sync.Mutex{},
	}
}

func (r *UserRepository) IsUsernameExist(username string) bool {
	r.mx.Lock()
	defer r.mx.Unlock()

	_, ok := r.users[username]
	return ok
}

func (r *UserRepository) InsertUser(user models.User) string {
	r.mx.Lock()
	defer r.mx.Unlock()

	user.ID = fmt.Sprintf("%d", r.lastID)
	r.users[user.Username] = user
	r.lastID++
	return fmt.Sprintf("%d", r.lastID-1)
}

func (r *UserRepository) CheckPassword(user models.User) (string, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	usr := r.users[user.Username]
	if usr.Password == user.Password {
		return usr.ID, nil
	}
	return "", fmt.Errorf("incorrect password")
}
