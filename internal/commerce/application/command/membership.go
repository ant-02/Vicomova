package command

import (
	"context"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
	"vicomova/pkg/log"
)

type MembershipCommandService struct {
	membershipRepo commerceRepo.MembershipRepository
}

func NewMembershipCommandService(membershipRepo commerceRepo.MembershipRepository) *MembershipCommandService {
	return &MembershipCommandService{
		membershipRepo: membershipRepo,
	}
}

func (s *MembershipCommandService) BuyMembership(ctx context.Context, cmd *BuyMembershipCommand) (*BuyMembershipResult, error) {
	membership := commerceEntity.NewMembership(cmd.UserID, commerceEntity.MembershipLevel(cmd.Level), 30) // 默认30天
	if err := s.membershipRepo.Create(ctx, membership); err != nil {
		log.Error.Printf("MembershipCommandService.BuyMembership: failed: %v", err)
		return nil, err
	}
	return &BuyMembershipResult{
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
