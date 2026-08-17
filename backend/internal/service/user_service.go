package service

import (
	"errors"
	"log/slog"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/repository"
	"github.com/flowershop/backend/internal/util"
)

// UserService 用户服务：注册、登录、资料、密码。
type UserService struct {
	repo     *repository.UserRepository
	auditSvc *AuditService
	cfg      *util.JWTConfig
	logger   *slog.Logger
}

func NewUserService(repo *repository.UserRepository, auditSvc *AuditService, cfg *util.JWTConfig, logger *slog.Logger) *UserService {
	return &UserService{repo: repo, auditSvc: auditSvc, cfg: cfg, logger: logger}
}

// Register 注册用户。
func (s *UserService) Register(req dto.RegisterRequest) (*dto.UserVO, error) {
	if _, err := s.repo.FindByUsername(req.Username); err == nil {
		msg := constants.MsgDuplicateUsername
		s.logger.Warn("register failed", "reason", "duplicate username", "username", req.Username)
		return nil, util.NewAppError(constants.CodeDuplicateUsername, msg)
	}
	hashed, err := util.HashPassword(req.Password)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, "hash password failed", err)
	}
	user := &model.User{
		Username: req.Username,
		Password: hashed,
		Nickname: req.Nickname,
		Email:    req.Email,
		Phone:    req.Phone,
		Role:     constants.RoleUser,
		Status:   constants.UserStatusActive,
	}
	if user.Nickname == "" {
		user.Nickname = req.Username
	}
	if err := s.repo.Create(user); err != nil {
		return nil, util.WrapAppError(constants.CodeDatabaseError, "register user failed", err)
	}
	s.logger.Info(constants.LogUserRegistered, "id", user.ID, "username", user.Username, "role", user.Role)
	s.auditSvc.Record(user.ID, user.Username, "CREATE", "user", uint64(user.ID), "register user "+user.Username)
	return ToUserVO(user), nil
}

// Login 登录并签发 JWT。
func (s *UserService) Login(req dto.LoginRequest) (*dto.TokenVO, error) {
	user, err := s.repo.FindByUsername(req.Username)
	if errors.Is(err, util.ErrNotFound) {
		return nil, util.NewAppError(constants.CodeInvalidCredentials, constants.MsgInvalidCredential)
	}
	if err != nil {
		return nil, err
	}
	if !util.CheckPassword(user.Password, req.Password) {
		s.logger.Warn("login failed", "reason", "bad password", "username", req.Username)
		return nil, util.NewAppError(constants.CodeInvalidCredentials, constants.MsgInvalidCredential)
	}
	if user.Status != constants.UserStatusActive {
		return nil, util.NewAppError(constants.CodeUserDisabled, "user="+user.Username+" is disabled")
	}
	token, err := util.GenerateToken(s.cfg.Secret, user.ID, user.Username, user.Role, s.cfg.AccessExpire)
	if err != nil {
		return nil, util.WrapAppError(constants.CodeInternalError, "generate token failed", err)
	}
	s.logger.Info(constants.LogUserLoggedIn, "id", user.ID, "username", user.Username, "role", user.Role)
	s.auditSvc.Record(user.ID, user.Username, "LOGIN", "user", uint64(user.ID), "user login")
	return &dto.TokenVO{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   s.cfg.AccessExpire,
		User:        *ToUserVO(user),
	}, nil
}

// GetProfile 获取当前用户信息。
func (s *UserService) GetProfile(userID uint) (*dto.UserVO, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	return ToUserVO(user), nil
}

// UpdateProfile 更新资料。
func (s *UserService) UpdateProfile(userID uint, req dto.UpdateProfileRequest) (*dto.UserVO, error) {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, err
	}
	if req.Nickname != "" {
		user.Nickname = req.Nickname
	}
	if req.Email != "" {
		user.Email = req.Email
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Avatar != "" {
		user.Avatar = req.Avatar
	}
	if err := s.repo.Update(user); err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogUserProfileUpdated, "id", user.ID, "nickname", user.Nickname)
	s.auditSvc.Record(user.ID, user.Username, "UPDATE", "user", uint64(user.ID), "update profile")
	return ToUserVO(user), nil
}

// ChangePassword 修改密码。
func (s *UserService) ChangePassword(userID uint, req dto.ChangePasswordRequest) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}
	if !util.CheckPassword(user.Password, req.OldPassword) {
		return util.NewAppError(constants.CodeInvalidCredentials, "old password is incorrect for user="+user.Username)
	}
	hashed, err := util.HashPassword(req.NewPassword)
	if err != nil {
		return util.WrapAppError(constants.CodeInternalError, "hash password failed", err)
	}
	if err := s.repo.UpdateFields(userID, map[string]interface{}{"password": hashed}); err != nil {
		return err
	}
	s.logger.Info(constants.LogPasswordChanged, "user_id", userID)
	s.auditSvc.Record(user.ID, user.Username, "UPDATE", "user", uint64(user.ID), "change password")
	return nil
}

// ToUserVO 模型转视图对象。
func ToUserVO(u *model.User) *dto.UserVO {
	return &dto.UserVO{
		ID:        u.ID,
		Username:  u.Username,
		Nickname:  u.Nickname,
		Email:     u.Email,
		Phone:     u.Phone,
		Avatar:    u.Avatar,
		Role:      u.Role,
		Status:    u.Status,
		CreatedAt: util.FormatTime(u.CreatedAt),
	}
}
