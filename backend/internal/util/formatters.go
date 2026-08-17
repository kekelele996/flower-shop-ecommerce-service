package util

import (
	"fmt"
	"strings"
	"time"
)

// FormatTime 格式化时间为 "2006-01-02 15:04:05"。
func FormatTime(t time.Time) string {
	if t.IsZero() {
		return ""
	}
	return t.Format("2006-01-02 15:04:05")
}

// OrderStatusText 订单状态 → 中文文本（与前端状态徽标、日志模板、错误码共同维护状态机）。
func OrderStatusText(status string) string {
	switch status {
	case "PENDING_PAYMENT":
		return "待付款"
	case "PENDING_SHIPMENT":
		return "待发货"
	case "SHIPPED":
		return "已发货"
	case "COMPLETED":
		return "已完成"
	case "CANCELLED":
		return "已取消"
	default:
		return "未知状态"
	}
}

// ProductStatusText 商品状态 → 中文文本。
func ProductStatusText(status string) string {
	switch status {
	case "ON_SALE":
		return "在售"
	case "OFF_SALE":
		return "下架"
	default:
		return "未知"
	}
}

// CouponTypeText 优惠券类型 → 中文文本。
func CouponTypeText(t string) string {
	switch t {
	case "FULL_REDUCTION":
		return "满减券"
	case "DISCOUNT":
		return "折扣券"
	default:
		return "未知"
	}
}

// LogisticsStatusText 物流状态 → 中文文本。
func LogisticsStatusText(status string) string {
	switch status {
	case "PENDING":
		return "待揽收"
	case "PICKED_UP":
		return "已揽收"
	case "IN_TRANSIT":
		return "运输中"
	case "OUT_FOR_DELIVERY":
		return "派送中"
	case "DELIVERED":
		return "已签收"
	default:
		return "未知"
	}
}

// JoinImages 图片 JSON 数组字符串 → 展示文本。
func JoinImages(imagesJSON string) string {
	imgs := UnmarshalStrings(imagesJSON)
	return strings.Join(imgs, ", ")
}

// PriceText 金额格式化（保留两位小数）。
func PriceText(amount float64) string {
	return fmt.Sprintf("¥%.2f", amount)
}

// MoneyToCents 元转分，便于金额计算避免浮点误差。
func MoneyToCents(amount float64) int64 {
	return int64(amount*100 + 0.5)
}

// FormatTimePtr 格式化时间指针。
func FormatTimePtr(t *time.Time) string {
	if t == nil {
		return ""
	}
	return FormatTime(*t)
}
