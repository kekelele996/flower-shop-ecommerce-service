package util

import (
	"testing"
	"time"
)

func TestOrderStatusText(t *testing.T) {
	tests := []struct {
		name   string
		status string
		want   string
	}{
		{"pending payment", "PENDING_PAYMENT", "待付款"},
		{"pending shipment", "PENDING_SHIPMENT", "待发货"},
		{"shipped", "SHIPPED", "已发货"},
		{"completed", "COMPLETED", "已完成"},
		{"cancelled", "CANCELLED", "已取消"},
		{"unknown", "WHATEVER", "未知状态"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := OrderStatusText(tt.status); got != tt.want {
				t.Errorf("OrderStatusText(%q) = %q, want %q", tt.status, got, tt.want)
			}
		})
	}
}

func TestMoneyToCents(t *testing.T) {
	tests := []struct {
		amount float64
		want   int64
	}{
		{0, 0},
		{12.34, 1234},
		{99.99, 9999},
	}
	for _, tt := range tests {
		if got := MoneyToCents(tt.amount); got != tt.want {
			t.Errorf("MoneyToCents(%v) = %d, want %d", tt.amount, got, tt.want)
		}
	}
}

func TestFormatTimePtr(t *testing.T) {
	now := time.Now()
	if got := FormatTimePtr(&now); got == "" {
		t.Error("FormatTimePtr(non-nil) should not be empty")
	}
	if got := FormatTimePtr(nil); got != "" {
		t.Errorf("FormatTimePtr(nil) = %q, want empty", got)
	}
}

func TestJoinImages(t *testing.T) {
	if got := JoinImages(`["a.jpg","b.jpg"]`); got != "a.jpg, b.jpg" {
		t.Errorf("JoinImages = %q", got)
	}
}
