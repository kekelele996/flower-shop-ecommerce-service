# Bug 复现说明

## Bug 是什么

购物车汇总金额与数量计算错误、加入购物车数量可超过 99、金额文案丢失小数位、OFF_SALE 被排除出合法商品状态。

## 如何触发

```bash
cd backend
go test ./internal/service/ -run 'TestCartSummary_TotalAndOrder|TestCartAdd_CapsAt99|TestPriceText_Decimals|TestValidProductStatuses_OffSale' -count=1
```

## 错误信息

```
--- FAIL: TestCartSummary_TotalAndOrder
    total price = 35, want 80
--- FAIL: TestCartAdd_CapsAt99
    quantity = 120, want 99
--- FAIL: TestPriceText_Decimals
    PriceText(129.5) = "¥130", want ¥129.50
```
