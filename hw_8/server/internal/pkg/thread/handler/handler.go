package handler

import (
	"log"
	"server/internal/pkg/domain"
	"server/tools"

	"github.com/labstack/echo"
)

type Handler struct {
	ThreadSvc domain.ThreadService
}

func (h Handler) GetThread(ctx echo.Context) error {
	context := tools.ConvertEchoContext(ctx)
	tid := ctx.Param("tid")

	t, err := h.ThreadSvc.Get(context, tid)
	if err != nil {
		return err
	}
	log.Printf("%s %s %v", ctx.Get("request_id"), "ThreadSvc.Get", err)

	return ctx.JSON(200, t)
}

func (h Handler) CreateThread(ctx echo.Context) error {
	context := tools.ConvertEchoContext(ctx)
	var thread domain.Thread

	err := ctx.Bind(&thread)
	if err != nil {
		return err
	}

	err = h.ThreadSvc.Create(context, thread)
	if err != nil {
		return err
	}
	log.Printf("%s %s %v", ctx.Get("request_id"), "ThreadSvc.Create", err)

	return ctx.NoContent(200)
}
