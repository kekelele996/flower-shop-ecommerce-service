package service

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/repository"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newStatsEnv004(t *testing.T) (*OrderService, *repository.OrderRepository, *model.User, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&model.User{}, &model.Product{}, &model.Category{}, &model.CartItem{}, &model.Order{}, &model.OrderItem{}, &model.Payment{}, &model.Logistics{}, &model.Coupon{}, &model.AuditLog{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	user := &model.User{Username: "u1", Password: "h", Status: constants.UserStatusActive}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	orderRepo := repository.NewOrderRepository(db)
	auditSvc := NewAuditService(repository.NewAuditRepository(db))
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	orderSvc := NewOrderService(db, orderRepo, repository.NewOrderItemRepository(db), repository.NewPaymentRepository(db), repository.NewLogisticsRepository(db), repository.NewCartRepository(db), repository.NewProductRepository(db), repository.NewCouponRepository(db), auditSvc, l)
	return orderSvc, orderRepo, user, db
}

func TestSumSales_IncludesCompleted(t *testing.T) {
	_, orderRepo, user, db := newStatsEnv004(t)
	now := time.Now()
	rows := []*model.Order{
		{OrderNo: "A", UserID: user.ID, Status: constants.OrderStatusPendingShipment, PayAmount: 100, PaidAt: &now},
		{OrderNo: "B", UserID: user.ID, Status: constants.OrderStatusShipped, PayAmount: 50, PaidAt: &now},
		{OrderNo: "C", UserID: user.ID, Status: constants.OrderStatusCompleted, PayAmount: 30, PaidAt: &now},
	}
	for _, o := range rows {
		if err := orderRepo.Create(db, o); err != nil {
			t.Fatalf("create order %s: %v", o.OrderNo, err)
		}
	}
	sum, err := orderRepo.SumSales()
	if err != nil {
		t.Fatalf("sum sales: %v", err)
	}
	if sum != 180 {
		t.Errorf("sum sales = %v, want 180", sum)
	}
}

func TestOrderList_UserScoped(t *testing.T) {
	svc, orderRepo, user, db := newStatsEnv004(t)
	other := &model.User{Username: "u2", Password: "h", Status: constants.UserStatusActive}
	if err := db.Create(other).Error; err != nil {
		t.Fatalf("create other user: %v", err)
	}
	for _, u := range []*model.User{user, other} {
		if err := orderRepo.Create(db, &model.Order{OrderNo: "O-" + u.Username, UserID: u.ID, Status: constants.OrderStatusPendingShipment, PayAmount: 10}); err != nil {
			t.Fatalf("create order: %v", err)
		}
	}
	list, _, err := svc.List(dto.OrderQuery{Page: 1, PageSize: 10}, user.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 || list[0].UserID != user.ID {
		t.Errorf("user should only see own orders: %+v", list)
	}
}

func TestOrderList_HasItems(t *testing.T) {
	svc, orderRepo, user, db := newStatsEnv004(t)
	o := &model.Order{OrderNo: "L1", UserID: user.ID, Status: constants.OrderStatusPendingShipment, PayAmount: 10}
	if err := orderRepo.Create(db, o); err != nil {
		t.Fatalf("create order: %v", err)
	}
	if err := db.Create(&model.OrderItem{OrderID: o.ID, ProductID: 1, ProductName: "花", Price: 10, Quantity: 1, TotalPrice: 10}).Error; err != nil {
		t.Fatalf("create item: %v", err)
	}
	list, _, err := svc.List(dto.OrderQuery{Page: 1, PageSize: 10}, user.ID)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 1 || len(list[0].Items) != 1 {
		t.Errorf("list items wrong: %+v", list)
	}
}
