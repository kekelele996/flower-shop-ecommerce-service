package util

import (
	"fmt"
	"time"
)

// GenerateOrderNo 生成订单号：yyyyMMddHHmmss + 6 位随机数字。
func GenerateOrderNo() string {
	return fmt.Sprintf("%s%06d", time.Now().Format("20060102150405"), time.Now().UnixNano()%1000000)
}

// GeneratePayNo 生成支付流水号。
func GeneratePayNo() string {
	return "PAY" + time.Now().Format("20060102150405") + fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)
}

// GenerateTrackingNo 生成模拟物流单号。
func GenerateTrackingNo() string {
	return "SF" + time.Now().Format("20060102150405") + fmt.Sprintf("%04d", time.Now().UnixNano()%10000)
}
