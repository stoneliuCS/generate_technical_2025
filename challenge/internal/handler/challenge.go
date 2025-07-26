package handler

import (
	"context"
	api "generate_technical_challenge_2025/internal/api"
	"generate_technical_challenge_2025/internal/database/models"
	"generate_technical_challenge_2025/internal/services"
	"generate_technical_challenge_2025/internal/utils"
	"net/url"
	"sort"

	"github.com/google/uuid"
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
	// Sort the keys and index through so that the users get the same order.
	keys := lo.Keys(waves)
	sort.Slice(keys, func(i, j int) bool {
		return keys[i].String() < keys[j].String()
	})
	states := lo.Map(keys, func(key uuid.UUID, _ int) api.APIV1ChallengeBackendIDAliensGetOKItem {
		val := waves[key]
		alienMap := lo.Map(val.SurveyRemainingAlienInvasion(), func(alien services.Alien, _ int) api.APIV1ChallengeBackendIDAliensGetOKItemAliensItem {
			return api.APIV1ChallengeBackendIDAliensGetOKItemAliensItem{Hp: alien.Hp, Atk: alien.Atk}
		})
		return api.APIV1ChallengeBackendIDAliensGetOKItem{ChallengeID: key, Aliens: alienMap, Hp: val.GetHpLeft()}
	})
	result := api.APIV1ChallengeBackendIDAliensGetOKApplicationJSON(states)
	return &result, nil
}

// APIV1ChallengeBackendIDAliensSubmitPost implements api.Handler.
func (h Handler) APIV1ChallengeBackendIDAliensSubmitPost(ctx context.Context, req []api.APIV1ChallengeBackendIDAliensSubmitPostReqItem, params api.APIV1ChallengeBackendIDAliensSubmitPostParams) (api.APIV1ChallengeBackendIDAliensSubmitPostRes, error) {
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
	mapVals := lo.SliceToMap(req, func(userSubmission api.APIV1ChallengeBackendIDAliensSubmitPostReqItem) (uuid.UUID, services.UserChallengeSubmission) {
		var commands []string
		for _, cmd := range userSubmission.State.Commands {
			commands = append(commands, string(cmd))
		}
		return userSubmission.ChallengeID.Value, services.UserChallengeSubmission{Hp: userSubmission.State.RemainingHP, Commands: commands, AliensLeft: userSubmission.State.RemainingAliens}
	})
	ans := h.challengeService.ScoreMemberSubmission(params.ID, mapVals)
	response := &api.APIV1ChallengeBackendIDAliensSubmitPostOK{Valid: ans.Valid, Message: ans.Message}
	if ans.Valid {
		response.Score = api.OptInt{Value: ans.Score, Set: true}
		valid := true
		score := models.CreateScore(params.ID, models.ALGORITHM_CHALLENGE_TYPE, ans.Score, valid)
		_, err := h.memberService.CreateScore(score)
		if err != nil {
			return &api.APIV1ChallengeBackendIDAliensSubmitPostInternalServerError{Message: "Database error when saving a score."}, err
		}
	} else {
		valid := false
		score := models.CreateScore(params.ID, models.ALGORITHM_CHALLENGE_TYPE, models.INVALID_SCORE, valid)
		_, err := h.memberService.CreateScore(score)
		if err != nil {
			return &api.APIV1ChallengeBackendIDAliensSubmitPostInternalServerError{Message: "Database error when saving a score."}, err
		}
	}
	return response, nil
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

// APIV1ChallengeBackendIDNgrokSubmitPost implements api.Handler.
func (h Handler) APIV1ChallengeBackendIDNgrokSubmitPost(ctx context.Context, req api.OptAPIV1ChallengeBackendIDNgrokSubmitPostReq, params api.APIV1ChallengeBackendIDNgrokSubmitPostParams) (api.APIV1ChallengeBackendIDNgrokSubmitPostRes, error) {
	exists, err := h.memberService.CheckMemberExistsById(params.ID)
	if err != nil {
		return &api.APIV1ChallengeBackendIDNgrokSubmitPostInternalServerError{Message: "Database error finding member Id."}, nil
	}
	if !exists {
		return &api.APIV1ChallengeBackendIDNgrokSubmitPostBadRequest{Message: "Unable to find member id."}, nil
	}

	uuidStr := params.ID.String()
	if !globalRateLimiter.Allow(uuidStr) {
		return &api.APIV1ChallengeBackendIDNgrokSubmitPostTooManyRequests{
			Message: "Rate limit exceeded: 10 requests per minute per challenge ID",
		}, nil
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
