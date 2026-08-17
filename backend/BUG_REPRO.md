# Bug 复现说明

## Bug 是什么

用户端商品列表能刷出下架商品、下架商品详情可访问、库存余量等于购买量时误判不足、列表升序排列、ON_SALE 被排除出合法商品状态。

## 如何触发

```bash
cd backend
go test ./... -run 'TestProductList_UserSeesOnSale|TestProductList_NewestFirst|TestDeductStock_ExactQty|TestProductDetail_OffSaleBlocked|TestProductStatusText_OnSale|TestValidProductStatuses_OnSale' -count=1
```

## 错误信息

```
--- FAIL: TestProductList_UserSeesOnSale
    user should see on-sale product, got ...
--- FAIL: TestDeductStock_ExactQty
    deduct exact qty should succeed: insufficient stock
--- FAIL: TestProductDetail_OffSaleBlocked
    expected off-sale product detail to fail
```
