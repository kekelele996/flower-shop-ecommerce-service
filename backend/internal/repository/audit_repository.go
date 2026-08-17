package repository

import (
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/util"
	"gorm.io/gorm"
)

// AuditRepository 审计日志仓储。
type AuditRepository struct {
	db *gorm.DB
}

func NewAuditRepository(db *gorm.DB) *AuditRepository {
	return &AuditRepository{db: db}
}

func (r *AuditRepository) Create(log *model.AuditLog) error {
	if err := r.db.Create(log).Error; err != nil {
		return util.WrapAppError(50001, "create audit log failed", err)
	}
	return nil
}

func (r *AuditRepository) List(q dto.AuditLogQuery) ([]model.AuditLog, int64, error) {
	query := r.db.Model(&model.AuditLog{})
	if q.Username != "" {
		query = query.Where("username LIKE ?", "%"+q.Username+"%")
	}
	if q.Action != "" {
		query = query.Where("action = ?", q.Action)
	}
	if q.Entity != "" {
		query = query.Where("entity = ?", q.Entity)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, util.WrapAppError(50001, "count audit logs failed", err)
	}

	page, pageSize := q.Page, q.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var list []model.AuditLog
	if err := query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, util.WrapAppError(50001, "list audit logs failed", err)
	}
	return list, total, nil
}
