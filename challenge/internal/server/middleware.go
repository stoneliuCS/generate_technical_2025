package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/ogen-go/ogen/middleware"
)

func logging(logger *slog.Logger) middleware.Middleware {
	return func(req middleware.Request, next func(req middleware.Request) (middleware.Response, error)) (middleware.Response, error) {
		start := time.Now()

		// Extract request information
		operationName := req.OperationName
		operationID := req.OperationID

		// Log the incoming request
		logger.Info("Incoming Request:",
			slog.String("operation", operationName),
			slog.String("operation_id", operationID),
			slog.Time("start_time", start),
		)

		// Call the next handler
		resp, err := next(req)

		// Calculate duration
		duration := time.Since(start)

		// Log based on response/error
		if err != nil {
			// Log error case
			logger.Info("request failed",
				slog.String("operation", operationName),
				slog.String("operation_id", operationID),
				slog.Duration("duration", duration),
				slog.Any("error", err),
			)
		} else {
			logger.Info("request completed",
				slog.String("operation", operationName),
				slog.String("operation_id", operationID),
				slog.Duration("duration", duration),
			)
		}

		return resp, err
	}
}

func slackErrorMiddleware(slackWebhookURI string) middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		messageHeader := "Runtime Error"
		defer func() {
			if r := recover(); r != nil {
				panicErr := fmt.Errorf("panic: %v", r)
				slackAlertError(panicErr, messageHeader, slackWebhookURI)
			}
		}()

		// Call the next handler.
		resp, err := next(req)

		if err != nil {
			slackAlertError(err, messageHeader, slackWebhookURI)
		}

		return resp, err
	}
}

func slowRequestMiddleware(threshold time.Duration, slackWebhookURI string) middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		start := time.Now()
		resp, err := next(req)
		duration := time.Since(start)

		if duration > threshold {
			slowErr := fmt.Errorf("slow request: %s took %v (threshold: %v)",
				req.Raw.URL, duration, threshold)
			messageHeader := "Slow Request"
			slackAlertError(slowErr, messageHeader, slackWebhookURI)
		}

		return resp, err
	}
}

func slackAlertError(err error, messageHeader string, slackWebhookURI string) {
	// Get stack trace.
	buf := make([]byte, 4096)
	stackSize := runtime.Stack(buf, false)
	stack := string(buf[:stackSize])

	message := fmt.Sprintf("%s\n```%v\n\nStack Trace:\n%s```", messageHeader, err, stack)

	payload := map[string]string{"text": message}
	jsonData, _ := json.Marshal(payload)

	go http.Post(slackWebhookURI, "application/json", bytes.NewBuffer(jsonData))
}
