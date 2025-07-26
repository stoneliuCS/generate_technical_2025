package handler

import (
	"context"
	challenge "generate_technical_challenge_2025"
	api "generate_technical_challenge_2025/internal/api"
	"generate_technical_challenge_2025/internal/services"
	"generate_technical_challenge_2025/internal/static"
	"log/slog"
	"strings"

	scalar "github.com/MarceloPetrucio/go-scalar-api-reference"
)

// Handles incoming API requests
type Handler struct {
	memberService    services.MemberService
	challengeService services.ChallengeService
	logger           *slog.Logger // event logger
}

// ChallengeGet implements api.Handler.
func (h Handler) ChallengeGet(ctx context.Context) (api.ChallengeGetRes, error) {
	var buf strings.Builder
	err := static.ChallengePage().Render(ctx, &buf)
	if err != nil {
		return &api.ChallengeGetInternalServerError{Message: "Internal Server Error"}, nil
	}
	return &api.ChallengeGetOK{Data: strings.NewReader(buf.String())}, nil
}

// APIV1ChallengeBackendIDNgrokSubmitPost implements api.Handler.
func (h Handler) APIV1ChallengeBackendIDNgrokSubmitPost(ctx context.Context, req api.OptAPIV1ChallengeBackendIDNgrokSubmitPostReq, params api.APIV1ChallengeBackendIDNgrokSubmitPostParams) (api.APIV1ChallengeBackendIDNgrokSubmitPostRes, error) {
	exists, err := h.memberService.CheckMemberExistsById(params.ID)
	if err != nil {
		return &api.APIV1ChallengeBackendIDNgrokSubmitPostInternalServerError{Message: "Database error finding member Id."}, nil
	}
	if !exists {
		return &api.APIV1ChallengeBackendIDNgrokSubmitPostBadRequest{Message: "Unable to find member id."}, nil
	}

	generatedRequests := h.challengeService.GenerateUniqueNgrokChallenge(params.ID)
	gradeResult := h.challengeService.GradeNgrokServer(req.Value.URL.Value, generatedRequests)

	if gradeResult.Valid {
		// Successful grading:
		result := api.APIV1ChallengeBackendIDNgrokSubmitPostOK{
			Type: api.APIV1ChallengeBackendIDNgrokSubmitPostOK0APIV1ChallengeBackendIDNgrokSubmitPostOK,
			APIV1ChallengeBackendIDNgrokSubmitPostOK0: api.APIV1ChallengeBackendIDNgrokSubmitPostOK0{
				Valid: api.NewOptBool(true),
				Score: api.NewOptInt(gradeResult.Score),
			},
		}
		return &result, nil
	} else {
		// Grading failed:
		result := api.APIV1ChallengeBackendIDNgrokSubmitPostOK{
			Type: api.APIV1ChallengeBackendIDNgrokSubmitPostOK1APIV1ChallengeBackendIDNgrokSubmitPostOK,
			APIV1ChallengeBackendIDNgrokSubmitPostOK1: api.APIV1ChallengeBackendIDNgrokSubmitPostOK1{
				Valid:  api.NewOptBool(false),
				Reason: api.NewOptString(gradeResult.Reason),
			},
		}
		return &result, nil
	}
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
func CreateHandler(logger *slog.Logger, memberService services.MemberService, challengeService services.ChallengeService) api.Handler {
	return Handler{
		memberService,
		challengeService,
		logger,
	}
}
