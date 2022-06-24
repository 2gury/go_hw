package session

import "lectures-2022-1/06_databases/99_hw/redditclone/internal/models"

type SessionRepository interface {
	Create(session *models.Session) error
	Get(sessValue string) (*models.Session, error)
	Delete(sessValue string) error
}