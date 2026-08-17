# Bug 复现说明

## Bug 是什么

订单详情/按单号查询查不到明细，物流轨迹解析为空。

## 如何触发

```bash
cd backend
go test ./internal/service/ -run 'TestOrderDetail_HasItems|TestFindByOrderNo_HasItems|TestLogistics_HasEvents'
```

## 错误信息

```
--- FAIL: TestOrderDetail_HasItems
    detail items = 0, want 1
--- FAIL: TestFindByOrderNo_HasItems
    order items = 0, want 1
--- FAIL: TestLogistics_HasEvents
    expected logistics events, got none
```
