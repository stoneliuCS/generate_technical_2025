package server

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"runtime"
	"time"

	"github.com/google/uuid"
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
		requestID := getRequestID(req)
		requestBody := getRequestBodyJSON(req)

		defer func() {
			if r := recover(); r != nil {
				panicErr := fmt.Errorf("panic: %v", r)
				slackAlertError(panicErr, messageHeader, slackWebhookURI, requestID, requestBody, req.OperationName)
			}
		}()

		// Call the next handler.
		resp, err := next(req)

		if err != nil {
			slackAlertError(err, messageHeader, slackWebhookURI, requestID, requestBody, req.OperationName)
		}

		return resp, err
	}
}

func slowRequestMiddleware(threshold time.Duration, slackWebhookURI string) middleware.Middleware {
	return func(req middleware.Request, next middleware.Next) (middleware.Response, error) {
		start := time.Now()
		resp, err := next(req)
		duration := time.Since(start)
		requestID := getRequestID(req)
		requestBody := getRequestBodyJSON(req)

		if duration > threshold {
			slowErr := fmt.Errorf("slow request: %s took %v (threshold: %v)",
				req.Raw.URL, duration, threshold)
			messageHeader := "Slow Request"
			slackAlertError(slowErr, messageHeader, slackWebhookURI, requestID, requestBody, req.OperationName)
		}

		return resp, err
	}
}

func slackAlertError(err error, messageHeader string, slackWebhookURI string, requestID string, requestBody string, operationName string) {
	// Get stack trace.
	buf := make([]byte, 8192)
	stackSize := runtime.Stack(buf, false)
	stack := string(buf[:stackSize])

	message := fmt.Sprintf("%s\n```Operation: %s\nRequesting User ID: %s\nTimestamp: %s\n\nRequest Body:\n%s\n\nError: %v\n\nStack Trace:\n%s```",
		messageHeader, operationName, requestID, time.Now().Format(time.RFC3339), requestBody, err, stack)

	payload := map[string]string{"text": message}
	jsonData, _ := json.Marshal(payload)

	http.Post(slackWebhookURI, "application/json", bytes.NewBuffer(jsonData))
}

// Extracts user ID for logging in the slack middleware.
func getRequestID(req middleware.Request) string {
	for paramKey, paramValue := range req.Params {
		if paramKey.Name == "id" && paramKey.In == "path" {
			if uuid, ok := paramValue.(uuid.UUID); ok {
				return uuid.String()
			}
		}
	}

	return "couldn't find user ID"
}

// Extracts request body for logging in the slack middleware.
func getRequestBodyJSON(req middleware.Request) string {
	if req.Body == nil {
		return "no body"
	}

	jsonData, err := json.MarshalIndent(req.Body, "", "  ")
	if err != nil {
		return fmt.Sprintf("failed to marshal body: %v", err)
	}

	return string(jsonData)
}
