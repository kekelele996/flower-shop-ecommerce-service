package repository

import (
	"errors"

	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/util"
	"gorm.io/gorm"
)

// ReviewRepository 评价仓储。
type ReviewRepository struct {
	db *gorm.DB
}

func NewReviewRepository(db *gorm.DB) *ReviewRepository {
	return &ReviewRepository{db: db}
}

func (r *ReviewRepository) Create(tx *gorm.DB, review *model.Review) error {
	if err := tx.Create(review).Error; err != nil {
		return util.WrapAppError(50001, "create review failed", err)
	}
	return nil
}

func (r *ReviewRepository) FindByID(id uint) (*model.Review, error) {
	var review model.Review
	err := r.db.First(&review, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, util.WrapAppError(50001, "find review by id failed", err)
	}
	return &review, nil
}

func (r *ReviewRepository) FindByOrderItemID(orderItemID uint) (*model.Review, error) {
	var review model.Review
	err := r.db.Where("order_item_id = ?", orderItemID).First(&review).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, util.WrapAppError(50001, "find review by order item failed", err)
	}
	return &review, nil
}

func (r *ReviewRepository) UpdateReply(id uint, reply string) error {
	return r.db.Model(&model.Review{}).Where("id = ?", id).
		Updates(map[string]interface{}{"reply": reply, "replied_at": gorm.Expr("NOW()")}).Error
}

func (r *ReviewRepository) List(q dto.ReviewQuery) ([]model.Review, int64, error) {
	query := r.db.Model(&model.Review{}).Preload("User").Preload("Product")
	if q.ProductID > 0 {
		query = query.Where("product_id = ?", q.ProductID)
	}
	if q.OrderID > 0 {
		query = query.Where("order_id = ?", q.OrderID)
	}
	if q.UserID > 0 {
		query = query.Where("user_id = ?", q.UserID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, util.WrapAppError(50001, "count reviews failed", err)
	}

	page, pageSize := q.Page, q.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var list []model.Review
	if err := query.Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, util.WrapAppError(50001, "list reviews failed", err)
	}
	return list, total, nil
}
