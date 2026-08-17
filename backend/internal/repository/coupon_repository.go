package repository

import (
	"errors"
	"time"

	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// CouponRepository 优惠券仓储。
type CouponRepository struct {
	db *gorm.DB
}

func NewCouponRepository(db *gorm.DB) *CouponRepository {
	return &CouponRepository{db: db}
}

func (r *CouponRepository) Create(tx *gorm.DB, c *model.Coupon) error {
	if err := tx.Create(c).Error; err != nil {
		return util.WrapAppError(50001, "create coupon failed", err)
	}
	return nil
}

func (r *CouponRepository) FindByID(id uint) (*model.Coupon, error) {
	var c model.Coupon
	err := r.db.First(&c, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, util.WrapAppError(50001, "find coupon by id failed", err)
	}
	return &c, nil
}

func (r *CouponRepository) FindByIDForUpdate(tx *gorm.DB, id uint) (*model.Coupon, error) {
	var c model.Coupon
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&c, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, util.WrapAppError(50001, "lock coupon failed", err)
	}
	return &c, nil
}

func (r *CouponRepository) ListByUser(userID uint) ([]model.Coupon, error) {
	var list []model.Coupon
	if err := r.db.Where("user_id = ?", userID).Order("id desc").Find(&list).Error; err != nil {
		return nil, util.WrapAppError(50001, "list user coupons failed", err)
	}
	return list, nil
}

func (r *CouponRepository) ListAvailable(userID uint) ([]model.Coupon, error) {
	now := time.Now()
	var list []model.Coupon
	if err := r.db.Where("user_id = ? AND status = ? AND valid_from <= ? AND valid_to >= ?",
		userID, "UNUSED", now, now).Order("id desc").Find(&list).Error; err != nil {
		return nil, util.WrapAppError(50001, "list available coupons failed", err)
	}
	return list, nil
}

// Claim 领取优惠券模板：生成用户券。
func (r *CouponRepository) Claim(tx *gorm.DB, templateID, userID uint, name, ctype string, threshold, amount, rate float64, validFrom, validTo time.Time) (*model.Coupon, error) {
	c := &model.Coupon{
		UserID:       userID,
		TemplateID:   templateID,
		Name:         name,
		Type:         ctype,
		Threshold:    threshold,
		Amount:       amount,
		DiscountRate: rate,
		Status:       "UNUSED",
		ValidFrom:    validFrom,
		ValidTo:      validTo,
	}
	if err := tx.Create(c).Error; err != nil {
		return nil, util.WrapAppError(50001, "claim coupon failed", err)
	}
	return c, nil
}

func (r *CouponRepository) MarkUsed(tx *gorm.DB, id uint, orderID uint) error {
	now := time.Now()
	if err := tx.Model(&model.Coupon{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": "USED", "used_at": &now, "order_id": orderID}).Error; err != nil {
		return util.WrapAppError(50001, "mark coupon used failed", err)
	}
	return nil
}

func (r *CouponRepository) MarkUnused(tx *gorm.DB, id uint) error {
	if err := tx.Model(&model.Coupon{}).Where("id = ?", id).
		Updates(map[string]interface{}{"status": "UNUSED", "used_at": nil, "order_id": 0}).Error; err != nil {
		return util.WrapAppError(50001, "mark coupon unused failed", err)
	}
	return nil
}
