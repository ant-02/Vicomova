package mysql

import (
	"context"
	"errors"

	"vicomova/internal/chat/domain/entity"
	repo "vicomova/internal/chat/domain/repository"
	sharedMysql "vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/log"

	"gorm.io/gorm"
)

type GroupMemberRepository struct {
	mysql *sharedMysql.Client
}

func NewGroupMemberRepository(mysqlClient *sharedMysql.Client) repo.GroupMemberRepository {
	return &GroupMemberRepository{mysql: mysqlClient}
}

func (r *GroupMemberRepository) Create(ctx context.Context, member *entity.GroupMember) error {
	var existing GroupMemberPO
	result := r.mysql.WithContext(ctx).Unscoped().
		Where("group_id = ? AND user_id = ?", member.GroupID, member.UserID).
		First(&existing)

	if result.Error == nil {
		if existing.DeletedAt.Valid {
			if err := r.mysql.WithContext(ctx).Unscoped().Exec("UPDATE group_members SET deleted_at = NULL WHERE id = ?", existing.ID).Error; err != nil {
				log.Error.Printf("GroupMemberRepository.Create: restore failed: %v", err)
				return err
			}
		}
		member.ID = existing.ID
		return nil
	}

	if !errors.Is(result.Error, gorm.ErrRecordNotFound) {
		log.Error.Printf("GroupMemberRepository.Create: query failed: %v", result.Error)
		return result.Error
	}

	po := &GroupMemberPO{
		GroupID: member.GroupID,
		UserID:  member.UserID,
		Role:    member.Role,
	}
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("GroupMemberRepository.Create: insert failed: %v", err)
		return err
	}
	member.ID = po.ID
	return nil
}

func (r *GroupMemberRepository) Delete(ctx context.Context, groupID, userID int64) error {
	if err := r.mysql.WithContext(ctx).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Delete(&GroupMemberPO{}).Error; err != nil {
		log.Error.Printf("GroupMemberRepository.Delete: failed: %v", err)
		return err
	}
	return nil
}

func (r *GroupMemberRepository) GetMember(ctx context.Context, groupID, userID int64) (*entity.GroupMember, error) {
	var po GroupMemberPO
	err := r.mysql.WithContext(ctx).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error.Printf("GroupMemberRepository.GetMember: failed: %v", err)
		return nil, err
	}
	return POToGroupMember(&po), nil
}

func (r *GroupMemberRepository) IsMember(ctx context.Context, groupID, userID int64) (bool, error) {
	var count int64
	if err := r.mysql.WithContext(ctx).Model(&GroupMemberPO{}).
		Where("group_id = ? AND user_id = ?", groupID, userID).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GroupMemberRepository) IsAdmin(ctx context.Context, groupID, userID int64) (bool, error) {
	var count int64
	if err := r.mysql.WithContext(ctx).Model(&GroupMemberPO{}).
		Where("group_id = ? AND user_id = ? AND role = ?", groupID, userID, entity.GroupRoleAdmin).
		Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *GroupMemberRepository) ListByGroup(ctx context.Context, groupID int64) ([]*entity.GroupMember, error) {
	var pos []GroupMemberPO
	if err := r.mysql.WithContext(ctx).
		Where("group_id = ?", groupID).
		Find(&pos).Error; err != nil {
		log.Error.Printf("GroupMemberRepository.ListByGroup: failed: %v", err)
		return nil, err
	}

	members := make([]*entity.GroupMember, len(pos))
	for i := range pos {
		members[i] = POToGroupMember(&pos[i])
	}
	return members, nil
}

func (r *GroupMemberRepository) ListByUser(ctx context.Context, userID int64) ([]int64, error) {
	var pos []GroupMemberPO
	if err := r.mysql.WithContext(ctx).
		Where("user_id = ?", userID).
		Find(&pos).Error; err != nil {
		log.Error.Printf("GroupMemberRepository.ListByUser: failed: %v", err)
		return nil, err
	}

	groupIDs := make([]int64, len(pos))
	for i := range pos {
		groupIDs[i] = pos[i].GroupID
	}
	return groupIDs, nil
}
