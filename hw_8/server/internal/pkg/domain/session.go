package domain

import (
	"context"
	"errors"
	"net/http"
)

var ErrNoSession = errors.New("no session")

type Session struct {
	UserID string
}

type SessionService interface {
	CheckSession(ctx context.Context, headers http.Header) (Session, error)
}
