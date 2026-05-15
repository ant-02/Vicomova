package query

import (
	"context"

	commerceRepo "vicomova/internal/commerce/domain/repository"
	"vicomova/pkg/log"
)

type SignInQueryService struct {
	signInRepo commerceRepo.SignInRecordRepository
}

func NewSignInQueryService(signInRepo commerceRepo.SignInRecordRepository) *SignInQueryService {
	return &SignInQueryService{signInRepo: signInRepo}
}

func (s *SignInQueryService) GetSignInInfo(ctx context.Context, query *GetSignInInfoQuery) (*GetSignInInfoResult, error) {
	todayRecord, err := s.signInRepo.GetTodaySignIn(ctx, query.UserID)
	if err != nil {
		log.Error.Printf("SignInQueryService.GetSignInInfo: failed: %v", err)
		return nil, err
	}

	lastRecord, err := s.signInRepo.GetLastSignIn(ctx, query.UserID)
	if err != nil {
		log.Error.Printf("SignInQueryService.GetSignInInfo: get last sign-in failed: %v", err)
		return nil, err
	}

	signedInToday := todayRecord != nil
	consecutiveDays := 0
	if lastRecord != nil {
		consecutiveDays = lastRecord.ConsecutiveDays
	}

	bonus := int64(0)
	if signedInToday {
		bonus = todayRecord.BonusPoints
	}

	return &GetSignInInfoResult{
		Info: &SignInInfoDTO{
			UserID:          query.UserID,
			ConsecutiveDays: consecutiveDays,
			LastSignInAt:    0, // TODO: convert from lastRecord
			TodayBonus:      bonus,
			SignedInToday:   signedInToday,
		},
	}, nil
}
