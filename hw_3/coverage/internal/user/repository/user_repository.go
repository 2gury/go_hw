package repository

import (
	"encoding/xml"
	"io/ioutil"
	"os"
	"sync"

	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/internal/models"
	"gitlab.com/mailru-go/lectures-2022-1/03/99_hw/coverage/internal/user"
)

type UserRepository struct {
	xmlFilename string
	mx          *sync.Mutex
}

func NewUserRepository(filename string) user.UserRepository {
	return &UserRepository{
		xmlFilename: filename,
		mx:          &sync.Mutex{},
	}
}

func (r *UserRepository) SelectUsers() (*models.SearchUsers, error) {
	r.mx.Lock()
	defer r.mx.Unlock()

	xmlFile, err := os.Open(r.xmlFilename)
	if err != nil {
		return nil, err
	}
	data, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		return nil, err
	}
	users := &models.SearchUsers{}
	if err := xml.Unmarshal(data, users); err != nil {
		return nil, err
	}

	return users, nil
}
