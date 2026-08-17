package repository

import (
	"testing"

	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func newTestDB(t *testing.T) *gorm.DB {
	t.Helper()
	db, err := gorm.Open(sqlite.Open("file::memory:?cache=shared"), &gorm.Config{Logger: logger.Default.LogMode(logger.Silent)})
	if err != nil {
		t.Fatalf("open sqlite failed: %v", err)
	}
	if err := db.AutoMigrate(&model.Product{}, &model.Category{}); err != nil {
		t.Fatalf("migrate failed: %v", err)
	}
	return db
}

func TestProductRepositorySearch(t *testing.T) {
	db := newTestDB(t)
	repo := NewProductRepository(db)
	_ = db.Create(&model.Category{ID: 1, Name: "鲜切花", Level: 1})
	products := []model.Product{
		{CategoryID: 1, Name: "玫瑰", Price: 100, Stock: 10, Sales: 5, FreeShipping: true, Status: "ON_SALE", ShippingFrom: "云南"},
		{CategoryID: 1, Name: "向日葵", Price: 50, Stock: 20, Sales: 8, FreeShipping: false, Status: "ON_SALE", ShippingFrom: "云南"},
		{CategoryID: 2, Name: "绿萝", Price: 30, Stock: 30, Sales: 20, FreeShipping: true, Status: "OFF_SALE", ShippingFrom: "福建"},
	}
	for i := range products {
		if err := repo.Create(&products[i]); err != nil {
			t.Fatalf("create product failed: %v", err)
		}
	}

	tests := []struct {
		name       string
		q          dto.ProductQuery
		onlyOnSale bool
		want       int64
	}{
		{"keyword rose", dto.ProductQuery{Keyword: "玫瑰", Page: 1, PageSize: 10}, false, 1},
		{"on sale only", dto.ProductQuery{Page: 1, PageSize: 10}, true, 2},
		{"category 1", dto.ProductQuery{CategoryID: 1, Page: 1, PageSize: 10}, false, 2},
		{"free shipping", func() dto.ProductQuery { fs := true; return dto.ProductQuery{FreeShipping: &fs, Page: 1, PageSize: 10} }(), true, 1},
		{"price range", dto.ProductQuery{MinPrice: 40, MaxPrice: 60, Page: 1, PageSize: 10}, false, 1},
		{"shipping from", dto.ProductQuery{ShippingFrom: "云南", Page: 1, PageSize: 10}, false, 2},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			list, total, err := repo.Search(tt.q, tt.onlyOnSale)
			if err != nil {
				t.Fatalf("search failed: %v", err)
			}
			if total != tt.want {
				t.Errorf("total = %d, want %d (list=%+v)", total, tt.want, list)
			}
		})
	}
}

func TestProductRepositoryFindByID(t *testing.T) {
	db := newTestDB(t)
	repo := NewProductRepository(db)
	p := &model.Product{CategoryID: 1, Name: "龟背竹", Price: 69, Stock: 10, Status: "ON_SALE"}
	if err := repo.Create(p); err != nil {
		t.Fatalf("create failed: %v", err)
	}
	got, err := repo.FindByID(p.ID)
	if err != nil {
		t.Fatalf("find failed: %v", err)
	}
	if got.Name != "龟背竹" {
		t.Errorf("name = %q", got.Name)
	}
}
