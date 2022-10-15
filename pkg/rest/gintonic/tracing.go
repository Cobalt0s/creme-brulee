package gintonic

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

func Respond(c *gin.Context, span trace.Span, statusCode int, responseBody interface{}) {
	span.SetStatus(
		resolveCode(statusCode),
		bodyToMessage(responseBody),
	)
	c.JSON(statusCode, responseBody)
}

func bodyToMessage(responseBody interface{}) string {
	if responseBody == nil {
		return ""
	}
	msg, err := json.Marshal(responseBody)
	if err != nil {
		msg = []byte("could not marshal")
	}
	return string(msg)
}

func resolveCode(code int) codes.Code {
	if code >= 400 {
		return codes.Error
	}
	// 300 are included as ok as well
	if code >= 200 {
		return codes.Ok
	}
	return codes.Unset
}
