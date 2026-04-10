package user

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"time"

	"vicomova/internal/domain/user"
	errorsPkg "vicomova/internal/pkg/errors"

	"github.com/google/uuid"
)

type RegisterCommand struct {
	Username string
	Password string
	Email    string
}

type LoginCommand struct {
	Username string
	Password string
}

type RefreshTokenCommand struct {
	RefreshToken string
}

type LogoutCommand struct {
	AccessToken string
}

type TokenResult struct {
	UserID       int64
	Username     string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int64
}

type UserResult struct {
	UserID   int64
	Username string
	Email    string
}

type UserCommandService struct {
	userRepo         user.UserRepository
	refreshTokenRepo user.RefreshTokenRepository
	tokenService     *user.TokenService
}

func NewUserCommandService(
	userRepo user.UserRepository,
	refreshTokenRepo user.RefreshTokenRepository,
	tokenService *user.TokenService,
) *UserCommandService {
	return &UserCommandService{
		userRepo:         userRepo,
		refreshTokenRepo: refreshTokenRepo,
		tokenService:     tokenService,
	}
}

func (s *UserCommandService) Register(ctx context.Context, cmd *RegisterCommand) (*UserResult, error) {
	existing, err := s.userRepo.GetByUsername(ctx, cmd.Username)
	if err != nil {
		return nil, errorsPkg.ErrInternalServer
	}
	if existing != nil {
		return nil, errorsPkg.ErrUserExists
	}

	hashedPassword := s.hashPassword(cmd.Password)

	u := &user.User{
		Username: cmd.Username,
		Password: hashedPassword,
		Email:    cmd.Email,
	}

	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, errorsPkg.ErrInternalServer
	}

	return &UserResult{
		UserID:   u.ID,
		Username: u.Username,
	}, nil
}

func (s *UserCommandService) Login(ctx context.Context, cmd *LoginCommand) (*TokenResult, error) {
	u, err := s.userRepo.GetByUsername(ctx, cmd.Username)
	if err != nil {
		return nil, errorsPkg.ErrInternalServer
	}
	if u == nil {
		return nil, errorsPkg.ErrUserNotFound
	}

	if !s.verifyPassword(cmd.Password, u.Password) {
		return nil, errorsPkg.ErrPasswordWrong
	}

	return s.generateTokenPair(ctx, u.ID, u.Username)
}

func (s *UserCommandService) RefreshToken(ctx context.Context, cmd *RefreshTokenCommand) (*TokenResult, error) {
	rt, err := s.refreshTokenRepo.GetByToken(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, errorsPkg.ErrInternalServer
	}
	if rt == nil {
		return nil, errorsPkg.ErrInvalidToken
	}
	if rt.Revoked {
		return nil, errorsPkg.ErrInvalidToken
	}
	if time.Now().After(rt.ExpiresAt) {
		return nil, errorsPkg.ErrTokenExpired
	}

	u, err := s.userRepo.GetByID(ctx, rt.UserID)
	if err != nil {
		return nil, errorsPkg.ErrInternalServer
	}
	if u == nil {
		return nil, errorsPkg.ErrUserNotFound
	}

	if err := s.refreshTokenRepo.Revoke(ctx, cmd.RefreshToken); err != nil {
		return nil, errorsPkg.ErrInternalServer
	}

	return s.generateTokenPair(ctx, u.ID, u.Username)
}

func (s *UserCommandService) Logout(ctx context.Context, cmd *LogoutCommand) error {
	claims, err := s.tokenService.ParseAccessToken(cmd.AccessToken)
	if err != nil {
		return errorsPkg.ErrInvalidToken
	}

	userID, ok := claims["user_id"].(float64)
	if !ok {
		return errorsPkg.ErrInvalidToken
	}

	if err := s.refreshTokenRepo.RevokeAllForUser(ctx, int64(userID)); err != nil {
		return errorsPkg.ErrInternalServer
	}

	return nil
}

func (s *UserCommandService) hashPassword(password string) string {
	hash := sha256.Sum256([]byte(password))
	return hex.EncodeToString(hash[:])
}

func (s *UserCommandService) verifyPassword(password, hash string) bool {
	return s.hashPassword(password) == hash
}

func (s *UserCommandService) generateTokenPair(ctx context.Context, userID int64, username string) (*TokenResult, error) {
	accessToken, err := s.tokenService.GenerateAccessToken(userID, username)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.New().String()

	rt := &user.RefreshToken{
		UserID:    userID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(user.RefreshTokenExpiry),
	}

	if err := s.refreshTokenRepo.Create(ctx, rt); err != nil {
		return nil, errorsPkg.ErrInternalServer
	}

	return &TokenResult{
		UserID:       userID,
		Username:     username,
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.tokenService.GetAccessTokenExpiry().Seconds()),
	}, nil
}
