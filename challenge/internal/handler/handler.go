package handler

import (
	"context"
	api "generate_technical_challenge_2025/internal/api"
	"generate_technical_challenge_2025/internal/services"
	"log/slog"
	"strings"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
)

var openapiSpec string = "../openapi.json"

// Handles incoming API requests
type Handler struct {
	services *services.Services
	logger   *slog.Logger // event logger
}

// APIV1RegisterGet implements api.Handler.
func (h Handler) APIV1RegisterGet(ctx context.Context, req api.OptAPIV1RegisterGetReq) (*api.APIV1RegisterGetCreated, error) {
	panic("unimplemented")
}

// HealtcheckGet implements api.Handler.
func (h Handler) HealthcheckGet(ctx context.Context) (*api.HealthcheckGetOK, error) {
	return &api.HealthcheckGetOK{Message: api.OptHealthcheckGetOKMessage{Value: "OK", Set: true}}, nil
}

// Creates a new handler for all defined API endpoints
func CreateHandler(logger *slog.Logger, services *services.Services) api.Handler {
	return Handler{
		services,
		logger,
	}
}

func (h Handler) Get(ctx context.Context) (api.GetOK, error) {
	html, err := scalar.ApiReferenceHTML(&scalar.Options{
		SpecURL: openapiSpec,
	})
	if err != nil {
		return api.GetOK{}, err
	}
	return api.GetOK{Data: strings.NewReader(html)}, nil
}
