package service

import (
	"log/slog"
	"os"
	"testing"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/repository"
	"github.com/flowershop/backend/internal/util"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newTestUserService(t *testing.T) (*UserService, *repository.UserRepository) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.AutoMigrate(&model.User{}, &model.AuditLog{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	userRepo := repository.NewUserRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	auditSvc := NewAuditService(auditRepo)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	svc := NewUserService(userRepo, auditSvc, &util.JWTConfig{Secret: "test-secret", AccessExpire: 3600}, logger)
	return svc, userRepo
}

func TestUserServiceRegisterAndLogin(t *testing.T) {
	svc, userRepo := newTestUserService(t)

	registerReq := dto.RegisterRequest{Username: "alice", Password: "secret123", Nickname: "爱丽丝"}
	user, err := svc.Register(registerReq)
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}
	if user.Username != "alice" || user.Role != constants.RoleUser {
		t.Errorf("register result mismatch: %+v", user)
	}
	if _, err := svc.Register(registerReq); err == nil {
		t.Error("duplicate register should fail")
	}

	tests := []struct {
		name     string
		username string
		password string
		wantErr  bool
	}{
		{"correct credentials", "alice", "secret123", false},
		{"wrong password", "alice", "wrongpass", true},
		{"unknown user", "bob", "secret123", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			token, err := svc.Login(dto.LoginRequest{Username: tt.username, Password: tt.password})
			if tt.wantErr {
				if err == nil {
					t.Error("expected error")
				}
				return
			}
			if err != nil {
				t.Fatalf("login failed: %v", err)
			}
			if token.AccessToken == "" {
				t.Error("access token empty")
			}
		})
	}

	// 密码修改后旧密码失效
	if err := svc.ChangePassword(user.ID, dto.ChangePasswordRequest{OldPassword: "secret123", NewPassword: "newpass123"}); err != nil {
		t.Fatalf("change password failed: %v", err)
	}
	if _, err := svc.Login(dto.LoginRequest{Username: "alice", Password: "secret123"}); err == nil {
		t.Error("old password should be invalid after change")
	}
	if _, err := userRepo.FindByUsername("alice"); err != nil {
		t.Errorf("find by username failed: %v", err)
	}
}
