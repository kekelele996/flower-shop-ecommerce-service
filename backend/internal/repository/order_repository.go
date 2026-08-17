package repository

import (
	"errors"
	"time"

	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/util"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// OrderRepository 订单仓储。
type OrderRepository struct {
	db *gorm.DB
}

func NewOrderRepository(db *gorm.DB) *OrderRepository {
	return &OrderRepository{db: db}
}

func (r *OrderRepository) Create(tx *gorm.DB, order *model.Order) error {
	if err := tx.Create(order).Error; err != nil {
		return util.WrapAppError(50001, "create order failed", err)
	}
	return nil
}

func (r *OrderRepository) Update(tx *gorm.DB, order *model.Order) error {
	if err := tx.Save(order).Error; err != nil {
		return util.WrapAppError(50001, "update order failed", err)
	}
	return nil
}

func (r *OrderRepository) UpdateStatus(tx *gorm.DB, id uint, status string, fields map[string]interface{}) error {
	updates := map[string]interface{}{"status": status}
	for k, v := range fields {
		updates[k] = v
	}
	if err := tx.Model(&model.Order{}).Where("id = ?", id).Updates(updates).Error; err != nil {
		return util.WrapAppError(50001, "update order status failed", err)
	}
	return nil
}

func (r *OrderRepository) FindByID(tx *gorm.DB, id uint) (*model.Order, error) {
	var order model.Order
	err := tx.Preload("Items").First(&order, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, util.WrapAppError(50001, "find order by id failed", err)
	}
	return &order, nil
}

// FindByIDForUpdate 行锁查询订单（并发状态流转保护）。
func (r *OrderRepository) FindByIDForUpdate(tx *gorm.DB, id uint) (*model.Order, error) {
	var order model.Order
	err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).Preload("Items").First(&order, id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, util.WrapAppError(50001, "lock order failed", err)
	}
	return &order, nil
}

func (r *OrderRepository) FindByOrderNo(tx *gorm.DB, orderNo string) (*model.Order, error) {
	var order model.Order
	err := tx.Preload("Items").Where("order_no = ?", orderNo).First(&order).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, util.ErrNotFound
	}
	if err != nil {
		return nil, util.WrapAppError(50001, "find order by order_no failed", err)
	}
	return &order, nil
}

// List 订单列表（用户/管理端共用，根据 userID 是否 >0 区分）。
func (r *OrderRepository) List(q dto.OrderQuery, userID uint) ([]model.Order, int64, error) {
	query := r.db.Model(&model.Order{})
	if userID > 0 {
		query = query.Where("user_id = ?", userID)
	}
	if q.Status != "" {
		query = query.Where("status = ?", q.Status)
	}
	if q.OrderNo != "" {
		query = query.Where("order_no LIKE ?", "%"+q.OrderNo+"%")
	}
	if q.Keyword != "" {
		query = query.Where("order_no LIKE ?", "%"+q.Keyword+"%")
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, util.WrapAppError(50001, "count orders failed", err)
	}

	page, pageSize := q.Page, q.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	var list []model.Order
	if err := query.Preload("Items").Order("id desc").Offset((page - 1) * pageSize).Limit(pageSize).Find(&list).Error; err != nil {
		return nil, 0, util.WrapAppError(50001, "list orders failed", err)
	}
	return list, total, nil
}

func (r *OrderRepository) CountByStatus(status string) (int64, error) {
	var count int64
	if err := r.db.Model(&model.Order{}).Where("status = ?", status).Count(&count).Error; err != nil {
		return 0, util.WrapAppError(50001, "count orders by status failed", err)
	}
	return count, nil
}

// SumSales 统计总销售额（已支付及以上状态）。
func (r *OrderRepository) SumSales() (float64, error) {
	var sum float64
	if err := r.db.Model(&model.Order{}).
		Where("status IN ?", []string{"PENDING_SHIPMENT", "SHIPPED", "COMPLETED"}).
		Select("COALESCE(SUM(pay_amount), 0)").Scan(&sum).Error; err != nil {
		return 0, util.WrapAppError(50001, "sum sales failed", err)
	}
	return sum, nil
}

// CountByStatusMap 各状态订单数。
func (r *OrderRepository) CountByStatusMap() (map[string]int64, error) {
	type row struct {
		Status string
		Count  int64
	}
	var rows []row
	if err := r.db.Model(&model.Order{}).Select("status, COUNT(*) as count").Group("status").Scan(&rows).Error; err != nil {
		return nil, util.WrapAppError(50001, "count orders group by status failed", err)
	}
	out := map[string]int64{}
	for _, rw := range rows {
		out[rw.Status] = rw.Count
	}
	return out, nil
}

// DailySales 近 7 天每日销售额。
func (r *OrderRepository) DailySales(days int) ([]dto.DailySalesVO, error) {
	type row struct {
		Date  string
		Sales float64
	}
	var rows []row
	since := time.Now().AddDate(0, 0, -(days - 1))
	if err := r.db.Model(&model.Order{}).
		Where("status IN ? AND created_at >= ?", []string{"PENDING_SHIPMENT", "SHIPPED", "COMPLETED"}, since).
		Select("TO_CHAR(created_at, 'YYYY-MM-DD') as date, SUM(pay_amount) as sales").
		Group("TO_CHAR(created_at, 'YYYY-MM-DD')").Order("date asc").Scan(&rows).Error; err != nil {
		return nil, util.WrapAppError(50001, "daily sales failed", err)
	}
	out := make([]dto.DailySalesVO, 0, len(rows))
	for _, rw := range rows {
		out = append(out, dto.DailySalesVO{Date: rw.Date, Sales: rw.Sales})
	}
	return out, nil
}

// TopProductsByRevenue 热销商品 TopN（按成交额）。
func (r *OrderRepository) TopProductsByRevenue(limit int) ([]dto.TopProductVO, error) {
	type row struct {
		ProductID   uint
		ProductName string
		Sales       int
		Revenue     float64
	}
	var rows []row
	if err := r.db.Model(&model.OrderItem{}).
		Joins("JOIN orders ON orders.id = order_items.order_id").
		Where("orders.status IN ?", []string{"PENDING_SHIPMENT", "SHIPPED", "COMPLETED"}).
		Select("order_items.product_id, MAX(order_items.product_name) as product_name, SUM(order_items.quantity) as sales, SUM(order_items.total_price) as revenue").
		Group("order_items.product_id").Order("revenue desc").Limit(limit).Scan(&rows).Error; err != nil {
		return nil, util.WrapAppError(50001, "top products failed", err)
	}
	out := make([]dto.TopProductVO, 0, len(rows))
	for _, rw := range rows {
		out = append(out, dto.TopProductVO{ProductID: rw.ProductID, ProductName: rw.ProductName, Sales: rw.Sales, Revenue: rw.Revenue})
	}
	return out, nil
}

func (r *OrderRepository) CountAll() (int64, error) {
	var count int64
	if err := r.db.Model(&model.Order{}).Count(&count).Error; err != nil {
		return 0, util.WrapAppError(50001, "count all orders failed", err)
	}
	return count, nil
}
