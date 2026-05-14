package mock

import (
	"context"

	"vicomova/internal/interaction/domain/entity"
)

type MockCommentRepository struct {
	Comments        map[int64]*entity.Comment
	CreateErr       error
	DeleteErr       error
	GetErr          error
	ListByVideoErr  error
	ListByParentErr error
	IncrementErr    error
}

func NewMockCommentRepository() *MockCommentRepository {
	return &MockCommentRepository{
		Comments: make(map[int64]*entity.Comment),
	}
}

func (m *MockCommentRepository) Create(ctx context.Context, comment *entity.Comment) error {
	if m.CreateErr != nil {
		return m.CreateErr
	}
	comment.ID = int64(len(m.Comments) + 1)
	m.Comments[comment.ID] = comment
	return nil
}

func (m *MockCommentRepository) Delete(ctx context.Context, id int64) error {
	if m.DeleteErr != nil {
		return m.DeleteErr
	}
	delete(m.Comments, id)
	return nil
}

func (m *MockCommentRepository) Get(ctx context.Context, id int64) (*entity.Comment, error) {
	if m.GetErr != nil {
		return nil, m.GetErr
	}
	return m.Comments[id], nil
}

func (m *MockCommentRepository) ListByVideo(ctx context.Context, videoID int64, cursor int64, limit int) ([]*entity.Comment, bool, error) {
	if m.ListByVideoErr != nil {
		return nil, false, m.ListByVideoErr
	}
	var result []*entity.Comment
	for _, c := range m.Comments {
		if c.VideoID == videoID && c.ParentID == 0 {
			result = append(result, c)
		}
	}
	if len(result) > limit {
		return result[:limit], true, nil
	}
	return result, false, nil
}

func (m *MockCommentRepository) ListByParent(ctx context.Context, parentID int64, cursor int64, limit int) ([]*entity.Comment, bool, error) {
	if m.ListByParentErr != nil {
		return nil, false, m.ListByParentErr
	}
	var result []*entity.Comment
	for _, c := range m.Comments {
		if c.ParentID == parentID {
			result = append(result, c)
		}
	}
	if len(result) > limit {
		return result[:limit], true, nil
	}
	return result, false, nil
}

func (m *MockCommentRepository) IncrementLike(ctx context.Context, id int64) error {
	if m.IncrementErr != nil {
		return m.IncrementErr
	}
	if c, ok := m.Comments[id]; ok {
		c.LikeCount++
	}
	return nil
}
