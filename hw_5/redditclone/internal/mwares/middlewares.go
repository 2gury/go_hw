package mwares

import (
	"context"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/consts"
	customContext "lectures-2022-1/05_web_app/99_hw/redditclone/internal/helpers/context"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/helpers/errors"
	customErrors "lectures-2022-1/05_web_app/99_hw/redditclone/internal/helpers/errors"
	"lectures-2022-1/05_web_app/99_hw/redditclone/tools/response"
	"lectures-2022-1/05_web_app/99_hw/redditclone/tools/session"
	"log"
	"net/http"
	"time"
)

func PanicCoverMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ctx := r.Context()
				customErr := errors.Get(consts.CodeInternalError)
				response.WriteErrorResponse(w, ctx, customErr)
			}
		}()
		next.ServeHTTP(w, r)
	})
}

func CheckAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		authToken := r.Header.Get("Authorization")
		user, err := session.CheckSession(session.GetTokenValue(authToken))
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

func AccessLogMiddleware(next http.Handler) http.Handler {
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
