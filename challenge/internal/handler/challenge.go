package handler

import (
	"context"
	api "generate_technical_challenge_2025/internal/api"
	"generate_technical_challenge_2025/internal/services"

	"github.com/samber/lo"
)

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
func (h Handler) APIV1ChallengeBackendIDAliensSubmitPost(ctx context.Context, req api.OptAPIV1ChallengeBackendIDAliensSubmitPostReq, params api.APIV1ChallengeBackendIDAliensSubmitPostParams) (api.APIV1ChallengeBackendIDAliensSubmitPostRes, error) {
	panic("unimplemented")
}

// APIV1ChallengeFrontendIDAliensGet implements api.Handler.
func (h Handler) APIV1ChallengeFrontendIDAliensGet(ctx context.Context, params api.APIV1ChallengeFrontendIDAliensGetParams) (api.APIV1ChallengeFrontendIDAliensGetRes, error) {
	panic("unimplemented")
}
