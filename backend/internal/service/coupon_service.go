package service

import (
	"log/slog"
	"time"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/repository"
	"github.com/flowershop/backend/internal/util"
	"gorm.io/gorm"
)

// CouponService 优惠券服务：模板创建、领取、查询。
type CouponService struct {
	db       *gorm.DB
	repo     *repository.CouponRepository
	auditSvc *AuditService
	logger   *slog.Logger
}

func NewCouponService(db *gorm.DB, repo *repository.CouponRepository, auditSvc *AuditService, logger *slog.Logger) *CouponService {
	return &CouponService{db: db, repo: repo, auditSvc: auditSvc, logger: logger}
}

// CreateTemplate 管理员创建优惠券模板并批量发放。
func (s *CouponService) CreateTemplate(req dto.CouponTemplateCreateRequest) (int, error) {
	now := time.Now()
	validFrom := now
	validTo := now.AddDate(0, 0, req.ValidDays)
	created := 0
	err := s.db.Transaction(func(tx *gorm.DB) error {
		for i := 0; i < req.Total; i++ {
			c, err := s.repo.Claim(tx, 0, 0, req.Name, req.Type, req.Threshold, req.Amount, req.DiscountRate, validFrom, validTo)
			if err != nil {
				return err
			}
			_ = c
			created++
		}
		return nil
	})
	if err != nil {
		return 0, err
	}
	s.logger.Info(constants.LogCouponCreated, "id", 0, "name", req.Name, "type", req.Type)
	return created, nil
}

// Claim 用户领取优惠券（模拟从模板领取，模板ID 0 表示公共券池）。
func (s *CouponService) Claim(userID uint, templateID uint) (*dto.CouponVO, error) {
	now := time.Now()
	var c *model.Coupon
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		c, err = s.repo.Claim(tx, templateID, userID, "新人满减券", constants.CouponTypeFullReduction, 100, 20, 0, now, now.AddDate(0, 0, 30))
		return err
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogCouponClaimed, "coupon_id", c.ID, "user_id", userID, "status", c.Status)
	return s.ToVO(c), nil
}

// ListByUser 我的优惠券。
func (s *CouponService) ListByUser(userID uint) ([]dto.CouponVO, error) {
	list, err := s.repo.ListByUser(userID)
	if err != nil {
		return nil, err
	}
	vos := make([]dto.CouponVO, 0, len(list))
	for i := range list {
		vos = append(vos, *s.ToVO(&list[i]))
	}
	return vos, nil
}

// ListAvailable 可用优惠券（复用 repository.ListAvailable）。
func (s *CouponService) ListAvailable(userID uint) ([]dto.CouponVO, error) {
	list, err := s.repo.ListAvailable(userID)
	if err != nil {
		return nil, err
	}
	vos := make([]dto.CouponVO, 0, len(list))
	for i := range list {
		vos = append(vos, *s.ToVO(&list[i]))
	}
	return vos, nil
}

func (s *CouponService) ToVO(c *model.Coupon) *dto.CouponVO {
	return &dto.CouponVO{
		ID:           c.ID,
		Name:         c.Name,
		Type:         c.Type,
		TypeText:     util.CouponTypeText(c.Type),
		Threshold:    c.Threshold,
		Amount:       c.Amount,
		DiscountRate: c.DiscountRate,
		Status:       c.Status,
		ValidFrom:    util.FormatTime(c.ValidFrom),
		ValidTo:      util.FormatTime(c.ValidTo),
	}
}
