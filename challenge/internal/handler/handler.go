package handler

import (
	"context"
	challenge "generate_technical_challenge_2025"
	api "generate_technical_challenge_2025/internal/api"
	"generate_technical_challenge_2025/internal/services"
	"log/slog"
	"strings"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
)

// Handles incoming API requests
type Handler struct {
	memberService services.MemberService
	logger        *slog.Logger // event logger
}

// Get implements api.Handler.
func (h Handler) Get(ctx context.Context) (api.GetRes, error) {
	html, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: challenge.GetSpecPath(),
	})
	if err != nil {
		return &api.GetInternalServerError{Message: "Error fetching API documentation."}, nil
	}
	return &api.GetOK{Data: strings.NewReader(html)}, nil
}

// HealthcheckGet implements api.Handler.
func (h Handler) HealthcheckGet(ctx context.Context) (*api.HealthcheckGetOK, error) {
	return &api.HealthcheckGetOK{Message: "OK"}, nil
}

// Creates a new handler for all defined API endpoints
func CreateHandler(logger *slog.Logger, memberService services.MemberService) api.Handler {
	return Handler{
		memberService,
		logger,
	}
}
