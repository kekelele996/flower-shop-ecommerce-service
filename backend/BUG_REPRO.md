# Bug 复现说明

## Bug 是什么

订单支付/取消状态机错乱：待付款订单不能支付、已发货订单可被取消、支付后销量未增加、状态更新未落库、状态文案错误。

## 如何触发

```bash
cd backend
go test -run 'TestPay_UpdatesStatus|TestPay_IncrementsSales|TestCancel_ShippedBlocked|TestOrderStatusText_PendingPayment|TestValidOrderStatuses_Shipped' ./... -count=1
```

## 错误信息

```
--- FAIL: TestPay_UpdatesStatus
    status = PENDING_PAYMENT, want PENDING_SHIPMENT
--- FAIL: TestPay_IncrementsSales
    sales = 0, want 2
--- FAIL: TestCancel_ShippedBlocked
    expected cancel of shipped order to be blocked
```
