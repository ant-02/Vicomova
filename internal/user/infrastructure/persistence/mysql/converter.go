package mysql

import (
	userEntity "vicomova/internal/user/domain/entity"
	userVO "vicomova/internal/user/domain/valueobject"
)

func UserToPO(u *userEntity.User) *UserPO {
	if u == nil {
		return nil
	}
	var username, password, email string
	if u.Username != nil {
		username = u.Username.String()
	}
	if u.Password != nil {
		password = u.Password.Hash()
	}
	if u.Email != nil {
		email = u.Email.String()
	}
	return &UserPO{
		ID:        u.ID,
		Username:  username,
		Password:  password,
		Email:     email,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}

func POToUser(po *UserPO) *userEntity.User {
	if po == nil {
		return nil
	}
	username, _ := userVO.NewUsername(po.Username)
	email, _ := userVO.NewEmail(po.Email)
	password := userVO.NewPasswordFromHash(po.Password)
	return &userEntity.User{
		ID:        po.ID,
		Username:  username,
		Password:  password,
		Email:     email,
		CreatedAt: po.CreatedAt,
		UpdatedAt: po.UpdatedAt,
	}
}