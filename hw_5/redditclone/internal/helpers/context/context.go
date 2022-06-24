package context

import (
	"context"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/consts"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/helpers/errors"
)

type ContextKey int

const (
	UserID ContextKey = 101 + iota
	Username
	StatusCode
)

func WriteStatusCodeContext(ctx context.Context, code int) {
	statusCode, ok := ctx.Value(StatusCode).(*int)
	if !ok {
		return
	}
	*statusCode = code
}

func GetUserID(ctx context.Context) (string, *errors.Error) {
	userID, ok := ctx.Value(UserID).(string)
	if !ok {
		return "", errors.Get(consts.CodeBadRequest)
	}
	return userID, nil
}

func GetUsername(ctx context.Context) (string, *errors.Error) {
	username, ok := ctx.Value(Username).(string)
	if !ok {
		return "", errors.Get(consts.CodeBadRequest)
	}
	return username, nil
}
