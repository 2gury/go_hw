package tools

import (
	"context"
	"fmt"
	"os"
	"runtime"

	"github.com/labstack/echo"
)

func CurrentFunction() string {
	counter, _, _, success := runtime.Caller(1)

	if !success {
		println("functionName: runtime.Caller: failed")
		os.Exit(1)
	}

	return runtime.FuncForPC(counter).Name()
}

func ConvertEchoContext(echoCtx echo.Context) context.Context {
	echoCtx.Get("request_id")
	return context.WithValue(echoCtx.Request().Context(), "request_id", echoCtx.Get("request_id"))
}

func GetRequestID(ctx context.Context) (string, error) {
	requestID, ok := ctx.Value("request_id").(string)
	if !ok {
		return "", fmt.Errorf("err get request_id")
	}
	return requestID, nil
}
