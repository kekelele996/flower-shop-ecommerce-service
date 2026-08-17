package service

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/repository"
	"github.com/flowershop/backend/internal/util"
	"github.com/redis/go-redis/v9"
)

// ProductService 商品服务：列表、搜索、推荐、上下架、浏览记录。
type ProductService struct {
	repo     *repository.ProductRepository
	cateRepo *repository.CategoryRepository
	auditSvc *AuditService
	redis    *redis.Client
	logger   *slog.Logger
}

func NewProductService(repo *repository.ProductRepository, cateRepo *repository.CategoryRepository, auditSvc *AuditService, rdb *redis.Client, logger *slog.Logger) *ProductService {
	return &ProductService{repo: repo, cateRepo: cateRepo, auditSvc: auditSvc, redis: rdb, logger: logger}
}

func (s *ProductService) Create(req dto.ProductCreateRequest) (*dto.ProductVO, error) {
	p := &model.Product{
		CategoryID:    req.CategoryID,
		Name:          req.Name,
		SubTitle:      req.SubTitle,
		Description:   req.Description,
		Detail:        req.Detail,
		CoverImage:    req.CoverImage,
		Images:        util.MarshalJSON(req.Images),
		Price:         req.Price,
		OriginalPrice: req.OriginalPrice,
		Stock:         req.Stock,
		ShippingFrom:  req.ShippingFrom,
		FreeShipping:  req.FreeShipping,
		Status:        constants.ProductStatusOnSale,
	}
	if req.Status != "" {
		p.Status = req.Status
	}
	if err := s.repo.Create(p); err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogProductCreated, "id", p.ID, "name", p.Name, "price", p.Price, "stock", p.Stock, "status", p.Status)
	s.auditSvc.Record(0, "admin", "CREATE", "product", uint64(p.ID), "create product "+p.Name)
	s.invalidateCache()
	return s.ToVO(p), nil
}

func (s *ProductService) Update(id uint, req dto.ProductUpdateRequest) (*dto.ProductVO, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.NewAppError(constants.CodeProductNotFound, "product not found, id="+util.UintString(id))
	}
	if req.CategoryID != nil {
		p.CategoryID = *req.CategoryID
	}
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.SubTitle != nil {
		p.SubTitle = *req.SubTitle
	}
	if req.Description != nil {
		p.Description = *req.Description
	}
	if req.Detail != nil {
		p.Detail = *req.Detail
	}
	if req.CoverImage != nil {
		p.CoverImage = *req.CoverImage
	}
	if req.Images != nil {
		p.Images = util.MarshalJSON(req.Images)
	}
	if req.Price != nil {
		p.Price = *req.Price
	}
	if req.OriginalPrice != nil {
		p.OriginalPrice = *req.OriginalPrice
	}
	if req.Stock != nil {
		p.Stock = *req.Stock
	}
	if req.ShippingFrom != nil {
		p.ShippingFrom = *req.ShippingFrom
	}
	if req.FreeShipping != nil {
		p.FreeShipping = *req.FreeShipping
	}
	if req.Status != nil {
		p.Status = *req.Status
	}
	if err := s.repo.Update(p); err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogProductUpdated, "id", p.ID, "name", p.Name, "price", p.Price, "stock", p.Stock, "status", p.Status)
	s.auditSvc.Record(0, "admin", "UPDATE", "product", uint64(p.ID), "update product "+p.Name)
	s.invalidateCache()
	return s.ToVO(p), nil
}

// ChangeStatus 商品上下架。
func (s *ProductService) ChangeStatus(id uint, status string) (*dto.ProductVO, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.NewAppError(constants.CodeProductNotFound, "product not found, id="+util.UintString(id))
	}
	p.Status = status
	if err := s.repo.Update(p); err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogProductStatusChanged, "id", p.ID, "status", p.Status, "operator", 0)
	s.auditSvc.Record(0, "admin", "UPDATE", "product", uint64(p.ID), "change product status to "+status)
	s.invalidateCache()
	return s.ToVO(p), nil
}

