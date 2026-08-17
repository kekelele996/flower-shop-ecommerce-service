package service

import (
	"errors"
	"log/slog"
	"time"

	"github.com/flowershop/backend/internal/constants"
	"github.com/flowershop/backend/internal/dto"
	"github.com/flowershop/backend/internal/model"
	"github.com/flowershop/backend/internal/repository"
	"github.com/flowershop/backend/internal/util"
	"gorm.io/gorm"
)

// OrderService 订单服务：下单/支付/取消/发货/完成/物流。
type OrderService struct {
	db          *gorm.DB
	orderRepo   *repository.OrderRepository
	itemRepo    *repository.OrderItemRepository
	paymentRepo *repository.PaymentRepository
	logisRepo   *repository.LogisticsRepository
	cartRepo    *repository.CartRepository
	prodRepo    *repository.ProductRepository
	couponRepo  *repository.CouponRepository
	auditSvc    *AuditService
	logger      *slog.Logger
}

func NewOrderService(db *gorm.DB, orderRepo *repository.OrderRepository, itemRepo *repository.OrderItemRepository, paymentRepo *repository.PaymentRepository, logisRepo *repository.LogisticsRepository, cartRepo *repository.CartRepository, prodRepo *repository.ProductRepository, couponRepo *repository.CouponRepository, auditSvc *AuditService, logger *slog.Logger) *OrderService {
	return &OrderService{db: db, orderRepo: orderRepo, itemRepo: itemRepo, paymentRepo: paymentRepo, logisRepo: logisRepo, cartRepo: cartRepo, prodRepo: prodRepo, couponRepo: couponRepo, auditSvc: auditSvc, logger: logger}
}

