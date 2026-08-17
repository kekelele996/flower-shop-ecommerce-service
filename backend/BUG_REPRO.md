# Bug 复现说明

## Bug 是什么

销售统计漏掉已完成订单、订单列表能看到其他用户的订单、列表明细未加载。

## 如何触发

```bash
cd backend
go test ./internal/service/ -run 'TestSumSales_IncludesCompleted|TestOrderList_UserScoped|TestOrderList_HasItems'
```

## 错误信息

```
--- FAIL: TestSumSales_IncludesCompleted
    sum sales = 150, want 180
--- FAIL: TestOrderList_UserScoped
    user should only see own orders: ...
--- FAIL: TestOrderList_HasItems
    list items wrong: ...
```
