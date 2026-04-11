package mysql

import (
	userEntity "vicomova/internal/user/domain/entity"
)

func UserToPO(u *userEntity.User) *UserPO {
	if u == nil {
		return nil
	}
	return &UserPO{
		ID:        u.ID,
		Username:  u.Username,
		Password:  u.Password,
		Email:     u.Email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func POToUser(po *UserPO) *userEntity.User {
	if po == nil {
		return nil
	}
	return &userEntity.User{
		ID:        po.ID,
		Username:  po.Username,
		Password:  po.Password,
		Email:     po.Email,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}