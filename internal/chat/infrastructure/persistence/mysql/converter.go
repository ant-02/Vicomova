package mysql

import (
	"vicomova/internal/chat/domain/entity"
)

func FollowToPO(f *entity.Follow) *FollowPO {
	return &FollowPO{
		ID:          f.ID,
		FollowerID:  f.FollowerID,
		FollowingID: f.FollowingID,
	}
}

func POToFollow(po *FollowPO) *entity.Follow {
	return &entity.Follow{
		ID:          po.ID,
		FollowerID:  po.FollowerID,
		FollowingID: po.FollowingID,
		CreatedAt:   po.CreatedAt,
	}
}

func ConversationToPO(c *entity.Conversation) *ConversationPO {
	return &ConversationPO{
		ID:        c.ID,
		Type:      c.Type,
		Name:      c.Name,
		Avatar:    c.Avatar,
		CreatedAt: c.CreatedAt,
		UpdatedAt: c.UpdatedAt,
	}
}

func POToConversation(po *ConversationPO) *entity.Conversation {
	return &entity.Conversation{
		ID:        po.ID,
		Type:      po.Type,
		Name:      po.Name,
		Avatar:    po.Avatar,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}

func GroupMemberToPO(m *entity.GroupMember) *GroupMemberPO {
	return &GroupMemberPO{
		ID:       m.ID,
		GroupID:  m.GroupID,
		UserID:   m.UserID,
		Role:     m.Role,
		JoinedAt: m.JoinedAt,
	}
}

func POToGroupMember(po *GroupMemberPO) *entity.GroupMember {
	return &entity.GroupMember{
		ID:       po.ID,
		GroupID:  po.GroupID,
		UserID:   po.UserID,
		Role:     po.Role,
		JoinedAt: po.JoinedAt,
	}
}

func MessageToPO(m *entity.Message) *MessagePO {
	return &MessagePO{
		ID:             m.ID,
		ConversationID: m.ConversationID,
		SenderID:       m.SenderID,
		Content:        m.Content,
		CreatedAt:      m.CreatedAt,
	}
}

func POToMessage(po *MessagePO) *entity.Message {
	return &entity.Message{
		ID:             po.ID,
		ConversationID: po.ConversationID,
		SenderID:       po.SenderID,
		Content:        po.Content,
		CreatedAt:      po.CreatedAt,
	}
}
