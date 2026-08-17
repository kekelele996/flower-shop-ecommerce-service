package service

import (
	"log/slog"
	"os"
	"testing"
	"time"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/repository"
	"github.com/flowershop/backend/internal/util"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newOrderEnv005(t *testing.T) (*OrderService, *repository.OrderRepository, *model.User, uint) {
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
	order := &model.Order{OrderNo: "O-005", UserID: user.ID, Status: constants.OrderStatusShipped, TotalAmount: 100, PayAmount: 100}
	if err := db.Create(order).Error; err != nil {
		t.Fatalf("create order: %v", err)
	}
	if err := db.Create(&model.OrderItem{OrderID: order.ID, ProductID: 1, ProductName: "花", Price: 100, Quantity: 1, TotalPrice: 100}).Error; err != nil {
		t.Fatalf("create item: %v", err)
	}
	events := util.MarshalJSON([]model.LogisticsEvent{{Time: time.Now(), Status: constants.LogisticsStatusPickedUp, Desc: "已揽收"}})
	if err := db.Create(&model.Logistics{OrderID: order.ID, OrderNo: "O-005", TrackingNo: "SF001", Carrier: "SF_EXPRESS", Status: constants.LogisticsStatusPickedUp, Events: events}).Error; err != nil {
		t.Fatalf("create logistics: %v", err)
	}
	auditSvc := NewAuditService(repository.NewAuditRepository(db))
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	orderSvc := NewOrderService(db, repository.NewOrderRepository(db), repository.NewOrderItemRepository(db), repository.NewPaymentRepository(db), repository.NewLogisticsRepository(db), repository.NewCartRepository(db), repository.NewProductRepository(db), repository.NewCouponRepository(db), auditSvc, l)
	return orderSvc, repository.NewOrderRepository(db), user, order.ID
}

func TestOrderDetail_HasItems(t *testing.T) {
	svc, _, user, orderID := newOrderEnv005(t)
	vo, err := svc.Detail(user.ID, orderID, false)
	if err != nil {
		t.Fatalf("detail: %v", err)
	}
	if len(vo.Items) != 1 {
		t.Errorf("detail items = %d, want 1", len(vo.Items))
	}
}

func TestFindByOrderNo_HasItems(t *testing.T) {
	svc, orderRepo, _, _ := newOrderEnv005(t)
	order, err := orderRepo.FindByOrderNo(svc.db, "O-005")
	if err != nil {
		t.Fatalf("find by order no: %v", err)
	}
	if len(order.Items) != 1 {
		t.Errorf("order items = %d, want 1", len(order.Items))
	}
}

func TestLogistics_HasEvents(t *testing.T) {
	svc, _, user, orderID := newOrderEnv005(t)
	vo, err := svc.Logistics(user.ID, orderID, false)
	if err != nil {
		t.Fatalf("logistics: %v", err)
	}
	if len(vo.Events) == 0 {
		t.Fatal("expected logistics events, got none")
	}
}
