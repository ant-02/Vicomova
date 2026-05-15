package query

import (
	"context"

	commerceRepo "vicomova/internal/commerce/domain/repository"
)

type MembershipQueryService struct {
	membershipRepo commerceRepo.MembershipRepository
}

func NewMembershipQueryService(membershipRepo commerceRepo.MembershipRepository) *MembershipQueryService {
	return &MembershipQueryService{membershipRepo: membershipRepo}
}

func (s *MembershipQueryService) GetMembership(ctx context.Context, query *GetMembershipQuery) (*GetMembershipResult, error) {
	membership, err := s.membershipRepo.GetActiveMembership(ctx, query.UserID)
	if err != nil {
		return nil, err
	}
	if membership == nil {
		return &GetMembershipResult{Membership: nil}, nil
	}
	return &GetMembershipResult{
		Membership: &MembershipDTO{
			ID:       membership.ID,
			UserID:   membership.UserID,
			Level:    int(membership.Level),
			StartAt:  membership.StartAt.Unix(),
			ExpireAt: membership.ExpireAt.Unix(),
			IsActive: membership.IsActive,
		},
	}, nil
}
