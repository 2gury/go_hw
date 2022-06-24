package repository

import (
	"encoding/json"
	"fmt"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/models"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/session"

	"github.com/gomodule/redigo/redis"
)

type SessionRdRepository struct {
	rdConn redis.Conn
}

func NewSessionRdRepository(conn redis.Conn) session.SessionRepository {
	return &SessionRdRepository {
		rdConn: conn,
	}
}

func (r *SessionRdRepository) Create(session *models.Session) error {
	sess, err := json.Marshal(session)
	if err != nil {
		return err
	}

	res, err := redis.String(r.rdConn.Do("SET", session.Value, sess, "EX", session.GetTime()))
	if err != nil {
		return err
	}
	if res != "OK" {
		return fmt.Errorf("redis: not OK")
	}

	return nil
}

func (r *SessionRdRepository) Get(sessValue string) (*models.Session, error) {
	sess := &models.Session{}

	bytes, err := redis.Bytes(r.rdConn.Do("GET", sessValue))
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(bytes, sess)
	if err != nil {
		return nil, err
	}

	return sess, nil
}

func (r *SessionRdRepository) Delete(sessValue string) error {
	_, err := redis.Int(r.rdConn.Do("DEL", sessValue))
	if err != nil {
		return err
	}

	return nil
}