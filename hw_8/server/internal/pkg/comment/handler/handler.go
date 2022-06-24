package handler

import (
	"log"
	"server/internal/pkg/domain"
	"server/tools"

	"github.com/labstack/echo"
)

type Handler struct {
	CommentSvc domain.CommentService
}

func (h Handler) Create(ctx echo.Context) error {
	var comment domain.Comment

	err := ctx.Bind(&comment)
	if err != nil {
		return err
	}

	tid := ctx.Param("tid")

	context := tools.ConvertEchoContext(ctx)
	err = h.CommentSvc.Create(context, tid, comment)
	log.Printf("%s %s %v", ctx.Get("request_id"), "CommentSvc.Create", err)

	return err
}

func (h Handler) Like(ctx echo.Context) error {
	tid := ctx.Param("tid")
	cid := ctx.Param("cid")

	context := tools.ConvertEchoContext(ctx)
	err := h.CommentSvc.Like(context, tid, cid)
	log.Printf("%s %s %v", ctx.Get("request_id"), "CommentSvc.Like", err)

	return err
}

func (h Handler) InternalServerError(ctx echo.Context) error {
	return ctx.JSON(500, "InternalServerError")
}
