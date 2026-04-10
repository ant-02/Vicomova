package service

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"vicomova/internal/domain/user"
	errorsPkg "vicomova/internal/pkg/errors"
	"vicomova/internal/repository"

	"github.com/golang-jwt/jwt/v5"
)

type UserService struct {
	repo      repository.UserRepository
	jwtSecret string
}

func NewUserService(repo repository.UserRepository, jwtSecret string) *UserService {
	return &UserService{
		repo:      repo,
		jwtSecret: jwtSecret,
	}
}

type RegisterInput struct {
	Username string
	Password string
	Email    string
}

type LoginInput struct {
	Username string
	Password string
}

type UserOutput struct {
	UserID   int64
	Username string
	Email    string
	Token    string
}

func (s *UserService) Register(ctx context.Context, input *RegisterInput) (*UserOutput, error) {
	// 检查用户是否已存在
	existing, err := s.repo.GetByUsername(ctx, input.Username)
	if err != nil {
		return nil, errorsPkg.ErrInternalServer
	}
	if existing != nil {
		return nil, errorsPkg.ErrUserExists
	}

	// 密码哈希
	hashedPassword := s.hashPassword(input.Password)

	user := &user.User{
		Username: input.Username,
		Password: hashedPassword,
		Email:    input.Email,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, errorsPkg.ErrInternalServer
	}

	return &UserOutput{
		UserID:   user.ID,
		Username: user.Username,
	}, nil
}

func (s *UserService) Login(ctx context.Context, input *LoginInput) (*UserOutput, error) {
	user, err := s.repo.GetByUsername(ctx, input.Username)
	if err != nil {
		return nil, errorsPkg.ErrInternalServer
	}
	if user == nil {
		return nil, errorsPkg.ErrUserNotFound
	}

	// 验证密码
	if !s.verifyPassword(input.Password, user.Password) {
		return nil, errorsPkg.ErrPasswordWrong
	}

	// 生成 JWT
	token, err := s.generateToken(user.ID, user.Username)
	if err != nil {
		return nil, errorsPkg.ErrInternalServer
	}

	return &UserOutput{
		UserID:   user.ID,
		Username: user.Username,
		Token:    token,
	}, nil
}

func (s *UserService) GetUser(ctx context.Context, userID int64) (*UserOutput, error) {
	user, err := s.repo.GetByID(ctx, userID)
	if err != nil {
		return nil, errorsPkg.ErrInternalServer
	}
	if user == nil {
		return nil, errorsPkg.ErrUserNotFound
	}

	return &UserOutput{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
	}, nil
}

func (s *UserService) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func (s *UserService) verifyPassword(password, hash string) bool {
	return s.hashPassword(password) == hash
}

func (s *UserService) generateToken(userID int64, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}
