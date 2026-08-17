package service

import (
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/repository"
)

// AuditService 审计服务：service 层埋点统一入口。
type AuditService struct {
	repo *repository.AuditRepository
}

func NewAuditService(repo *repository.AuditRepository) *AuditService {
	return &AuditService{repo: repo}
}

// Record 异步记录审计日志（不阻塞主流程）。
func (s *AuditService) Record(userID uint, username, action, entity string, entityID uint64, detail string) {
	if s == nil || s.repo == nil {
		return
	}
	go func() {
		_ = s.repo.Create(&model.AuditLog{
			UserID:   userID,
			Username: username,
			Action:   action,
			Entity:   entity,
			EntityID: uintToStr(entityID),
			Detail:   detail,
		})
	}()
}

func uintToStr(v uint64) string {
	if v == 0 {
		return ""
	}
	buf := make([]byte, 0, 20)
	for v > 0 {
		buf = append([]byte{byte('0' + v%10)}, buf...)
		v /= 10
	}
	return string(buf)
}
