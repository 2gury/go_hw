package models

import (
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/consts"
	"time"
)

type Session struct {
	Value        string        `json:"value"`
	TimeDuration time.Duration `json:"time_duration"`
}

func NewSession(sessValue string) *Session {
	expires := consts.ExpiresDuration
	return &Session{
		Value:        sessValue,
		TimeDuration: expires,
	}
}

func (sess *Session) GetTime() int {
	return int(sess.TimeDuration.Seconds())
}
