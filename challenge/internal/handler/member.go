// Package for user registrations during technical challenge.
package handler

import (
	"context"
	api "generate_technical_challenge_2025/internal/api"
	models "generate_technical_challenge_2025/internal/database/models"
	"strings"
)

func validateNEUEmail(email string) bool {
	splitEmail := strings.SplitN(email, "@", 2)
	return splitEmail[1] == "northeastern.edu"
}

func validateNUID(nuid string) bool {
	return len(nuid) == 9
}

// APIV1MemberGet implements api.Handler.
func (h Handler) APIV1MemberGet(ctx context.Context, params api.APIV1MemberGetParams) (api.APIV1MemberGetRes, error) {
	if !validateNEUEmail(params.Email) {
		return &api.APIV1MemberGetBadRequest{Message: "Not a valid northeastern email address"}, nil
	}
	if !validateNUID(params.Nuid) {
		return &api.APIV1MemberGetBadRequest{Message: "Not a valid NUID"}, nil
	}
	id, err := h.memberService.GetMember(params.Email, params.Nuid)
	if err != nil {
		return &api.APIV1MemberGetNotFound{Message: "Could not find a northeastern email address or nuid associated."}, nil
	}
	return &api.APIV1MemberGetOK{ID: *id}, nil
}

// APIV1MemberRegisterPost implements api.Handler.
func (h Handler) APIV1MemberRegisterPost(ctx context.Context, req api.OptAPIV1MemberRegisterPostReq) (api.APIV1MemberRegisterPostRes, error) {
	email := req.Value.GetEmail()
	nuid := req.Value.GetNuid()
	// Validate EMAIL and NUID
	if !validateNEUEmail(email) {
		return &api.APIV1MemberRegisterPostBadRequest{Message: "Not a valid northeastern email address."}, nil
	}
	if !validateNUID(nuid) {
		return &api.APIV1MemberRegisterPostBadRequest{Message: "Not a valid NUID."}, nil
	}
	exists, err := h.memberService.CheckMemberExists(email, nuid)
	if err != nil {
		return &api.APIV1MemberRegisterPostInternalServerError{Message: "Database error querying for user."}, nil
	}
	if exists {
		return &api.APIV1MemberRegisterPostConflict{Message: "Member already exists."}, nil
	}
	// Deserialize input into internal model of users.
	member := models.CreateMember(email, nuid)
	id, err := h.memberService.CreateMember(member)
	if err != nil {
		return &api.APIV1MemberRegisterPostInternalServerError{Message: "Database error when creating a new user."}, err
	}
	return &api.APIV1MemberRegisterPostCreated{ID: *id}, nil
}
