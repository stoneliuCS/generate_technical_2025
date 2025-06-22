package server

import (
	"log/slog"
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
