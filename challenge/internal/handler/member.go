// Package for user registrations during technical challenge.
package handler

import (
	"context"
	api "generate_technical_challenge_2025/internal/api"
	models "generate_technical_challenge_2025/internal/database/models"
	"strings"
)

// APIV1MemberGet implements api.Handler.
func (h Handler) APIV1MemberGet(ctx context.Context, params api.APIV1MemberGetParams) (api.APIV1MemberGetRes, error) {
	panic("unimplemented")
}

// APIV1MemberRegisterPost implements api.Handler.
func (h Handler) APIV1MemberRegisterPost(ctx context.Context, req api.OptAPIV1MemberRegisterPostReq) (api.APIV1MemberRegisterPostRes, error) {
	email := req.Value.GetEmail()
	nuid := req.Value.GetNuid()
	splitEmail := strings.SplitN(email, "@", 2)
	if splitEmail[1] != "northeastern.edu" {
		return &api.APIV1MemberRegisterPostBadRequest{Message: "Not a valid northeastern email address."}, nil
	}
	if len(nuid) != 9 {
		return &api.APIV1MemberRegisterPostBadRequest{Message: "Not a valid NUID."}, nil
	}
	// Deserialize input into internal model of users.
	member := models.CreateMember(email, nuid)
	id, err := h.memberService.CreateMember(member)
	if err != nil {
		return &api.APIV1MemberRegisterPostInternalServerError{Message: "Database error when creating a new user."}, err
	}
	return &api.APIV1MemberRegisterPostCreated{ID: *id}, nil
}
