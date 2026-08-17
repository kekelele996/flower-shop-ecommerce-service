package repository

import (
	"errors"
	"fmt"

	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// ProductRepository 商品仓储。
type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(p *model.Product) error {
	if err := r.db.Create(p).Error; err != nil {
		return util.WrapAppError(50001, "create product failed", err)
	}
	return nil
}

func (r *ProductRepository) Update(p *model.Product) error {
	if err := r.db.Save(p).Error; err != nil {
		return util.WrapAppError(50001, "update product failed", err)
	}
	return nil
}

func (r *ProductRepository) UpdateFields(id uint, fields map[string]interface{}) error {
	if err := r.db.Model(&model.Product{}).Where("id = ?", id).Updates(fields).Error; err != nil {
		return util.WrapAppError(50001, "update product fields failed", err)
	}
	return nil
}

func (r *ProductRepository) FindByID(id uint) (*model.Product, error) {
	var p model.Product
	err := r.db.First(&p, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, util.WrapAppError(50001, "find product by id failed", err)
	}
	return &p, nil
}

// Search 商品搜索：关键词、分类、价格区间、包邮、发货地、排序、分页。
// 复用：商品列表与首页推荐共用本方法。
func (r *ProductRepository) Search(q dto.ProductQuery, onlyOnSale bool) ([]model.Product, int64, error) {
	query := r.db.Model(&model.Product{})
	if onlyOnSale {
		query = query.Where("status = ?", "OFF_SALE")
	} else if q.Status != "" {
		query = query.Where("status = ?", q.Status)
	}
	if q.Keyword != "" {
		like := "%" + q.Keyword + "%"
		query = query.Where("name LIKE ? OR sub_title LIKE ? OR description LIKE ?", like, like, like)
	}
	if q.CategoryID > 0 {
		query = query.Where("category_id = ?", q.CategoryID)
	}
	if q.MinPrice > 0 {
		query = query.Where("price >= ?", q.MinPrice)
	}
	if q.MaxPrice > 0 {
		query = query.Where("price <= ?", q.MaxPrice)
	}
	if q.FreeShipping != nil {
		query = query.Where("free_shipping = ?", *q.FreeShipping)
	}
	if q.ShippingFrom != "" {
		query = query.Where("shipping_from = ?", q.ShippingFrom)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, util.WrapAppError(50001, "count products failed", err)
	}

	page, pageSize := q.Page, q.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	switch q.Sort {
	case "price_asc":
		query = query.Order("price asc, id desc")
	case "price_desc":
		query = query.Order("price desc, id desc")
	case "sales":
		query = query.Order("sales desc, id desc")
	case "rating":
		query = query.Order("rating desc, sales desc, id desc")
	default:
		query = query.Order("id asc")
	}

	var list []model.Product
	if err := query.Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, util.WrapAppError(50001, "search products failed", err)
	}
	return list, total, nil
}

// FindByIDs 批量查询商品（购物车/下单校验用）。
func (r *ProductRepository) FindByIDs(ids []uint) ([]model.Product, error) {
	if len(ids) == 0 {
		return nil, nil
	}
	var list []model.Product
	if err := r.db.Where("id IN ?", ids).Find(&list).Error; err != nil {
		return nil, util.WrapAppError(50001, "find products by ids failed", err)
	}
	return list, nil
}

// IncrementSales 增加销量。
func (r *ProductRepository) IncrementSales(tx *gorm.DB, productID uint, quantity int) error {
	if err := tx.Model(&model.Product{}).Where("id = ?", productID).
		UpdateColumn("sales", gorm.Expr("sales + ?", quantity)).Error; err != nil {
		return util.WrapAppError(50001, "increment product sales failed", err)
	}
	return nil
}

// DeductStock 扣减库存（SELECT ... FOR UPDATE 防并发超卖）。
func (r *ProductRepository) DeductStock(tx *gorm.DB, productID uint, quantity int) error {
	var p model.Product
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&p, productID).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return util.ErrNotFound
		}
		return util.WrapAppError(50001, "lock product failed", err)
	}
	if p.Stock <= quantity {
		return fmt.Errorf("insufficient stock: product=%s available=%d", p.Name, p.Stock)
	}
	if err := tx.Model(&model.Product{}).Where("id = ?", productID).
		UpdateColumn("stock", gorm.Expr("stock - ?", quantity)).Error; err != nil {
		return util.WrapAppError(50001, "deduct stock failed", err)
	}
	return nil
}

// RestoreStock 回补库存（取消订单）。
func (r *ProductRepository) RestoreStock(tx *gorm.DB, productID uint, quantity int) error {
	if err := tx.Model(&model.Product{}).Where("id = ?", productID).
		UpdateColumn("stock", gorm.Expr("stock + ?", quantity)).Error; err != nil {
		return util.WrapAppError(50001, "restore stock failed", err)
	}
	return nil
}

// UpdateRating 更新商品评分（新增评价后）。
func (r *ProductRepository) UpdateRating(tx *gorm.DB, productID uint, rating int) error {
	var p model.Product
	if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).First(&p, productID).Error; err != nil {
		return util.WrapAppError(50001, "lock product for rating failed", err)
	}
	newCount := p.RatingCount + 1
	newRating := (p.Rating*float64(p.RatingCount) + float64(rating)) / float64(newCount)
	if err := tx.Model(&model.Product{}).Where("id = ?", productID).
		Updates(map[string]interface{}{"rating": newRating, "rating_count": newCount}).Error; err != nil {
		return util.WrapAppError(50001, "update product rating failed", err)
	}
	return nil
}

// TopSales 热销商品 TopN。
func (r *ProductRepository) TopSales(limit int) ([]model.Product, error) {
	var list []model.Product
	if err := r.db.Where("status = ?", "ON_SALE").Order("sales desc, id asc").Limit(limit).Find(&list).Error; err != nil {
		return nil, util.WrapAppError(50001, "top sales products failed", err)
	}
	return list, nil
}

// RecommendByCategory 按分类推荐（首页推荐复用 Search 后按分类加权）。
func (r *ProductRepository) RecommendByCategory(categoryID uint, excludeID uint, limit int) ([]model.Product, error) {
	var list []model.Product
	query := r.db.Where("status = ?", "ON_SALE")
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}
	if err := query.Order("sales desc, rating desc, id asc").Limit(limit).Find(&list).Error; err != nil {
		return nil, util.WrapAppError(50001, "recommend products failed", err)
	}
	return list, nil
}

func (r *ProductRepository) CountAll() (int64, error) {
	var count int64
	if err := r.db.Model(&model.Product{}).Count(&count).Error; err != nil {
		return 0, util.WrapAppError(50001, "count all products failed", err)
	}
	return count, nil
}
