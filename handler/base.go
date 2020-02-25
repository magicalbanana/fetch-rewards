package handler

import (
	"net/http"
	"strings"
	"time"

	"github.com/nats-io/nuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Handler represents the handler that will be used for the Base
type Handler func(*http.Request) (*Response, error)

// Base is a custom http.Handler ...
type Base struct {
	H      Handler
	Logger *zap.Logger
}

// ServeHTTP allows our Base type to satisfy http.Handler.
func (b Base) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	startingTime := time.Now().UTC()

	zapFields := []zapcore.Field{
		zap.String("event", "web request"),
		zap.String("method", req.Method),
		zap.String("path", req.URL.Path),
		zap.String("query", req.URL.RawQuery),
	}

	b.Logger.Info("Request Started", zapFields...)

	response, err := b.H(req)
	if err != nil {
		b.Logger.Error("Request Error", zap.String("error", err.Error()))
		if response == nil {
			response = NewJSONResponse()
		}
		if response.AppCode == "" {
			response.AppCode = "app.internal_server_error"
		}
		if response.HTTPStatus == 0 {
			response.HTTPStatus = http.StatusInternalServerError
		}
		respErr := &Error{
			ID:     nuid.Next(),
			Code:   "app.internal_server_error",
			Detail: "internal server error",
		}
		response.AddError(respErr)
	}

	rw.WriteHeader(response.HTTPStatus)

	respByte, err := response.Marshal()
	if err != nil {
		b.Logger.Error("Marshal Error", zap.String("error", err.Error()))
		respByte = []byte(strings.Join([]string{
			"An internal server error occured, error_id: ",
			nuid.Next(),
		}, " "),
		)
	}
	_, err = rw.Write(respByte)
	if err != nil {
		b.Logger.Error("Writing Response", zap.String("error", err.Error()))
	}

	// Calculate duration of this request
	reqDuration := time.Now().UTC().Sub(startingTime)

	zapFields = append(zapFields, zap.Duration("request_time", reqDuration))
	zapFields = append(zapFields, zap.Int("response_code", response.HTTPStatus))
	b.Logger.Info("Request Completed", zapFields...)
}
