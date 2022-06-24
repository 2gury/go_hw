package main

import (
	"fmt"
	"server/internal/api/middleware"
	"server/internal/pkg/comment/handler"
	commentrepo "server/internal/pkg/comment/repository"
	commentsvc "server/internal/pkg/comment/service"
	"server/internal/pkg/session"
	"server/internal/pkg/http_client"
	threadhttp "server/internal/pkg/thread/handler"
	threadrepo "server/internal/pkg/thread/repository"
	threadsvc "server/internal/pkg/thread/service"

	"github.com/labstack/echo"
)

func main() {
	e := echo.New()
	monitoring := middleware.NewMonitoring(e)
	httpCli := http_client.NewHttpClient(monitoring)
	sessionSvc := session.NewService(httpCli)

	mwares := middleware.NewMiddlwareManager(sessionSvc, monitoring)

	e.Use(mwares.AccessLogMiddleware())

	threadRepo := threadrepo.NewRepository(httpCli)
	threadSvc := threadsvc.NewService(threadRepo)
	threadHandler := threadhttp.Handler{ThreadSvc: threadSvc}

	commentRepo := commentrepo.NewRepository(httpCli)
	commentSvc := commentsvc.NewService(commentRepo, threadRepo)
	commentHandler := handler.Handler{CommentSvc: commentSvc}

	e.GET("/thread/:tid", threadHandler.GetThread, mwares.AuthEchoMiddleware())
	e.POST("/thread", threadHandler.CreateThread, mwares.AuthEchoMiddleware())
	e.POST("/thread/:tid/comment", commentHandler.Create, mwares.AuthEchoMiddleware())
	e.POST("/thread/:tid/comment/:cid/like", commentHandler.Like, mwares.AuthEchoMiddleware())
	e.GET("/500", commentHandler.InternalServerError)

	fmt.Print(e.Start(":8000"))
}