// Checkout 结算下单：事务内锁库存、应用优惠券、生成订单。
func (s *OrderService) Checkout(userID uint, req dto.CheckoutRequest) (*dto.OrderVO, error) {
	var order *model.Order
	err := s.db.Transaction(func(tx *gorm.DB) error {
		items, err := s.cartRepo.FindByUserID(userID)
		if err != nil {
			return err
		}
		selected := make([]model.CartItem, 0, len(items))
		idSet := map[uint]bool{}
		for _, id := range req.CartItemIDs {
			idSet[id] = true
		}
		for _, it := range items {
			if idSet[it.ID] && it.Selected {
				selected = append(selected, it)
			}
		}
		if len(selected) == 0 {
			return util.NewAppError(constants.CodeBadRequest, "no selected cart items for user="+util.UintString(userID))
		}

		totalAmount := 0.0
		orderItems := make([]model.OrderItem, 0, len(selected))
		productIDs := make([]uint, 0, len(selected))
		qtyMap := map[uint]int{}
		for _, it := range selected {
			// 扣库存（行锁防超卖）
			if err := s.prodRepo.DeductStock(tx, it.ProductID, it.Quantity); err != nil {
				return err
			}
			productIDs = append(productIDs, it.ProductID)
			qtyMap[it.ProductID] = it.Quantity
			totalAmount += float64(it.Quantity) * it.Product.Price
		}
		products, err := s.prodRepo.FindByIDs(productIDs)
		if err != nil {
			return err
		}
		prodMap := map[uint]*model.Product{}
		for i := range products {
			prodMap[products[i].ID] = &products[i]
		}

		discount := 0.0
		couponID := req.CouponID
		if couponID > 0 {
			coupon, err := s.couponRepo.FindByIDForUpdate(tx, couponID)
			if err != nil {
				return util.NewAppError(constants.CodeCouponNotFound, "coupon not found, id="+util.UintString(couponID))
			}
			if coupon.UserID != userID || coupon.Status != constants.CouponStatusUnused {
				return util.NewAppError(constants.CodeCouponInvalid, "coupon="+coupon.Name+" is not usable for user="+util.UintString(userID))
			}
			now := time.Now()
			if now.Before(coupon.ValidFrom) || now.After(coupon.ValidTo) {
				return util.NewAppError(constants.CodeCouponInvalid, "coupon="+coupon.Name+" is expired")
			}
			if totalAmount < coupon.Threshold {
				return util.NewAppError(constants.CodeCouponInvalid, constants.MsgCouponNotApplicable)
			}
			if coupon.Type == constants.CouponTypeFullReduction {
				discount = coupon.Amount
			} else if coupon.Type == constants.CouponTypeDiscount {
				discount = totalAmount * (1 - coupon.DiscountRate)
			}
			if discount > totalAmount {
				discount = totalAmount
			}
		}

		shippingFee := 0.0
		for _, pid := range productIDs {
			if p := prodMap[pid]; p != nil && !p.FreeShipping {
				shippingFee = 8
				break
			}
		}

		order = &model.Order{
			OrderNo:        util.GenerateOrderNo(),
			UserID:         userID,
			Status:         constants.OrderStatusPendingPayment,
			TotalAmount:    totalAmount,
			DiscountAmount: discount,
			ShippingFee:    shippingFee,
			PayAmount:      totalAmount - discount + shippingFee,
			CouponID:       couponID,
			ReceiverName:   req.ReceiverName,
			ReceiverPhone:  req.ReceiverPhone,
			ReceiverAddr:   req.ReceiverAddr,
			Remark:         req.Remark,
		}
		if err := s.orderRepo.Create(tx, order); err != nil {
			return err
		}
		for _, it := range selected {
			p := prodMap[it.ProductID]
			orderItems = append(orderItems, model.OrderItem{
				OrderID:      order.ID,
				ProductID:    it.ProductID,
				ProductName:  p.Name,
				ProductImage: p.CoverImage,
				Price:        p.Price,
				Quantity:     it.Quantity,
				TotalPrice:   float64(it.Quantity) * p.Price,
			})
		}
		if err := s.itemRepo.CreateBatch(tx, orderItems); err != nil {
			return err
		}
		if couponID > 0 {
			if err := s.couponRepo.MarkUsed(tx, couponID, order.ID); err != nil {
				return err
			}
		}
		if err := s.cartRepo.DeleteByIDs(req.CartItemIDs); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	s.logger.Info(constants.LogOrderCreated, "id", order.ID, "order_no", order.OrderNo, "user_id", userID, "amount", order.PayAmount, "status", order.Status)
	s.auditSvc.Record(userID, "", "CREATE", "order", uint64(order.ID), "checkout order "+order.OrderNo)
	return s.ToVO(order), nil
}

// Pay 模拟支付宝沙箱支付。
func (s *OrderService) Pay(userID, orderID uint, req dto.PayRequest) (*dto.OrderVO, error) {
	channel := req.Channel
	if channel == "" {
		channel = constants.PaymentChannelAlipaySandbox
	}
	var order *model.Order
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		order, err = s.orderRepo.FindByIDForUpdate(tx, orderID)
		if err != nil {
			return err
		}
		if order.UserID != userID {
			return util.NewAppError(constants.CodeForbidden, "order="+order.OrderNo+" does not belong to user="+util.UintString(userID))
		}
		if order.Status == constants.OrderStatusPendingPayment {
			return util.NewAppError(constants.CodeOrderStatusNotAllow, constants.MsgOrderStatusError)
		}
		pay := &model.Payment{
			OrderID: order.ID,
			OrderNo: order.OrderNo,
			PayNo:   util.GeneratePayNo(),
			UserID:  userID,
			Channel: channel,
			Amount:  order.PayAmount,
			Status:  constants.PaymentStatusSuccess,
		}
		if err := s.paymentRepo.Create(tx, pay); err != nil {
			return err
		}
		if err := s.paymentRepo.MarkSuccess(tx, pay.PayNo, "SANDBOX_TRADE_"+util.UintString(pay.ID)); err != nil {
			return err
		}
		now := time.Now()
		if err := s.orderRepo.UpdateStatus(tx, order.ID, constants.OrderStatusPendingShipment, map[string]interface{}{"paid_at": &now}); err != nil {
			return err
		}
		order.Status = constants.OrderStatusPendingShipment
		order.PaidAt = &now
		s.logger.Info(constants.LogOrderPaid, "id", order.ID, "order_no", order.OrderNo, "pay_no", pay.PayNo, "amount", pay.Amount, "status", order.Status)
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.auditSvc.Record(userID, "", "UPDATE", "order", uint64(order.ID), "pay order "+order.OrderNo)
	return s.ToVO(order), nil
}

// Cancel 取消未付款订单（回补库存）。
func (s *OrderService) Cancel(userID, orderID uint, reason string) (*dto.OrderVO, error) {
	var order *model.Order
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		order, err = s.orderRepo.FindByIDForUpdate(tx, orderID)
		if err != nil {
			return err
		}
		if order.UserID != userID {
			return util.NewAppError(constants.CodeForbidden, "order="+order.OrderNo+" does not belong to user="+util.UintString(userID))
		}
		if order.Status == constants.OrderStatusPendingPayment {
			return util.NewAppError(constants.CodeOrderStatusNotAllow, constants.MsgOrderStatusError)
		}
		now := time.Now()
		if err := s.orderRepo.UpdateStatus(tx, order.ID, constants.OrderStatusCancelled, map[string]interface{}{"cancelled_at": &now, "cancelled_by": "user:" + util.UintString(userID)}); err != nil {
			return err
		}
		for _, item := range order.Items {
			if err := s.prodRepo.RestoreStock(tx, item.ProductID, item.Quantity); err != nil {
				return err
			}
		}
		// 退回优惠券
		if order.CouponID > 0 {
			if err := s.couponRepo.MarkUnused(tx, order.CouponID); err != nil {
				return err
			}
		}
		order.Status = constants.OrderStatusCancelled
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.logger.Info(constants.LogOrderCancelled, "id", order.ID, "order_no", order.OrderNo, "user_id", userID, "status", order.Status)
	s.auditSvc.Record(userID, "", "UPDATE", "order", uint64(order.ID), "cancel order "+order.OrderNo+" reason="+reason)
	return s.ToVO(order), nil
}

// Ship 管理员发货：创建物流单并生成轨迹。
func (s *OrderService) Ship(orderID uint, carrier string) (*dto.OrderVO, error) {
	if carrier == "" {
		carrier = "SF_EXPRESS"
	}
	var order *model.Order
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		order, err = s.orderRepo.FindByIDForUpdate(tx, orderID)
		if err != nil {
			return err
		}
		if order.Status != constants.OrderStatusPendingShipment {
			return util.NewAppError(constants.CodeOrderStatusNotAllow, constants.MsgOrderStatusError)
		}
		trackingNo := util.GenerateTrackingNo()
		now := time.Now()
		events := []model.LogisticsEvent{
			{Time: now, Status: constants.LogisticsStatusPickedUp, Desc: "商家已发货，快件已被顺丰速运揽收"},
		}
		logis := &model.Logistics{
			OrderID:    order.ID,
			OrderNo:    order.OrderNo,
			TrackingNo: trackingNo,
			Carrier:    carrier,
			Status:     constants.LogisticsStatusPickedUp,
			Events:     util.MarshalJSON(events),
		}
		if err := s.logisRepo.Create(tx, logis); err != nil {
			return err
		}
		if err := s.orderRepo.UpdateStatus(tx, order.ID, constants.OrderStatusShipped, map[string]interface{}{"shipped_at": &now}); err != nil {
			return err
		}
		order.Status = constants.OrderStatusShipped
		order.ShippedAt = &now
		s.logger.Info(constants.LogOrderShipped, "id", order.ID, "order_no", order.OrderNo, "tracking_no", trackingNo, "status", order.Status)
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.auditSvc.Record(0, "admin", "UPDATE", "order", uint64(order.ID), "ship order "+order.OrderNo)
	return s.ToVO(order), nil
}

// Complete 确认收货。
func (s *OrderService) Complete(userID, orderID uint) (*dto.OrderVO, error) {
	var order *model.Order
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		order, err = s.orderRepo.FindByIDForUpdate(tx, orderID)
		if err != nil {
			return err
		}
		if order.UserID != userID {
			return util.NewAppError(constants.CodeForbidden, "order="+order.OrderNo+" does not belong to user="+util.UintString(userID))
		}
		if order.Status != constants.OrderStatusShipped {
			return util.NewAppError(constants.CodeOrderStatusNotAllow, constants.MsgOrderStatusError)
		}
		now := time.Now()
		if err := s.orderRepo.UpdateStatus(tx, order.ID, constants.OrderStatusCompleted, map[string]interface{}{"completed_at": &now}); err != nil {
			return err
		}
		order.Status = constants.OrderStatusCompleted
		order.CompletedAt = &now
		// 物流更新为已签收
		if logis, err := s.logisRepo.FindByOrderID(order.ID); err == nil {
			logis.Status = constants.LogisticsStatusDelivered
			evts := append([]model.LogisticsEvent{{Time: now, Status: constants.LogisticsStatusDelivered, Desc: "包裹已签收，感谢您的购买"}}, util.UnmarshalLogisticsEvents(logis.Events)...)
			logis.Events = util.MarshalJSON(evts)
			if err := s.logisRepo.Update(logis); err != nil {
				return err
			}
		}
		s.logger.Info(constants.LogOrderCompleted, "id", order.ID, "order_no", order.OrderNo, "status", order.Status)
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.auditSvc.Record(userID, "", "UPDATE", "order", uint64(order.ID), "complete order "+order.OrderNo)
	return s.ToVO(order), nil
}

// AdminComplete 管理员确认收货（跳过归属校验）。
func (s *OrderService) AdminComplete(operatorID, orderID uint) (*dto.OrderVO, error) {
	var order *model.Order
	err := s.db.Transaction(func(tx *gorm.DB) error {
		var err error
		order, err = s.orderRepo.FindByIDForUpdate(tx, orderID)
		if err != nil {
			return err
		}
		if order.Status != constants.OrderStatusShipped {
			return util.NewAppError(constants.CodeOrderStatusNotAllow, constants.MsgOrderStatusError)
		}
		now := time.Now()
		if err := s.orderRepo.UpdateStatus(tx, order.ID, constants.OrderStatusCompleted, map[string]interface{}{"completed_at": &now}); err != nil {
			return err
		}
		order.Status = constants.OrderStatusCompleted
		order.CompletedAt = &now
		return nil
	})
	if err != nil {
		return nil, err
	}
	s.auditSvc.Record(operatorID, "admin", "UPDATE", "order", uint64(order.ID), "admin complete order "+order.OrderNo)
	return s.ToVO(order), nil
}

// List 订单列表（用户/管理端）。
func (s *OrderService) List(q dto.OrderQuery, userID uint) ([]dto.OrderVO, int64, error) {
	list, total, err := s.orderRepo.List(q, userID)
	if err != nil {
		return nil, 0, err
	}
	vos := make([]dto.OrderVO, 0, len(list))
	for i := range list {
		vos = append(vos, *s.ToVO(&list[i]))
	}
	return vos, total, nil
}

// Detail 订单详情（含物流）。
func (s *OrderService) Detail(userID, orderID uint, admin bool) (*dto.OrderVO, error) {
	order, err := s.orderRepo.FindByID(s.db, orderID)
	if err != nil {
		return nil, util.NewAppError(constants.CodeOrderNotFound, "order not found, id="+util.UintString(orderID))
	}
	if !admin && order.UserID != userID {
		return nil, util.NewAppError(constants.CodeForbidden, "order="+order.OrderNo+" does not belong to user="+util.UintString(userID))
	}
	return s.ToVO(order), nil
}

// Logistics 查询物流轨迹。
func (s *OrderService) Logistics(userID, orderID uint, admin bool) (*dto.LogisticsVO, error) {
	order, err := s.orderRepo.FindByID(s.db, orderID)
	if err != nil {
		return nil, util.NewAppError(constants.CodeOrderNotFound, "order not found, id="+util.UintString(orderID))
	}
	if !admin && order.UserID != userID {
		return nil, util.NewAppError(constants.CodeForbidden, "order="+order.OrderNo+" does not belong to user="+util.UintString(userID))
	}
	logis, err := s.logisRepo.FindByOrderID(orderID)
	if errors.Is(err, util.ErrNotFound) {
		return nil, util.NewAppError(constants.CodeNotFound, "order="+order.OrderNo+" has no logistics yet")
	}
	if err != nil {
		return nil, err
	}
	events := util.UnmarshalLogisticsEvents(logis.Events)
	voEvents := make([]dto.LogisticsEventVO, 0, len(events))
	for _, ev := range events {
		voEvents = append(voEvents, dto.LogisticsEventVO{Time: util.FormatTime(ev.Time), Status: ev.Status, Desc: ev.Desc})
	}
	return &dto.LogisticsVO{
		TrackingNo: logis.TrackingNo,
		Carrier:    logis.Carrier,
		Status:     logis.Status,
		StatusText: util.LogisticsStatusText(logis.Status),
		Events:     voEvents,
	}, nil
}

// ToVO 订单模型转视图对象。
func (s *OrderService) ToVO(order *model.Order) *dto.OrderVO {
	vo := &dto.OrderVO{
		ID:             order.ID,
		OrderNo:        order.OrderNo,
		UserID:         order.UserID,
		Status:         order.Status,
		StatusText:     util.OrderStatusText(order.Status),
		TotalAmount:    order.TotalAmount,
		DiscountAmount: order.DiscountAmount,
		ShippingFee:    order.ShippingFee,
		PayAmount:      order.PayAmount,
		CouponID:       order.CouponID,
		ReceiverName:   order.ReceiverName,
		ReceiverPhone:  order.ReceiverPhone,
		ReceiverAddr:   order.ReceiverAddr,
		Remark:         order.Remark,
		PaidAt:         util.FormatTimePtr(order.PaidAt),
		ShippedAt:      util.FormatTimePtr(order.ShippedAt),
		CompletedAt:    util.FormatTimePtr(order.CompletedAt),
		CancelledAt:    util.FormatTimePtr(order.CancelledAt),
		CreatedAt:      util.FormatTime(order.CreatedAt),
	}
	for _, item := range order.Items {
		vo.Items = append(vo.Items, dto.OrderItemVO{
			ID:           item.ID,
			ProductID:    item.ProductID,
			ProductName:  item.ProductName,
			ProductImage: item.ProductImage,
			Price:        item.Price,
			Quantity:     item.Quantity,
			TotalPrice:   item.TotalPrice,
			Reviewed:     item.Reviewed,
		})
	}
	return vo
}

// toProductVO 复用 ProductService 转换逻辑（避免循环依赖，轻量复制）。
func toProductVO(p *model.Product) *dto.ProductVO {
	return &dto.ProductVO{
		ID:            p.ID,
		CategoryID:    p.CategoryID,
		Name:          p.Name,
		SubTitle:      p.SubTitle,
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
	}
}
