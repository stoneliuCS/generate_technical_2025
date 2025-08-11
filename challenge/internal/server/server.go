package server

import (
	"fmt"
	api "generate_technical_challenge_2025/internal/api"
	"generate_technical_challenge_2025/internal/utils"
	"log/slog"
	"net/http"
	"time"
)

// Runs the server api with the given handler.
func RunServer(handler api.Handler, cfg utils.EnvConfig, logger *slog.Logger) {
	// Create middleware for logging.
	opts := api.WithMiddleware(
		logging(logger),
		slackErrorMiddleware(cfg.SLACK_WEBHOOK),
		// I would set this lower than 10 seconds, but the ngrok challenge is slow because it
		// has to make many http requests to the ngrok server, averaging 6-7 seconds.
		slowRequestMiddleware(10*time.Second, cfg.SLACK_WEBHOOK))

	// Create server
	srvFunc := func() (*api.Server, error) { return api.NewServer(handler, opts) }
	srv := utils.FatalCall(srvFunc)
	addr := fmt.Sprintf(":%d", cfg.PORT)
	servFunc := func() error {
		logger.Info("Started server on http://localhost" + addr)
		return http.ListenAndServe(addr, srv)
	}

	// Run server indefinitely
	utils.FatalCallErrorSupplier(servFunc)
}
