package service

import (
	"log/slog"
	"time"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/repository"
	"github.com/flowershop/backend/internal/util"
	"gorm.io/gorm"
)

// ReviewService 评价服务：完成后评价、商家回复。
type ReviewService struct {
	db        *gorm.DB
	repo      *repository.ReviewRepository
	itemRepo  *repository.OrderItemRepository
	orderRepo *repository.OrderRepository
	prodRepo  *repository.ProductRepository
	userRepo  *repository.UserRepository
	auditSvc  *AuditService
	logger    *slog.Logger
}

func NewReviewService(db *gorm.DB, repo *repository.ReviewRepository, itemRepo *repository.OrderItemRepository, orderRepo *repository.OrderRepository, prodRepo *repository.ProductRepository, userRepo *repository.UserRepository, auditSvc *AuditService, logger *slog.Logger) *ReviewService {
	return &ReviewService{db: db, repo: repo, itemRepo: itemRepo, orderRepo: orderRepo, prodRepo: prodRepo, userRepo: userRepo, auditSvc: auditSvc, logger: logger}
}

// Create 提交评价：订单已完成且该明细未评价。
func (s *ReviewService) Create(userID uint, req dto.ReviewCreateRequest) (*dto.ReviewVO, error) {
	var review *model.Review
	err := s.db.Transaction(func(tx *gorm.DB) error {
		item, err := s.itemRepo.FindByID(req.OrderItemID)
		if err != nil {
			return util.NewAppError(constants.CodeNotFound, "order item not found, id="+util.UintString(req.OrderItemID))
		}
		order, err := s.orderRepo.FindByID(tx, item.OrderID)
		if err != nil {
			return err
		}
		if order.UserID != userID {
			return util.NewAppError(constants.CodeForbidden, "order="+order.OrderNo+" does not belong to user="+util.UintString(userID))
		}
		if order.Status != constants.OrderStatusCompleted {
			return util.NewAppError(constants.CodeReviewNotAllow, constants.MsgReviewNotAllow)
		}
		if item.Reviewed {
			return util.NewAppError(constants.CodeConflict, "order item="+util.UintString(req.OrderItemID)+" already reviewed")
		}
		review = &model.Review{
			UserID:      userID,
			OrderID:     order.ID,
			OrderItemID: item.ID,
			ProductID:   item.ProductID,
			Rating:      req.Rating,
			Content:     req.Content,
			Images:      util.MarshalJSON(req.Images),
		}
		if err := s.repo.Create(tx, review); err != nil {
			return err
		}
		if err := s.itemRepo.MarkReviewed(tx, item.ID); err != nil {
			return err
		}
		return s.prodRepo.UpdateRating(tx, item.ProductID, req.Rating)
	})
	if err != nil {
		return nil, err
	}
	user, _ := s.userRepo.FindByID(userID)
	username := ""
	if user != nil {
		username = user.Username
	}
	s.logger.Info(constants.LogReviewCreated, "id", review.ID, "user_id", userID, "product_id", review.ProductID, "rating", review.Rating, "status", "COMPLETED")
	s.auditSvc.Record(userID, username, "CREATE", "review", uint64(review.ID), "create review for order_item="+util.UintString(req.OrderItemID))
	return s.ToVO(review), nil
}

// Reply 商家回复评价。
func (s *ReviewService) Reply(id uint, req dto.ReviewReplyRequest) (*dto.ReviewVO, error) {
	review, err := s.repo.FindByID(id)
	if err != nil {
		return nil, util.NewAppError(constants.CodeReviewNotFound, "review not found, id="+util.UintString(id))
	}
	if err := s.repo.UpdateReply(id, req.Reply); err != nil {
		return nil, err
	}
	now := time.Now()
	review.Reply = req.Reply
	review.RepliedAt = &now
	s.logger.Info(constants.LogReviewReplied, "id", review.ID, "reply", req.Reply)
	s.auditSvc.Record(0, "admin", "UPDATE", "review", uint64(review.ID), "reply review")
	return s.ToVO(review), nil
}

func (s *ReviewService) List(q dto.ReviewQuery) ([]dto.ReviewVO, int64, error) {
	list, total, err := s.repo.List(q)
	if err != nil {
		return nil, 0, err
	}
	vos := make([]dto.ReviewVO, 0, len(list))
	for i := range list {
		vos = append(vos, *s.ToVO(&list[i]))
	}
	return vos, total, nil
}

func (s *ReviewService) ToVO(r *model.Review) *dto.ReviewVO {
	username := ""
	if r.User != nil {
		username = r.User.Username
	}
	productName := ""
	if r.Product != nil {
		productName = r.Product.Name
	}
	return &dto.ReviewVO{
		ID:          r.ID,
		UserID:      r.UserID,
		Username:    username,
		OrderID:     r.OrderID,
		OrderItemID: r.OrderItemID,
		ProductID:   r.ProductID,
		ProductName: productName,
		Rating:      r.Rating,
		Content:     r.Content,
		Images:      util.UnmarshalStrings(r.Images),
		Reply:       r.Reply,
		RepliedAt:   util.FormatTimePtr(r.RepliedAt),
		CreatedAt:   util.FormatTime(r.CreatedAt),
	}
}
