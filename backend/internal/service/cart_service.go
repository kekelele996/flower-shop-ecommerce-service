package service

import (
	"errors"
	"log/slog"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/repository"
	"github.com/flowershop/backend/internal/util"
)

// CartService 购物车服务。
type CartService struct {
	repo     *repository.CartRepository
	prodRepo *repository.ProductRepository
	auditSvc *AuditService
	logger   *slog.Logger
}

func NewCartService(repo *repository.CartRepository, prodRepo *repository.ProductRepository, auditSvc *AuditService, logger *slog.Logger) *CartService {
	return &CartService{repo: repo, prodRepo: prodRepo, auditSvc: auditSvc, logger: logger}
}

func (s *CartService) Add(userID uint, req dto.CartAddRequest) (*dto.CartItemVO, error) {
	product, err := s.prodRepo.FindByID(req.ProductID)
	if err != nil {
		return nil, util.NewAppError(constants.CodeProductNotFound, "product not found, id="+util.UintString(req.ProductID))
	}
	if product.Status != constants.ProductStatusOnSale {
		return nil, util.NewAppError(constants.CodeProductNotFound, "product="+product.Name+" is not on sale")
	}
	if product.Stock < req.Quantity {
		return nil, util.NewAppError(constants.CodeInsufficientStock, constants.MsgInsufficientStock)
	}
	item, err := s.repo.FindByUserAndProduct(userID, req.ProductID)
	if errors.Is(err, util.ErrNotFound) {
		item = &model.CartItem{UserID: userID, ProductID: req.ProductID, Quantity: req.Quantity, Selected: true}
		if err := s.repo.Create(item); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	} else {
		newQty := item.Quantity + req.Quantity
		if newQty > constants.CartMaxQuantity {
			newQty = constants.CartMaxQuantity
		}
		if product.Stock < newQty {
			return nil, util.NewAppError(constants.CodeInsufficientStock, constants.MsgInsufficientStock)
		}
		item.Quantity = newQty
		if err := s.repo.Update(item); err != nil {
			return nil, err
		}
	}
	s.logger.Info(constants.LogCartItemAdded, "user_id", userID, "product_id", req.ProductID, "quantity", req.Quantity)
	s.auditSvc.Record(userID, "", "CREATE", "cart_item", uint64(item.ID), "add product to cart")
	return s.toVO(item, product), nil
}

func (s *CartService) Update(userID, itemID uint, req dto.CartUpdateRequest) error {
	item, err := s.repo.FindByID(itemID)
	if err != nil {
		return util.NewAppError(constants.CodeCartItemNotFound, "cart item not found, id="+util.UintString(itemID))
	}
	if item.UserID != userID {
		return util.NewAppError(constants.CodeForbidden, "cart item="+util.UintString(itemID)+" does not belong to user="+util.UintString(userID))
	}
	if req.Quantity != nil {
		product, err := s.prodRepo.FindByID(item.ProductID)
		if err != nil {
			return err
		}
		if product.Stock < *req.Quantity {
			return util.NewAppError(constants.CodeInsufficientStock, constants.MsgInsufficientStock)
		}
		item.Quantity = *req.Quantity
	}
	if req.Selected != nil {
		item.Selected = *req.Selected
	}
	if err := s.repo.Update(item); err != nil {
		return err
	}
	s.logger.Info(constants.LogCartItemUpdated, "user_id", userID, "item_id", itemID, "quantity", item.Quantity)
	return nil
}

func (s *CartService) Remove(userID, itemID uint) error {
	item, err := s.repo.FindByID(itemID)
	if err != nil {
		return util.NewAppError(constants.CodeCartItemNotFound, "cart item not found, id="+util.UintString(itemID))
	}
	if item.UserID != userID {
		return util.NewAppError(constants.CodeForbidden, "cart item="+util.UintString(itemID)+" does not belong to user="+util.UintString(userID))
	}
	if err := s.repo.Delete(itemID); err != nil {
		return err
	}
	s.logger.Info(constants.LogCartItemRemoved, "user_id", userID, "item_id", itemID)
	return nil
}

// Summary 购物车汇总（复用 cartRepository.FindByUserID）。
func (s *CartService) Summary(userID uint) (*dto.CartSummaryVO, error) {
	items, err := s.repo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	vo := &dto.CartSummaryVO{Items: make([]dto.CartItemVO, 0, len(items))}
	for _, item := range items {
		if item.Product == nil {
			continue
		}
		vo.Items = append(vo.Items, *s.toVO(&item, item.Product))
		vo.TotalQuantity += item.Quantity
		vo.TotalPrice += float64(item.Quantity) * item.Product.Price
	}
	return vo, nil
}

func (s *CartService) toVO(item *model.CartItem, product *model.Product) *dto.CartItemVO {
	return &dto.CartItemVO{
		ID:        item.ID,
		UserID:    item.UserID,
		ProductID: item.ProductID,
		Quantity:  item.Quantity,
		Selected:  item.Selected,
		Subtotal:  float64(item.Quantity) * product.Price,
		Product:   *toProductVO(product),
	}
}
