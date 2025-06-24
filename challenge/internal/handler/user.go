// Package for user registrations during technical challenge.
package handler

import (
	"context"
	api "generate_technical_challenge_2025/internal/api"
	models "generate_technical_challenge_2025/internal/database/models"
	"strings"
)

// APIV1RegisterPost implements api.Handler.
func (h Handler) APIV1RegisterPost(ctx context.Context, req api.OptAPIV1RegisterPostReq) (api.APIV1RegisterPostRes, error) {
	email := req.Value.GetEmail()
	nuid := req.Value.GetNuid()
	splitEmail := strings.SplitN(email, "@", 2)
	if splitEmail[1] != "northeastern.edu" {
		return &api.APIV1RegisterPostBadRequest{Message: "Not a valid northeastern email address."}, nil
	}
	if len(nuid) != 9 {
		return &api.APIV1RegisterPostBadRequest{Message: "Not a valid NUID."}, nil
	}
	// Deserialize input into internal model of users.
	user := models.CreateUser(email, nuid)
	id, err := h.userService.CreateUser(user)
	if err != nil {
		return &api.APIV1RegisterPostInternalServerError{Message: "Database error when creating a new user."}, err
	}
	return &api.APIV1RegisterPostCreated{ID: *id}, nil
}
