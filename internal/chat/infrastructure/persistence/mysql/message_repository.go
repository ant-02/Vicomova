package mysql

import (
	"context"
	"errors"
	"time"

	"vicomova/internal/chat/domain/entity"
	repo "vicomova/internal/chat/domain/repository"
	sharedMysql "vicomova/pkg/infrastructure/mysql"
	"vicomova/pkg/log"

	"gorm.io/gorm"
)

type MessageRepository struct {
	mysql *sharedMysql.Client
}

func NewMessageRepository(mysqlClient *sharedMysql.Client) repo.MessageRepository {
	return &MessageRepository{mysql: mysqlClient}
}

func (r *MessageRepository) Create(ctx context.Context, msg *entity.Message) error {
	po := &MessagePO{
		ConversationID: msg.ConversationID,
		SenderID:       msg.SenderID,
		Content:        msg.Content,
	}
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("MessageRepository.Create: insert failed: %v", err)
		return err
	}
	msg.ID = po.ID
	return nil
}

func (r *MessageRepository) GetByID(ctx context.Context, id int64) (*entity.Message, error) {
	var po MessagePO
	err := r.mysql.WithContext(ctx).
		Where("id = ?", id).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error.Printf("MessageRepository.GetByID: failed: %v", err)
		return nil, err
	}
	return POToMessage(&po), nil
}

func (r *MessageRepository) ListByConversation(ctx context.Context, conversationID int64, cursor int64, limit int) ([]*entity.Message, bool, error) {
	var pos []MessagePO

	db := r.mysql.WithContext(ctx).Model(&MessagePO{}).
		Where("conversation_id = ?", conversationID)

	if cursor > 0 {
		db = db.Where("created_at < ?", time.UnixMilli(cursor))
	}

	if err := db.Order("created_at DESC").Limit(limit + 1).Find(&pos).Error; err != nil {
		return nil, false, err
	}

	hasMore := len(pos) > limit
	if hasMore {
		pos = pos[:limit]
	}

	messages := make([]*entity.Message, len(pos))
	for i := range pos {
		messages[i] = POToMessage(&pos[i])
	}
	return messages, hasMore, nil
}
