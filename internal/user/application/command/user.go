package command

import (
	"context"

	userEntity "vicomova/internal/user/domain/entity"
	userVO "vicomova/internal/user/domain/valueobject"
	constants "vicomova/internal/shared/pkg/constants"
	errorsPkg "vicomova/internal/shared/pkg/errors"
	"vicomova/internal/shared/pkg/log"

	"github.com/google/uuid"
)

func (s *UserCommandService) SendVerificationCode(ctx context.Context, cmd *SendVerificationCodeCommand) (*SendCodeResult, error) {
	// 检查用户是否已存在
	existing, err := s.userRepo.GetByUsername(ctx, cmd.Username)
	if err != nil {
		log.Error.Printf("SendVerificationCode: GetByUsername failed: %v", err)
		return nil, errorsPkg.ErrInternalServer
	}
	if existing != nil {
		log.Warn.Printf("SendVerificationCode: username %s already exists", cmd.Username)
		return &SendCodeResult{
			Success: false,
			Message: "username already exists",
		}, nil
	}

	// 生成验证码
	code := userVO.GenerateEmailCode(cmd.Email)

	// 存储验证码到 Redis
	if err := s.emailCodeRepo.Store(ctx, cmd.Email, code.Code, constants.EmailCodeTTL); err != nil {
		log.Error.Printf("SendVerificationCode: Redis Store failed for %s: %v", cmd.Email, err)
		return nil, errorsPkg.ErrInternalServer
	}

	// 发送邮件
	if s.emailService == nil {
		log.Error.Printf("SendVerificationCode: emailService is nil")
		return &SendCodeResult{
			Success: false,
			Message: "email service not configured",
		}, errorsPkg.ErrInternalServer
	}
	if err := s.emailService.SendVerificationCode(ctx, cmd.Email, code.Code); err != nil {
		log.Error.Printf("SendVerificationCode: SendEmail failed for %s: %v", cmd.Email, err)
		return &SendCodeResult{
			Success: false,
			Message: "failed to send email",
		}, err
	}

	log.Info.Printf("SendVerificationCode: code sent to %s", cmd.Email)

	return &SendCodeResult{
		Success: true,
 		Message: "verification code sent",
	}, nil
}

func (s *UserCommandService) VerifyAndRegister(ctx context.Context, cmd *VerifyAndRegisterCommand) (*UserResult, error) {
	// 验证验证码
	valid, err := s.emailCodeRepo.Verify(ctx, cmd.Email, cmd.Code)
	if err != nil {
		return nil, errorsPkg.ErrInternalServer
	}
	if !valid {
		return nil, errorsPkg.ErrInvalidToken
	}

	// 检查用户是否已存在
	existing, err := s.userRepo.GetByUsername(ctx, cmd.Username)
	if err != nil {
		return nil, errorsPkg.ErrInternalServer
	}
	if existing != nil {
		return nil, errorsPkg.ErrUserExists
	}

	// 创建用户
	u := userEntity.NewUser(cmd.Username, cmd.Password, cmd.Email)
	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, errorsPkg.ErrInternalServer
	}

	// 删除已使用的验证码
	_ = s.emailCodeRepo.Delete(ctx, cmd.Email)

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

	if !u.CanLogin(cmd.Password) {
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
	if rt.IsRevoked() {
		return nil, errorsPkg.ErrInvalidToken
	}
	if rt.IsExpired() {
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

func (s *UserCommandService) generateTokenPair(ctx context.Context, userID int64, username string) (*TokenResult, error) {
	accessToken, err := s.tokenService.GenerateAccessToken(userID, username)
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.New().String()

	rt := userVO.NewRefreshToken(userID, refreshToken, s.tokenService.GetRefreshTokenExpiry())

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
