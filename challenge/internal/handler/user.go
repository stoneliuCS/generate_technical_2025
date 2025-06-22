// Package for user registrations during technical challenge.
package handler

import (
	"context"
	api "generate_technical_challenge_2025/internal/api"
)

// Handles user registration via their northeastern email.
func (h Handler) APIV1RegisterGet(ctx context.Context, req api.OptAPIV1RegisterGetReq) (*api.APIV1RegisterGetCreated, error) {
	// email := req.Value.GetEmail()
	// nuid := req.Value.GetNuid()
	panic("")
}
