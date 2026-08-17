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

func newCartEnv001(t *testing.T) (*CartService, *model.User, *model.Product, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&model.User{}, &model.Product{}, &model.Category{}, &model.CartItem{}, &model.AuditLog{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	user := &model.User{Username: "u1", Password: "h", Status: constants.UserStatusActive}
	if err := db.Create(user).Error; err != nil {
		t.Fatalf("create user: %v", err)
	}
	prod := &model.Product{CategoryID: 1, Name: "玫瑰", Price: 10, Stock: 999, FreeShipping: true, Status: constants.ProductStatusOnSale}
	if err := db.Create(prod).Error; err != nil {
		t.Fatalf("create product: %v", err)
	}
	auditSvc := NewAuditService(repository.NewAuditRepository(db))
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	svc := NewCartService(repository.NewCartRepository(db), repository.NewProductRepository(db), auditSvc, l)
	return svc, user, prod, db
}

func TestCartSummary_TotalAndOrder(t *testing.T) {
	svc, user, prod, db := newCartEnv001(t)
	if _, err := svc.Add(user.ID, dto.CartAddRequest{ProductID: prod.ID, Quantity: 2}); err != nil {
		t.Fatalf("add1: %v", err)
	}
	prod2 := &model.Product{CategoryID: 1, Name: "百合", Price: 20, Stock: 999, FreeShipping: true, Status: constants.ProductStatusOnSale}
	if err := db.Create(prod2).Error; err != nil {
		t.Fatalf("create product2: %v", err)
	}
	if _, err := svc.Add(user.ID, dto.CartAddRequest{ProductID: prod2.ID, Quantity: 3}); err != nil {
		t.Fatalf("add2: %v", err)
	}
	sum, err := svc.Summary(user.ID)
	if err != nil {
		t.Fatalf("summary: %v", err)
	}
	if sum.TotalQuantity != 5 {
		t.Errorf("total quantity = %d, want 5", sum.TotalQuantity)
	}
	if sum.TotalPrice != 80 {
		t.Errorf("total price = %v, want 80", sum.TotalPrice)
	}
	if len(sum.Items) != 2 || sum.Items[0].ProductID != prod2.ID {
		t.Errorf("summary order wrong: %+v", sum.Items)
	}
}

func TestCartAdd_CapsAt99(t *testing.T) {
	svc, user, prod, _ := newCartEnv001(t)
	for i := 0; i < 3; i++ {
		if _, err := svc.Add(user.ID, dto.CartAddRequest{ProductID: prod.ID, Quantity: 40}); err != nil {
			t.Fatalf("add: %v", err)
		}
	}
	item, err := svc.repo.FindByUserAndProduct(user.ID, prod.ID)
	if err != nil {
		t.Fatalf("find: %v", err)
	}
	if item.Quantity != 99 {
		t.Errorf("quantity = %d, want 99", item.Quantity)
	}
}

func TestPriceText_Decimals(t *testing.T) {
	if got := util.PriceText(129.5); got != "¥129.50" {
		t.Errorf("PriceText(129.5) = %q, want ¥129.50", got)
	}
}

func TestOffSaleStatusValid(t *testing.T) {
	if !constants.ValidProductStatuses[constants.ProductStatusOffSale] {
		t.Error("OFF_SALE should be a valid product status")
	}
}
