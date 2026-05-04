package command

import (
	"context"

	userEntity "vicomova/internal/user/domain/entity"
	userVO "vicomova/internal/user/domain/valueobject"
	constants "vicomova/internal/shared/pkg/constants"
	errorsPkg "vicomova/internal/shared/pkg/errors"
	"vicomova/internal/shared/pkg/log"
	"vicomova/internal/shared/pkg/utils"

	"github.com/google/uuid"
)

func (s *UserCommandService) SendVerificationCode(ctx context.Context, cmd *SendVerificationCodeCommand) (*SendCodeResult, error) {
	// 验证值对象格式
	username, err := userVO.NewUsername(cmd.Username)
	if err != nil {
		log.Warn.Printf("SendVerificationCode: invalid username: %v", err)
		return &SendCodeResult{
			Success: false,
			Message: "invalid username",
		}, nil
	}

	_, err = userVO.NewEmail(cmd.Email)
	if err != nil {
		log.Warn.Printf("SendVerificationCode: invalid email: %v", err)
		return &SendCodeResult{
			Success: false,
			Message: "invalid email",
		}, nil
	}

	// 检查用户是否已存在
	existing, err := s.userRepo.GetByUsername(ctx, username)
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
	code := utils.GenerateEmailCode(constants.EmailCodeLength)

	// 存储验证码到 Redis
	if err := s.emailCodeRepo.Store(ctx, cmd.Email, code, constants.EmailCodeTTL); err != nil {
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
	if err := s.emailService.SendVerificationCode(ctx, cmd.Email, code); err != nil {
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

	// 创建值对象
	username, err := userVO.NewUsername(cmd.Username)
	if err != nil {
		return nil, errorsPkg.ErrInvalidUsername
	}

	// 使用 hasher 生成密码哈希
	password, err := userVO.NewPasswordFromPlain(cmd.Password, s.passwordHasher)
	if err != nil {
		return nil, errorsPkg.ErrInvalidPassword
	}

	email, err := userVO.NewEmail(cmd.Email)
	if err != nil {
		return nil, errorsPkg.ErrInvalidEmail
	}

	// 检查用户是否已存在
	existing, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errorsPkg.ErrInternalServer
	}
	if existing != nil {
		return nil, errorsPkg.ErrUserExists
	}

	// 创建用户
	u := userEntity.NewUser(username, password, email)
	if err := s.userRepo.Create(ctx, u); err != nil {
		return nil, errorsPkg.ErrInternalServer
	}

	// 删除已使用的验证码
	_ = s.emailCodeRepo.Delete(ctx, cmd.Email)

	return &UserResult{
		UserID:   u.ID,
		Username: u.Username.String(),
	}, nil
}

func (s *UserCommandService) Login(ctx context.Context, cmd *LoginCommand) (*TokenResult, error) {
	username, err := userVO.NewUsername(cmd.Username)
	if err != nil {
		return nil, errorsPkg.ErrUserNotFound
	}

	u, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		return nil, errorsPkg.ErrInternalServer
	}
	if u == nil {
		return nil, errorsPkg.ErrUserNotFound
	}

	// 使用 hasher 验证密码
	if !s.passwordHasher.Verify(cmd.Password, u.Password.Hash()) {
		return nil, errorsPkg.ErrPasswordWrong
	}

	return s.generateTokenPair(ctx, u.ID, u.Username)
}

func (s *UserCommandService) RefreshToken(ctx context.Context, cmd *RefreshTokenCommand) (*TokenResult, error) {
	exists, err := s.refreshTokenRepo.Exists(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, errorsPkg.ErrInternalServer
	}
	if !exists {
		return nil, errorsPkg.ErrInvalidToken
	}

	userID, err := s.refreshTokenRepo.GetUserID(ctx, cmd.RefreshToken)
	if err != nil {
		return nil, errorsPkg.ErrInternalServer
	}

	u, err := s.userRepo.GetByID(ctx, userID)
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
	if err := s.refreshTokenRepo.Revoke(ctx, cmd.RefreshToken); err != nil {
		return errorsPkg.ErrInternalServer
	}
	return nil
}

func (s *UserCommandService) generateTokenPair(ctx context.Context, userID int64, username *userVO.Username) (*TokenResult, error) {
	accessToken, err := s.tokenService.GenerateAccessToken(userID, username.String())
	if err != nil {
		return nil, err
	}

	refreshToken := uuid.New().String()
	ttl := s.tokenService.GetRefreshTokenExpiry()

	if err := s.refreshTokenRepo.Create(ctx, userID, refreshToken, ttl); err != nil {
		return nil, errorsPkg.ErrInternalServer
	}

	return &TokenResult{
		UserID:       userID,
		Username:     username.String(),
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
		ExpiresIn:    int64(s.tokenService.GetAccessTokenExpiry().Seconds()),
	}, nil
}
