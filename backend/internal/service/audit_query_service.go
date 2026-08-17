package service

import (
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/repository"
	"github.com/flowershop/backend/internal/util"
)

// AuditQueryService 审计日志查询服务（与 AuditService 分离职责）。
type AuditQueryService struct {
	repo *repository.AuditRepository
}

func NewAuditQueryService(repo *repository.AuditRepository) *AuditQueryService {
	return &AuditQueryService{repo: repo}
}

func (s *AuditQueryService) List(q dto.AuditLogQuery) ([]dto.AuditLogVO, int64, error) {
	list, total, err := s.repo.List(q)
	if err != nil {
		return nil, 0, err
	}
	vos := make([]dto.AuditLogVO, 0, len(list))
	for _, l := range list {
		vos = append(vos, toAuditVO(&l))
	}
	return vos, total, nil
}

func toAuditVO(l *model.AuditLog) dto.AuditLogVO {
	return dto.AuditLogVO{
		ID:        l.ID,
		UserID:    l.UserID,
		Username:  l.Username,
		Action:    l.Action,
		Method:    l.Method,
		Path:      l.Path,
		Entity:    l.Entity,
		EntityID:  l.EntityID,
		Detail:    l.Detail,
		RequestID: l.RequestID,
		IP:        l.IP,
		CreatedAt: util.FormatTime(l.CreatedAt),
	}
}
