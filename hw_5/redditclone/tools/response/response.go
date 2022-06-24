package response

import (
	"context"
	"encoding/json"
	customContext "lectures-2022-1/05_web_app/99_hw/redditclone/internal/helpers/context"
	"lectures-2022-1/05_web_app/99_hw/redditclone/internal/helpers/errors"
	"log"
	"net/http"
)

type Body map[string]interface{}

type Response struct {
	Code         int    `json:"code,omitempty"`
	Error *errors.Error `json:"error,omitempty"`
	Body         *Body  `json:"body,omitempty"`
}

func WriteStatusCode(w http.ResponseWriter, ctx context.Context, statusCode int) {
	w.WriteHeader(statusCode)
	customContext.WriteStatusCodeContext(ctx, statusCode)
}

func WriteErrorResponse(w http.ResponseWriter, ctx context.Context, err *errors.Error) {
	WriteStatusCode(w, ctx, err.HTTPCode)
	cusErr := json.NewEncoder(w).Encode(&Body{
		"message": err.Message,
	})

	if cusErr != nil {
		log.Println(err)
	}
}
