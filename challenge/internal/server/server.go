package server

import (
	"fmt"
	api "generate_technical_challenge_2025/internal/api"
	"generate_technical_challenge_2025/internal/utils"
	"log/slog"
	"net/http"
)

// Runs the server api with the given handler.
func RunServer(handler api.Handler, cfg utils.EnvConfig, logger *slog.Logger) {
	// Create middleware for logging.
	opts := api.WithMiddleware(
		logging(logger),
		slackErrorMiddleware(cfg.SLACK_WEBHOOK))

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
