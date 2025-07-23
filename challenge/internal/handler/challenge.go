package handler

import (
	"context"
	api "generate_technical_challenge_2025/internal/api"
	"generate_technical_challenge_2025/internal/services"
	"generate_technical_challenge_2025/internal/utils"
	"net/url"

	"github.com/samber/lo"
)

var globalRateLimiter = utils.NewRateLimiter()

// APIV1ChallengeBackendIDAliensGet implements api.Handler.
func (h Handler) APIV1ChallengeBackendIDAliensGet(ctx context.Context, params api.APIV1ChallengeBackendIDAliensGetParams) (api.APIV1ChallengeBackendIDAliensGetRes, error) {
	exists, err := h.memberService.CheckMemberExistsById(params.ID)
	if err != nil {
		return &api.APIV1ChallengeBackendIDAliensGetInternalServerError{Message: "Database error finding member Id."}, nil
	}
	if !exists {
		return &api.APIV1ChallengeBackendIDAliensGetNotFound{Message: "Unable to find member id."}, nil
	}
	waves := h.challengeService.GenerateUniqueAlienChallenge(params.ID)
	states := lo.Map(waves, func(state services.InvasionState, _ int) api.APIV1ChallengeBackendIDAliensGetOKItem {
		alienMap := lo.Map(state.SurveyRemainingAlienInvasion(), func(alien services.Alien, _ int) api.APIV1ChallengeBackendIDAliensGetOKItemAliensItem {
			return api.APIV1ChallengeBackendIDAliensGetOKItemAliensItem{Hp: alien.Hp, Atk: alien.Atk}
		})
		return api.APIV1ChallengeBackendIDAliensGetOKItem{Aliens: alienMap, Hp: state.GetHpLeft()}
	})
	result := api.APIV1ChallengeBackendIDAliensGetOKApplicationJSON(states)
	return &result, nil
}

// APIV1ChallengeBackendIDAliensSubmitPost implements api.Handler.
func (h Handler) APIV1ChallengeBackendIDAliensSubmitPost(ctx context.Context, req [][]api.APIV1ChallengeBackendIDAliensSubmitPostReqItemItem, params api.APIV1ChallengeBackendIDAliensSubmitPostParams) (api.APIV1ChallengeBackendIDAliensSubmitPostRes, error) {
	exists, err := h.memberService.CheckMemberExistsById(params.ID)
	if err != nil {
		return &api.APIV1ChallengeBackendIDAliensSubmitPostInternalServerError{Message: "Database error finding member Id."}, nil
	}
	if !exists {
		return &api.APIV1ChallengeBackendIDAliensSubmitPostNotFound{Message: "Unable to find member id."}, nil
	}

	uuidStr := params.ID.String()
	if !globalRateLimiter.Allow(uuidStr) {
		return &api.APIV1ChallengeBackendIDAliensSubmitPostTooManyRequests{
			Message: "Rate limit exceeded: 10 requests per minute per challenge ID",
		}, nil
	}

	panic("TO BE CREATED")
}

// APIV1ChallengeFrontendIDAliensGet implements api.Handler.
// Note:
// generates a random number of aliens between LOWER_DETAILED_ALIEN_AMOUNT and UPPER_DETAILED_ALIEN_AMOUNT, and then
// limits/offsets it.
func (h Handler) APIV1ChallengeFrontendIDAliensGet(ctx context.Context, params api.APIV1ChallengeFrontendIDAliensGetParams) (api.APIV1ChallengeFrontendIDAliensGetRes, error) {
	exists, err := h.memberService.CheckMemberExistsById(params.ID)
	if err != nil {
		return &api.APIV1ChallengeFrontendIDAliensGetInternalServerError{Message: "Database error finding member Id."}, nil
	}
	if !exists {
		return &api.APIV1ChallengeFrontendIDAliensGetNotFound{Message: "Unable to find member id."}, nil
	}

	detailedAliens := h.challengeService.GenerateUniqueFrontendChallenge(params.ID)

	start := 0
	if params.Offset.Set {
		start = params.Offset.Value
	}

	end := len(detailedAliens)
	if params.Limit.Set {
		end = min(start+params.Limit.Value, len(detailedAliens))
	}

	var pagedAliens []services.DetailedAlien
	if start < len(detailedAliens) {
		pagedAliens = detailedAliens[start:end]
	}

	colony := lo.Map(pagedAliens, func(alien services.DetailedAlien, _ int) api.APIV1ChallengeFrontendIDAliensGetOKItem {
		profileURLParsed, err := url.Parse(alien.ProfileURL)
		if err != nil {
			profileURLParsed = &url.URL{}
		}

		return api.APIV1ChallengeFrontendIDAliensGetOKItem{
			ID:        alien.ID,
			FirstName: alien.FirstName,
			LastName:  alien.LastName,
			Type:      alien.Type.ToAPI(),
			Stats: api.APIV1ChallengeFrontendIDAliensGetOKItemStats{
				Atk: alien.BaseAlien.Atk,
				Hp:  alien.BaseAlien.Hp,
				Spd: alien.Spd,
			},
			URL: *profileURLParsed,
		}
	})

	response := api.APIV1ChallengeFrontendIDAliensGetOKApplicationJSON(colony)
	return &response, nil
}
