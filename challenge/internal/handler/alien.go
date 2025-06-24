package handler

import (
	"context"
	api "generate_technical_challenge_2025/internal/api"
)

// APIV1AliensGet implements api.Handler.
func (h Handler) APIV1AliensGet(ctx context.Context) (*api.APIV1AliensGetOK, error) {
	return &api.APIV1AliensGetOK{Waves: []api.Alien{}, Budget: 100, Health: 100}, nil
}
