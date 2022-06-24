package middleware

import (
	"log"
	"math/rand"
	"server/internal/pkg/domain"
	"server/tools"
	"strconv"
	"time"

	"github.com/labstack/echo"
)

type MiddlwareManager struct {
	sessSvc    domain.SessionService
	monitoring *Monitoring
}

func NewMiddlwareManager(service domain.SessionService, mnt *Monitoring) *MiddlwareManager {
	return &MiddlwareManager{
		sessSvc:    service,
		monitoring: mnt,
	}
}

func (mwm *MiddlwareManager) AuthEchoMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			ctx := tools.ConvertEchoContext(context)
			_, err := mwm.sessSvc.CheckSession(ctx, context.Request().Header)
			if err != nil {
				return context.NoContent(401)
			}

			return next(context)
		}
	}
}

func (mwm *MiddlwareManager) AccessLogMiddleware() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(context echo.Context) error {
			context.Set("request_id", RandStringRunes(7))
			start := time.Now()
			next(context)
			timing := time.Since(start)

			path := context.Request().URL.Path
			method := context.Request().Method
			status := strconv.Itoa(context.Response().Status)

			mwm.monitoring.Hits.WithLabelValues(method, path, status).Inc()
			mwm.monitoring.Duration.WithLabelValues(method, path, status).Observe(timing.Seconds())

			log.Printf("%v %s %s %s %s", context.Get("request_id"), method, path, status, timing.String())

			return nil
		}
	}
}

func RandStringRunes(n int) string {
	letterRunes := []rune("1234567890")

	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}