// List 商品列表（管理端可含下架商品）。
func (s *ProductService) List(q dto.ProductQuery, admin bool) ([]dto.ProductVO, int64, error) {
	list, total, err := s.repo.Search(q, !admin)
	if err != nil {
		return nil, 0, err
	}
	return s.ToVOs(list), total, nil
}

// Detail 商品详情。
func (s *ProductService) Detail(id uint) (*dto.ProductVO, error) {
	p, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.NewAppError(constants.CodeProductNotFound, "product not found, id="+util.UintString(id))
	}
	if p.Status != constants.ProductStatusOnSale {
		return nil, util.NewAppError(constants.CodeProductNotFound, "product not on sale, id="+util.UintString(id))
	}
	return s.ToVO(p), nil
}

// RecordView 记录浏览行为（Redis 存储最近浏览，用于首页推荐）。
func (s *ProductService) RecordView(ctx context.Context, userID, productID uint) {
	key := fmt.Sprintf("flowershop:view:%d", userID)
	s.redis.ZAdd(ctx, key, redis.Z{Score: float64(time.Now().Unix()), Member: productID})
	s.redis.ZRemRangeByRank(ctx, key, 0, -101)
	s.redis.Expire(ctx, key, 7*24*time.Hour)
	s.logger.Info(constants.LogProductViewed, "id", productID, "user_id", userID)
}

// Recommendations 首页推荐：基于最近浏览分类 + 热销兜底。
func (s *ProductService) Recommendations(ctx context.Context, userID uint) ([]dto.ProductVO, error) {
	cacheKey := "flowershop:recommend"
	if userID > 0 {
		cacheKey = fmt.Sprintf("flowershop:recommend:%d", userID)
	}
	cached, err := s.redis.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var vos []dto.ProductVO
		if util.UnmarshalJSONTo(cached, &vos) == nil && len(vos) > 0 {
			return vos, nil
		}
	}

	categoryID := uint(0)
	if userID > 0 {
		ids, _ := s.redis.ZRevRange(ctx, fmt.Sprintf("flowershop:view:%d", userID), 0, 5).Result()
		for _, member := range ids {
			var pid uint
			fmt.Sscanf(member, "%d", &pid)
			if p, err := s.repo.FindByID(pid); err == nil && p.Status == constants.ProductStatusOnSale {
				categoryID = p.CategoryID
				break
			}
		}
	}

	list, err := s.repo.RecommendByCategory(categoryID, 0, 12)
	if err != nil {
		return nil, err
	}
	vos := s.ToVOs(list)
	if len(vos) > 0 {
		s.redis.Set(ctx, cacheKey, util.MarshalJSON(vos), 5*time.Minute)
	}
	return vos, nil
}

// Stats 基础统计（复用 Search 获取商品数量）。
func (s *ProductService) invalidateCache() {
	ctx := context.Background()
	s.redis.Del(ctx, "flowershop:recommend")
}

func (s *ProductService) ToVO(p *model.Product) *dto.ProductVO {
	categoryName := ""
	if c, err := s.cateRepo.FindByID(p.CategoryID); err == nil {
		categoryName = c.Name
	}
	return &dto.ProductVO{
		ID:            p.ID,
		CategoryID:    p.CategoryID,
		CategoryName:  categoryName,
		Name:          p.Name,
		SubTitle:      p.SubTitle,
		Description:   p.Description,
		Detail:        p.Detail,
		CoverImage:    p.CoverImage,
		Images:        util.UnmarshalStrings(p.Images),
		Price:         p.Price,
		OriginalPrice: p.OriginalPrice,
		Stock:         p.Stock,
		Sales:         p.Sales,
		Rating:        p.Rating,
		RatingCount:   p.RatingCount,
		ShippingFrom:  p.ShippingFrom,
		FreeShipping:  p.FreeShipping,
		Status:        p.Status,
		CreatedAt:     util.FormatTime(p.CreatedAt),
	}
}

func (s *ProductService) ToVOs(list []model.Product) []dto.ProductVO {
	out := make([]dto.ProductVO, 0, len(list))
	for i := range list {
		out = append(out, *s.ToVO(&list[i]))
	}
	return out
}
