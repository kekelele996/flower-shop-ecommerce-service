package service

import (
	"log/slog"
	"os"
	"testing"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/repository"
	"github.com/flowershop/backend/internal/util"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newOrderEnv003(t *testing.T) (*OrderService, *repository.OrderRepository, *repository.ProductRepository, *model.User, uint, uint) {
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
	prod := &model.Product{CategoryID: 1, Name: "香槟玫瑰", Price: 100, Stock: 100, FreeShipping: true, Status: constants.ProductStatusOnSale}
	if err := db.Create(prod).Error; err != nil {
		t.Fatalf("create product: %v", err)
	}
	order := &model.Order{OrderNo: "O-001", UserID: user.ID, Status: constants.OrderStatusPendingPayment, TotalAmount: 200, PayAmount: 200}
	if err := db.Create(order).Error; err != nil {
		t.Fatalf("create order: %v", err)
	}
	item := &model.OrderItem{OrderID: order.ID, ProductID: prod.ID, ProductName: "香槟玫瑰", Price: 100, Quantity: 2, TotalPrice: 200}
	if err := db.Create(item).Error; err != nil {
		t.Fatalf("create order item: %v", err)
	}
	auditSvc := NewAuditService(repository.NewAuditRepository(db))
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	orderSvc := NewOrderService(db, repository.NewOrderRepository(db), repository.NewOrderItemRepository(db), repository.NewPaymentRepository(db), repository.NewLogisticsRepository(db), repository.NewCartRepository(db), repository.NewProductRepository(db), repository.NewCouponRepository(db), auditSvc, l)
	return orderSvc, repository.NewOrderRepository(db), repository.NewProductRepository(db), user, order.ID, prod.ID
}

func TestPay_UpdatesStatus(t *testing.T) {
	svc, orderRepo, _, user, orderID, _ := newOrderEnv003(t)
	if _, err := svc.Pay(user.ID, orderID, dto.PayRequest{Channel: "ALIPAY_SANDBOX"}); err != nil {
		t.Fatalf("pay: %v", err)
	}
	order, err := orderRepo.FindByID(svc.db, orderID)
	if err != nil {
		t.Fatalf("find order: %v", err)
	}
	if order.Status != constants.OrderStatusPendingShipment {
		t.Errorf("status = %s, want PENDING_SHIPMENT", order.Status)
	}
}

func TestPay_IncrementsSales(t *testing.T) {
	svc, _, prodRepo, user, orderID, prodID := newOrderEnv003(t)
	if _, err := svc.Pay(user.ID, orderID, dto.PayRequest{Channel: "ALIPAY_SANDBOX"}); err != nil {
		t.Fatalf("pay: %v", err)
	}
	prod, err := prodRepo.FindByID(prodID)
	if err != nil {
		t.Fatalf("find product: %v", err)
	}
	if prod.Sales != 2 {
		t.Errorf("sales = %d, want 2", prod.Sales)
	}
}

func TestCancel_ShippedBlocked(t *testing.T) {
	svc, _, _, user, orderID, _ := newOrderEnv003(t)
	if err := svc.db.Model(&model.Order{}).Where("id = ?", orderID).Update("status", constants.OrderStatusShipped).Error; err != nil {
		t.Fatalf("force shipped: %v", err)
	}
	if _, err := svc.Cancel(user.ID, orderID, "不想要了"); err == nil {
		t.Fatal("expected cancel of shipped order to be blocked")
	}
}

func TestOrderStatusText_PendingPayment(t *testing.T) {
	if got := util.OrderStatusText("PENDING_PAYMENT"); got != "待付款" {
		t.Errorf("OrderStatusText(PENDING_PAYMENT) = %q, want 待付款", got)
	}
}

func TestValidOrderStatuses_Shipped(t *testing.T) {
	if !constants.ValidOrderStatuses[constants.OrderStatusShipped] {
		t.Error("SHIPPED should be a valid order status")
	}
}
