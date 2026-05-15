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

type ConversationRepository struct {
	mysql *sharedMysql.Client
}

func NewConversationRepository(mysqlClient *sharedMysql.Client) repo.ConversationRepository {
	return &ConversationRepository{mysql: mysqlClient}
}

func (r *ConversationRepository) Create(ctx context.Context, conv *entity.Conversation) error {
	po := &ConversationPO{
		Type:   conv.Type,
		Name:   conv.Name,
		Avatar: conv.Avatar,
	}
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("ConversationRepository.Create: failed: %v", err)
		return err
	}
	conv.ID = po.ID
	return nil
}

func (r *ConversationRepository) Create1v1(ctx context.Context, userID1, userID2 int64) (*entity.Conversation, error) {
	po := &ConversationPO{
		Type: entity.ConversationType1v1,
	}
	if err := r.mysql.WithContext(ctx).Create(po).Error; err != nil {
		log.Error.Printf("ConversationRepository.Create1v1: failed: %v", err)
		return nil, err
	}
	return POToConversation(po), nil
}

func (r *ConversationRepository) Get1v1(ctx context.Context, userID1, userID2 int64) (*entity.Conversation, error) {
	var pos []ConversationPO
	err := r.mysql.WithContext(ctx).
		Where("type = ? AND deleted_at IS NULL", entity.ConversationType1v1).
		Find(&pos).Error
	if err != nil {
		log.Error.Printf("ConversationRepository.Get1v1: failed: %v", err)
		return nil, err
	}

	for _, po := range pos {
		if po.ID > 0 {
			return POToConversation(&po), nil
		}
	}
	return nil, nil
}

func (r *ConversationRepository) GetByID(ctx context.Context, id int64) (*entity.Conversation, error) {
	var po ConversationPO
	err := r.mysql.WithContext(ctx).
		Where("id = ?", id).
		First(&po).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		log.Error.Printf("ConversationRepository.GetByID: failed: %v", err)
		return nil, err
	}
	return POToConversation(&po), nil
}
