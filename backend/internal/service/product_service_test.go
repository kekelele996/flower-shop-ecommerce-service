package service

import (
	"context"
	"log/slog"
	"os"
	"testing"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/repository"
	"github.com/glebarez/sqlite"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newTestProductService(t *testing.T) *ProductService {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.AutoMigrate(&model.Product{}, &model.Category{}, &model.AuditLog{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	prodRepo := repository.NewProductRepository(db)
	cateRepo := repository.NewCategoryRepository(db)
	auditRepo := repository.NewAuditRepository(db)
	auditSvc := NewAuditService(auditRepo)
	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelError}))
	rdb := redis.NewClient(&redis.Options{Addr: "127.0.0.1:1"})
	return NewProductService(prodRepo, cateRepo, auditSvc, rdb, logger)
}

func TestProductServiceCreateAndList(t *testing.T) {
	svc := newTestProductService(t)
	_ = svc.cateRepo.Create(&model.Category{Name: "鲜切花", Level: 1})
	cats, _ := svc.cateRepo.ListAll()

	req := dto.ProductCreateRequest{
		CategoryID:   cats[0].ID,
		Name:         "香槟玫瑰",
		Price:        129,
		Stock:        100,
		ShippingFrom: "云南",
		FreeShipping: true,
	}
	created, err := svc.Create(req)
	if err != nil {
		t.Fatalf("create failed: %v", err)
	}
	if created.Name != "香槟玫瑰" || created.Status != constants.ProductStatusOnSale {
		t.Errorf("create mismatch: %+v", created)
	}

	list, total, err := svc.List(dto.ProductQuery{Keyword: "玫瑰", Page: 1, PageSize: 10}, false)
	if err != nil {
		t.Fatalf("list failed: %v", err)
	}
	if total != 1 || len(list) != 1 {
		t.Errorf("total=%d len=%d", total, len(list))
	}

	// 下架后非管理员不可见
	if _, err := svc.ChangeStatus(created.ID, constants.ProductStatusOffSale); err != nil {
		t.Fatalf("change status failed: %v", err)
	}
	if _, err := svc.Detail(created.ID); err == nil {
		t.Error("off-sale product detail should fail")
	}
	_, adminTotal, err := svc.List(dto.ProductQuery{Page: 1, PageSize: 10}, true)
	if err != nil {
		t.Fatalf("admin list failed: %v", err)
	}
	if adminTotal != 1 {
		t.Errorf("admin should see off-sale product, total=%d", adminTotal)
	}
}

func TestProductServiceRecommendations(t *testing.T) {
	svc := newTestProductService(t)
	_ = svc.cateRepo.Create(&model.Category{Name: "盆栽", Level: 1})
	cats, _ := svc.cateRepo.ListAll()
	for _, name := range []string{"琴叶榕", "龟背竹"} {
		if _, err := svc.Create(dto.ProductCreateRequest{CategoryID: cats[0].ID, Name: name, Price: 100, Stock: 10}); err != nil {
			t.Fatalf("create %s failed: %v", name, err)
		}
	}
	list, err := svc.Recommendations(context.Background(), 0)
	if err != nil {
		t.Fatalf("recommendations failed: %v", err)
	}
	if len(list) == 0 {
		t.Error("recommendations should not be empty")
	}
}
