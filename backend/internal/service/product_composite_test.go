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
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newProductEnv002(t *testing.T) (*ProductService, *repository.ProductRepository, *gorm.DB) {
	t.Helper()
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}
	sqlDB, _ := db.DB()
	sqlDB.SetMaxOpenConns(1)
	if err := db.AutoMigrate(&model.Product{}, &model.Category{}, &model.AuditLog{}); err != nil {
		t.Fatalf("migrate: %v", err)
	}
	cateRepo := repository.NewCategoryRepository(db)
	_ = cateRepo.Create(&model.Category{Name: "鲜切花", Level: 1})
	prodRepo := repository.NewProductRepository(db)
	auditSvc := NewAuditService(repository.NewAuditRepository(db))
	l := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	rdb := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
	svc := NewProductService(prodRepo, cateRepo, auditSvc, rdb, l)
	return svc, prodRepo, db
}

func TestProductList_UserSeesOnSale(t *testing.T) {
	svc, _, _ := newProductEnv002(t)
	p1, err := svc.Create(dto.ProductCreateRequest{CategoryID: 1, Name: "在售玫瑰", Price: 50, Stock: 10, Status: constants.ProductStatusOnSale})
	if err != nil {
		t.Fatalf("create1: %v", err)
	}
	p2, err := svc.Create(dto.ProductCreateRequest{CategoryID: 1, Name: "下架百合", Price: 60, Stock: 10, Status: constants.ProductStatusOnSale})
	if err != nil {
		t.Fatalf("create2: %v", err)
	}
	if _, err := svc.ChangeStatus(p2.ID, constants.ProductStatusOffSale); err != nil {
		t.Fatalf("change status: %v", err)
	}
	list, total, err := svc.List(dto.ProductQuery{Page: 1, PageSize: 10}, false)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Fatalf("user list total=%d len=%d, want 1", total, len(list))
	}
	if list[0].ID != p1.ID {
		t.Errorf("user should see on-sale product, got %+v", list)
	}
}

func TestProductList_NewestFirst(t *testing.T) {
	svc, _, _ := newProductEnv002(t)
	p1, _ := svc.Create(dto.ProductCreateRequest{CategoryID: 1, Name: "玫瑰1", Price: 50, Stock: 10, Status: constants.ProductStatusOnSale})
	p2, _ := svc.Create(dto.ProductCreateRequest{CategoryID: 1, Name: "玫瑰2", Price: 50, Stock: 10, Status: constants.ProductStatusOnSale})
	list, _, err := svc.List(dto.ProductQuery{Page: 1, PageSize: 10}, true)
	if err != nil {
		t.Fatalf("list: %v", err)
	}
	if len(list) != 2 || list[0].ID != p2.ID {
		t.Errorf("newest first wrong: %+v (p1=%d p2=%d)", list, p1.ID, p2.ID)
	}
}

func TestDeductStock_ExactQty(t *testing.T) {
	_, prodRepo, db := newProductEnv002(t)
	prod := &model.Product{CategoryID: 1, Name: "郁金香", Price: 30, Stock: 5, Status: constants.ProductStatusOnSale}
	if err := db.Create(prod).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := prodRepo.DeductStock(db, prod.ID, 5); err != nil {
		t.Fatalf("deduct exact qty should succeed: %v", err)
	}
}

func TestProductDetail_OffSaleBlocked(t *testing.T) {
	svc, _, db := newProductEnv002(t)
	prod := &model.Product{CategoryID: 1, Name: "玫瑰", Price: 50, Stock: 10, Status: constants.ProductStatusOffSale}
	if err := db.Create(prod).Error; err != nil {
		t.Fatalf("create: %v", err)
	}
	if _, err := svc.Detail(prod.ID); err == nil {
		t.Fatal("expected off-sale product detail to fail")
	}
}

func TestProductStatusText_OnSale(t *testing.T) {
	if got := util.ProductStatusText(constants.ProductStatusOnSale); got != "在售" {
		t.Errorf("ProductStatusText(ON_SALE) = %q, want 在售", got)
	}
}

func TestValidProductStatuses_OnSale(t *testing.T) {
	if !constants.ValidProductStatuses[constants.ProductStatusOnSale] {
		t.Error("ON_SALE should be a valid product status")
	}
}
