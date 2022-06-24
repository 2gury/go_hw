package mwares

import (
	"context"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/consts"
	customContext "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/context"
	customErrors "lectures-2022-1/06_databases/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/06_databases/99_hw/redditclone/internal/session"
	"lectures-2022-1/06_databases/99_hw/redditclone/tools/response"
	"log"
	"net/http"
	"strings"
	"time"
)

type MiddlewareManager struct {
	sessUse session.SessionUsecase
}

func NewMiddlewareManager(sessionUsecase session.SessionUsecase) *MiddlewareManager {
	return &MiddlewareManager{
		sessUse: sessionUsecase,
	}
}

func (mm *MiddlewareManager) PanicCoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ctx := r.Context()
				customErr := customErrors.Get(consts.CodeInternalError)
				response.WriteErrorResponse(w, ctx, customErr)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func (mm *MiddlewareManager) CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authToken := r.Header.Get("Authorization")
		user, err := mm.sessUse.Check(GetTokenValue(authToken))
		if err != nil {
			response.WriteErrorResponse(w, ctx, customErrors.Get(consts.CodeBadRequest))
			return
		}
		ctx = context.WithValue(ctx,
			customContext.UserID, user.ID,
		)
		ctx = context.WithValue(ctx,
			customContext.Username, user.Username,
		)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (mm *MiddlewareManager) AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		var code int
		ctx = context.WithValue(ctx,
			customContext.StatusCode, &code,
		)
		start := time.Now()
		next.ServeHTTP(w, r.WithContext(ctx))
		log.Printf("Status: %d, Method: %s; URL: %s, Time: %s\n", code, r.Method, r.URL.Path, time.Since(start))
	})
}

func GetTokenValue(rawToken string) string {
	return strings.Split(rawToken, " ")[1]
}
