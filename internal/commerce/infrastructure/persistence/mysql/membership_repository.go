package mysql

import (
	"context"
	"time"

	commerceEntity "vicomova/internal/commerce/domain/entity"
	commerceRepo "vicomova/internal/commerce/domain/repository"
	sharedMysql "vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/log"

	"gorm.io/gorm"
)

type MembershipRepository struct {
	mysql *sharedMysql.Client
}

func NewMembershipRepository(mysqlClient *sharedMysql.Client) commerceRepo.MembershipRepository {
	return &MembershipRepository{mysql: mysqlClient}
}

func (r *MembershipRepository) Create(ctx context.Context, m *commerceEntity.Membership) error {
	po := MembershipToPO(m)
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("MembershipRepository.Create: failed: %v", err)
		return err
	}
	m.ID = po.ID
	return nil
}

func (r *MembershipRepository) GetByID(ctx context.Context, id int64) (*commerceEntity.Membership, error) {
	var po MembershipPO
	err := r.mysql.WithContext(ctx).First(&po, id).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("MembershipRepository.GetByID: failed: %v", err)
		return nil, err
	}
	return POToMembership(&po), nil
}

func (r *MembershipRepository) GetActiveMembership(ctx context.Context, userID int64) (*commerceEntity.Membership, error) {
	var po MembershipPO
	err := r.mysql.WithContext(ctx).Where("user_id = ? AND is_active = ? AND expire_at > ?",
		userID, true, time.Now()).Order("level DESC").First(&po).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		log.Error.Printf("MembershipRepository.GetActiveMembership: failed: %v", err)
		return nil, err
	}
	return POToMembership(&po), nil
}

func (r *MembershipRepository) Update(ctx context.Context, m *commerceEntity.Membership) error {
	po := MembershipToPO(m)
	if err := r.mysql.WithContext(ctx).Save(po).Error; err != nil {
		log.Error.Printf("MembershipRepository.Update: failed: %v", err)
		return err
	}
	return nil
}

func (r *MembershipRepository) ListByUser(ctx context.Context, userID int64) ([]*commerceEntity.Membership, error) {
	var pos []MembershipPO
	if err := r.mysql.WithContext(ctx).Where("user_id = ?", userID).Order("created_at DESC").Find(&pos).Error; err != nil {
		log.Error.Printf("MembershipRepository.ListByUser: failed: %v", err)
		return nil, err
	}
	memberships := make([]*commerceEntity.Membership, len(pos))
	for i := range pos {
		memberships[i] = POToMembership(&pos[i])
	}
	return memberships, nil
}
