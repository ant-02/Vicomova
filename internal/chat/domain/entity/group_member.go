package entity

import (
	"time"

	"gorm.io/gorm"
)

const (
	GroupRoleMember = 0
	GroupRoleAdmin  = 1
)

type GroupMember struct {
	ID        int64          `json:"id"`
	GroupID   int64          `json:"group_id"`
	UserID    int64          `json:"user_id"`
	Role      int8           `json:"role"` // 0=member, 1=admin
	JoinedAt  time.Time      `json:"joined_at"`
	DeletedAt gorm.DeletedAt `json:"-"`
}
