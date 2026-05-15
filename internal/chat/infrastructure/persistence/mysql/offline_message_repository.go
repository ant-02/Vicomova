package mysql

import (
	"context"

	"vicomova/internal/chat/domain/entity"
	repo "vicomova/internal/chat/domain/repository"
	sharedMysql "vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/log"
)

type OfflineMessageRepository struct {
	mysql *sharedMysql.Client
}

func NewOfflineMessageRepository(mysqlClient *sharedMysql.Client) repo.OfflineMessageRepository {
	return &OfflineMessageRepository{mysql: mysqlClient}
}

func (r *OfflineMessageRepository) Create(ctx context.Context, msg *entity.OfflineMessage) error {
	po := &OfflineMessagePO{
		ToUserID:   msg.ToUserID,
		FromUserID: msg.FromUserID,
		Content:    msg.Content,
	}
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("OfflineMessageRepository.Create: insert failed: %v", err)
		return err
	}
	msg.ID = po.ID
	return nil
}

func (r *OfflineMessageRepository) ListByUser(ctx context.Context, userID int64, limit int) ([]*entity.OfflineMessage, error) {
	var pos []OfflineMessagePO
	if err := r.mysql.WithContext(ctx).Model(&OfflineMessagePO{}).
		Where("to_user_id = ? AND is_read = ?", userID, false).
		Order("created_at DESC").
		Limit(limit).
		Find(&pos).Error; err != nil {
		log.Error.Printf("OfflineMessageRepository.ListByUser: failed: %v", err)
		return nil, err
	}

	messages := make([]*entity.OfflineMessage, len(pos))
	for i := range pos {
		messages[i] = POToOfflineMessage(&pos[i])
	}
	return messages, nil
}

func (r *OfflineMessageRepository) MarkAsRead(ctx context.Context, userID int64) error {
	if err := r.mysql.WithContext(ctx).Model(&OfflineMessagePO{}).
		Where("to_user_id = ? AND is_read = ?", userID, false).
		Update("is_read", true).Error; err != nil {
		log.Error.Printf("OfflineMessageRepository.MarkAsRead: failed: %v", err)
		return err
	}
	return nil
}

func (r *OfflineMessageRepository) UnreadCount(ctx context.Context, userID int64) (int64, error) {
	var count int64
	if err := r.mysql.WithContext(ctx).Model(&OfflineMessagePO{}).
		Where("to_user_id = ? AND is_read = ?", userID, false).
		Count(&count).Error; err != nil {
		return 0, err
	}
	return count, nil
}

func POToOfflineMessage(po *OfflineMessagePO) *entity.OfflineMessage {
	return &entity.OfflineMessage{
		ID:         po.ID,
		ToUserID:   po.ToUserID,
		FromUserID: po.FromUserID,
		Content:    po.Content,
		CreatedAt:  po.CreatedAt,
		IsRead:     po.IsRead,
	}
}

func GetOfflineMessageRepository(mysqlClient *sharedMysql.Client) repo.OfflineMessageRepository {
	return &OfflineMessageRepository{mysql: mysqlClient}
}
