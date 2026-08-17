# 花语商城（flowershop）· 花卉 B2C 电商平台

面向花卉爱好者的 B2C 电商平台：支持花卉、绿植、花盆及园艺工具的浏览、搜索、加购、下单、支付宝沙箱模拟支付、物流跟踪、评价、优惠券与商家后台管理。

## 快速启动（Docker Compose 一键部署）

```bash
cd 项目根目录
docker compose up -d --build
```

启动后访问：

| 入口 | 地址 |
| --- | --- |
| 前端 | http://localhost:8101 |
| 后端 API | http://localhost:3101/api/v1 |
| 健康检查 | http://localhost:3101/healthz |
| MinIO 控制台 | http://localhost:47025 （minioadmin / minioadmin） |

测试账号：

- 普通用户：`user` / `user123`
- 管理员：`admin` / `admin123`

## 本地开发（备选）

```bash
# 后端
cd backend
export GOPROXY=https://goproxy.cn,direct
go mod tidy
go run ./cmd/server
# 构建：go build ./...

# 前端
cd frontend
npm install --registry=https://registry.npmmirror.com
npm run dev
```

## 项目主要功能

1. 商品多级分类浏览（鲜切花/盆栽/多肉/花盆/园艺工具/营养土肥），按价格、销量、评分排序，多图与图文详情。
2. 关键词搜索 + 价格区间、发货地、是否包邮筛选；首页基于浏览记录推荐（Redis）。
3. 购物车实时计价、数量修改、勾选结算；订单状态机（待付款→待发货→已发货→已完成/已取消），支持取消未付款订单（回补库存）。
4. 支付宝沙箱模拟支付：支付成功自动流转状态并生成支付记录。
5. 物流跟踪：发货生成运单与轨迹，支持查询物流轨迹。
6. 用户评价（1-5 星 + 图文），商家可回复。
7. 满减/折扣优惠券，结算自动抵扣，取消订单自动退回。
8. 商家后台：商品上下架、库存价格管理、订单发货、销售数据统计、审计日志。

## 技术栈

| 层 | 技术 |
| --- | --- |
| 前端 | React 18 + TypeScript + Ant Design，构建工具 Vite |
| 后端 | Go 1.22 + Gin + GORM |
| 数据库 | PostgreSQL 16 |
| 缓存 | Redis 7 |
| 对象存储 | MinIO |
| 认证 | JWT（HS256）+ RBAC |
| 部署 | Docker Compose + Nginx |

## 项目目录结构

```text
.
├── backend/                 # Go 后端（cmd + internal + pkg 标准布局）
│   ├── cmd/server/main.go   # 装配入口
│   ├── internal/
│   │   ├── config/          # 环境变量配置
│   │   ├── database/        # 连接、迁移、种子数据
│   │   ├── constants/       # 枚举 / 错误码 / 日志模板 / 消息文案
│   │   ├── model/           # 每个实体一个文件
│   │   ├── dto/             # 每个实体一个 DTO 文件
│   │   ├── repository/      # 每个实体一个仓储文件
│   │   ├── service/         # 每个实体一个服务文件
│   │   ├── handler/         # 每个实体一个处理器文件
│   │   ├── router/          # 路由聚合 + 每实体路由文件
│   │   ├── middleware/      # auth / rbac / audit / request_id / error_handler / logger / cors
│   │   └── util/            # jwt / logger / formatters / response / app_error 等
│   ├── migrations/          # 初始化建表 SQL
│   ├── api/                 # OpenAPI 摘要
│   └── Dockerfile
├── frontend/                # React 18 + TS + Ant Design
│   ├── src/
│   │   ├── api/             # 每个实体一个 API 文件
│   │   ├── components/      # StatusBadge / EmptyState / DataTable / ConfirmDialog / StarRating / ProductCard / PriceTag / AuthGuard
│   │   ├── hooks/           # useAuth / usePagination / useRequest
│   │   ├── stores/          # user / cart / product / order / review / coupon store
│   │   ├── pages/           # 首页、详情、登录、注册、购物车、结算、订单、个人中心、优惠券、管理后台
│   │   ├── constants/       # 与后端对应枚举
│   │   └── utils/           # request（拦截器）/ format
│   ├── Dockerfile
│   └── nginx.conf
├── database/                # 数据库脚本
├── docker-compose.yml
├── .env / .env.example
└── README.md
```

**文件结构强制清单（严禁合并职责到单一文件）**：

- 后端每个实体独立拆分 `model/<entity>.go`、`dto/<entity>_dto.go`、`repository/<entity>_repository.go`、`service/<entity>_service.go`、`handler/<entity>_handler.go`、`router/<entity>.go`、`constants/*`，严禁把多实体逻辑写进同一文件。
- 前端按模块拆分 `api/<entity>.ts`、`stores/<entity>Store.ts`、`pages/<module>/`、`components/`、`hooks/`，禁止把所有页面写进 App.tsx。
- 依赖方向单向：handler → service → repository → model，构造器注入，禁止全局 init 连接数据库。

## 环境变量说明

| 变量 | 默认值 | 说明 |
| --- | --- | --- |
| COMPOSE_PROJECT_NAME | flowershop | 容器名/网络前缀 |
| DB_NAME / DB_USER / DB_PASSWORD | flowershop_db / flowershop_user / flowershop_pwd | PostgreSQL 库名/用户/密码 |
| JWT_SECRET | change_me_to_a_long_random_string | JWT 签名密钥（生产必改） |
| FRONTEND_PORT / BACKEND_PORT | 8101 / 3101 | 前端/后端宿主机端口 |
| DB_PORT / REDIS_PORT | 5601 / 6501 | 数据库/Redis 宿主机端口 |
| MINIO_PORT / MINIO_CONSOLE_PORT | 47024 / 47025 | MinIO API/控制台端口 |
| REDIS_ADDR / REDIS_PASSWORD | redis:6379 / 空 | Redis 连接 |
| MINIO_ENDPOINT / MINIO_ROOT_USER / MINIO_ROOT_PASSWORD | minio:9000 / minioadmin / minioadmin | MinIO 连接 |

## API 清单（统一前缀 /api/v1，统一响应 `{"code":0,"message":"ok","data":...}`）

### 认证与用户
| 方法 | 路径 | 说明 | 权限 |
| --- | --- | --- | --- |
| POST | /auth/register | 注册 | 公开 |
| POST | /auth/login | 登录（返回 JWT） | 公开 |
| GET | /auth/profile | 当前用户 | 登录 |
| PUT | /auth/profile | 更新资料 | 登录 |
| PUT | /auth/password | 修改密码 | 登录 |

### 分类与商品
| 方法 | 路径 | 说明 | 权限 |
| --- | --- | --- | --- |
| GET | /categories | 分类树 | 公开 |
| GET | /products | 列表/搜索（keyword/category_id/min_price/max_price/free_shipping/shipping_from/sort） | 公开 |
| GET | /products/recommendations | 首页推荐（复用 productRepository.Search / RecommendByCategory） | 公开 |
| GET | /products/:id | 商品详情 | 公开 |
| POST | /products/:id/view | 记录浏览（Redis） | 公开 |

### 购物车 / 订单 / 评价 / 优惠券（登录）
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET/POST | /cart | 购物车汇总 / 添加 |
| PUT/DELETE | /cart/:id | 修改数量/勾选 / 删除 |
| POST | /orders | 结算下单（事务 + SELECT FOR UPDATE 扣库存） |
| GET | /orders | 我的订单列表 |
| GET | /orders/:id | 订单详情 |
| POST | /orders/:id/pay | 模拟支付宝沙箱支付 |
| POST | /orders/:id/cancel | 取消未付款订单（回补库存、退券） |
| POST | /orders/:id/complete | 确认收货 |
| GET | /orders/:id/logistics | 物流轨迹 |
| POST | /reviews | 提交评价（完成后可评） |
| GET | /reviews | 评价列表（product_id/order_id/user_id） |
| GET | /coupons | 我的优惠券 |
| GET | /coupons/available | 可用优惠券 |
| POST | /coupons/:id/claim | 领取优惠券 |
| POST | /uploads | 图片上传（kind=product/review/avatar） |

### 管理员（ADMIN）
| 方法 | 路径 | 说明 |
| --- | --- | --- |
| GET | /admin/products | 商品列表（含下架） |
| POST | /admin/products | 创建商品 |
| PUT | /admin/products/:id | 更新商品 |
| PUT | /admin/products/:id/status | 上下架 |
| GET | /admin/orders | 全部订单 |
| POST | /admin/orders/:id/ship | 发货（生成物流单） |
| POST | /admin/orders/:id/complete | 代确认收货 |
| POST | /admin/orders/:id/pay | 代支付 |
| GET | /admin/orders/:id/logistics | 查看物流 |
| GET | /admin/reviews | 全部评价 |
| POST | /admin/reviews/:id/reply | 回复评价 |
| POST | /admin/coupons | 创建优惠券模板 |
| GET | /admin/audit-logs | 审计日志 |
| GET | /admin/stats | 销售统计 |

复用关系：`/cart`（购物车汇总）与 `/orders`（结算）复用 `cartRepository.FindByUserID`；`/products` 与 `/products/recommendations` 复用 `productRepository.Search / RecommendByCategory`；`/orders` 与 `/admin/orders` 复用 `orderRepository.List`。

### curl 调用示例

```bash
# 登录获取 token
TOKEN=$(curl -s -X POST http://localhost:3101/api/v1/auth/login \
  -H 'Content-Type: application/json' \
  -d '{"username":"user","password":"user123"}' | python3 -c 'import sys,json;print(json.load(sys.stdin)["data"]["access_token"])')

# 带 JWT 的请求
curl -s http://localhost:3101/api/v1/cart -H "Authorization: Bearer $TOKEN"

# 下单
curl -s -X POST http://localhost:3101/api/v1/orders -H "Authorization: Bearer $TOKEN" \
  -H 'Content-Type: application/json' \
  -d '{"cart_item_ids":[1],"receiver_name":"张三","receiver_phone":"13800138000","receiver_addr":"上海市浦东新区花木路100号"}'
```

## Docker 部署说明

- 端口映射：前端 `8101:80`，后端 `3101:8080`，数据库 `5601:5432`，Redis `6501:6379`，MinIO `47024:9000` / `47025:9001`。
- 数据卷：`db_data`、`redis_data`、`minio_data` 命名卷持久化；不在任何中文路径绑定挂载，因此项目放在中文目录名下也能正常启动。
- 健康检查：db / redis / minio / backend / frontend 均配置 healthcheck，后端依赖数据库、缓存、对象存储就绪后才启动。
- 常见问题：
  - 端口被占用：修改 `.env` 中的 `*_PORT` 后 `docker compose up -d` 重建。
  - 后端启动失败：`docker compose logs backend` 查看日志。
  - 清空数据：`docker compose down -v`（会删除命名卷数据）。

## 枚举出现位置清单

### 1. 订单状态（PENDING_PAYMENT / PENDING_SHIPMENT / SHIPPED / COMPLETED / CANCELLED）

后端：
- `backend/internal/constants/enums.go`（定义 + ValidOrderStatuses）
- `backend/internal/model/order.go`（Status 字段默认值）
- `backend/internal/dto/order_dto.go`（OrderVO.Status/StatusText）
- `backend/internal/service/order_service.go`（状态机：Checkout/Pay/Cancel/Ship/Complete）
- `backend/internal/handler/order_handler.go`（状态相关错误处理）
- `backend/internal/repository/order_repository.go`（List/CountByStatus/SumSales 按状态过滤）
- `backend/internal/util/formatters.go`（OrderStatusText）
- `backend/internal/constants/error_codes.go`（CodeOrderStatusNotAllow）
- `backend/internal/constants/log_templates.go`（LogOrderCreated/Paid/Cancelled/Shipped/Completed）
- `backend/internal/constants/messages.go`（MsgOrderStatusError）

前端：
- `frontend/src/constants/enums.ts`（ORDER_STATUS / ORDER_STATUS_TEXT）
- `frontend/src/api/order.ts`（类型定义）
- `frontend/src/components/StatusBadge.tsx`（状态徽标颜色）
- `frontend/src/pages/Orders/index.tsx`（筛选 Tab + 按钮显隐）
- `frontend/src/pages/OrderDetail/index.tsx`（Steps + 操作按钮）
- `frontend/src/pages/admin/OrderManage.tsx`（筛选 + 发货/收货按钮）

### 2. 商品状态（ON_SALE / OFF_SALE）

后端：
- `backend/internal/constants/enums.go`（定义 + ValidProductStatuses）
- `backend/internal/model/product.go`（Status 默认值）
- `backend/internal/dto/product_dto.go`（Create/Update/StatusRequest oneof 校验）
- `backend/internal/service/product_service.go`（Detail 下架拦截、ChangeStatus）
- `backend/internal/repository/product_repository.go`（Search 按状态过滤）
- `backend/internal/util/formatters.go`（ProductStatusText）
- `backend/internal/constants/log_templates.go`（LogProductCreated/Updated/StatusChanged）
- `backend/internal/constants/error_codes.go`（CodeProductNotFound 下架提示复用）

前端：
- `frontend/src/constants/enums.ts`（PRODUCT_STATUS / PRODUCT_STATUS_TEXT）
- `frontend/src/api/product.ts`（类型）
- `frontend/src/components/StatusBadge.tsx`
- `frontend/src/pages/ProductDetail/index.tsx`（状态 Tag）
- `frontend/src/pages/admin/ProductManage.tsx`（上下架按钮/下拉）

### 3. 用户角色（USER / ADMIN）

后端：
- `backend/internal/constants/enums.go`
- `backend/internal/model/user.go`（Role 默认值）
- `backend/internal/dto/user_dto.go`（UserVO.Role）
- `backend/internal/middleware/rbac.go`（RBAC 校验）
- `backend/internal/middleware/auth.go`（JWT Claims.Role）
- `backend/internal/util/jwt.go`（Claims 定义）
- `backend/internal/service/user_service.go`（注册默认角色）
- `backend/internal/constants/log_templates.go`（LogUserRegistered/LoggedIn）

前端：
- `frontend/src/constants/enums.ts`（ROLE）
- `frontend/src/api/auth.ts`（UserVO.role）
- `frontend/src/stores/userStore.ts`（isAdmin）
- `frontend/src/components/AuthGuard.tsx`（路由守卫）
- `frontend/src/layouts/MainLayout.tsx`（后台入口显隐）

### 4. 优惠券类型/状态（FULL_REDUCTION / DISCOUNT；UNUSED / USED / EXPIRED）

后端：
- `backend/internal/constants/enums.go`
- `backend/internal/model/coupon.go`（Type/Status）
- `backend/internal/dto/coupon_dto.go`（oneof 校验）
- `backend/internal/service/coupon_service.go`（创建/领取）
- `backend/internal/service/order_service.go`（结算抵扣、MarkUsed/MarkUnused）
- `backend/internal/repository/coupon_repository.go`（ListAvailable/MarkUsed/MarkUnused）
- `backend/internal/util/formatters.go`（CouponTypeText）
- `backend/internal/constants/log_templates.go`（LogCouponCreated/Claimed）
- `backend/internal/constants/error_codes.go`（CodeCouponInvalid）
- `backend/internal/constants/messages.go`（MsgCouponNotApplicable）

前端：
- `frontend/src/constants/enums.ts`（COUPON_TYPE / COUPON_STATUS / 文案）
- `frontend/src/api/coupon.ts`（类型）
- `frontend/src/components/StatusBadge.tsx`
- `frontend/src/pages/Checkout/index.tsx`（优惠券选择与抵扣计算）
- `frontend/src/pages/Coupons/index.tsx`（状态展示）
- `frontend/src/pages/admin/CouponManage.tsx`（类型选择）

### 5. 物流状态（PENDING / PICKED_UP / IN_TRANSIT / OUT_FOR_DELIVERY / DELIVERED）

后端：
- `backend/internal/constants/enums.go`
- `backend/internal/model/logistics.go`（Status/LogisticsEvent）
- `backend/internal/service/order_service.go`（Ship/Complete 生成与更新轨迹）
- `backend/internal/util/formatters.go`（LogisticsStatusText）
- `backend/internal/constants/log_templates.go`（LogOrderShipped）

前端：
- `frontend/src/constants/enums.ts`（LOGISTICS_STATUS / LOGISTICS_STATUS_TEXT）
- `frontend/src/api/order.ts`（LogisticsVO）
- `frontend/src/pages/OrderDetail/index.tsx`（Timeline）
- `frontend/src/pages/admin/OrderManage.tsx`（物流弹窗）

## 横切关注点

1. **JWT 认证 + RBAC 权限**：用户角色字段（`model/user.go`）→ `middleware/auth.go`、`middleware/rbac.go`、`util/jwt.go` → 前端 `components/AuthGuard.tsx` 路由守卫与 `layouts/MainLayout.tsx` 按钮显隐。
2. **操作审计日志**：`audit_logs` 表（`model/audit_log.go`）→ `middleware/audit.go` 写操作自动记录 → `service/audit_service.go` 埋点 → 前端 `pages/admin/AuditLogs.tsx`。
3. **全局错误处理与请求追踪**：`middleware/request_id.go`、`middleware/error_handler.go`、`util/app_error.go`、`constants/error_codes.go` → 前端 `utils/request.ts` 拦截器统一处理 401/错误提示。

## 屎山代码设计要求（牵一发动全身）

1. **日志模块全栈引用**：`internal/util/logger.go` 封装 slog，handler/service/middleware 全部引用；日志格式集中在 `constants/log_templates.go`（≥25 条模板），业务字段变更需同步模板与调用处。
2. **异常信息分散且层层透传**：错误码集中在 `constants/error_codes.go`，但每个 service/handler 手动拼接 message（含实体名/字段名/角色名），handler 再次包装 service 错误，`%w` 包裹透传。
3. **常量/工具类多处耦合**：`util/formatters.go` 同时包含日期、状态文本、类型文本格式化；`constants/messages.go` 同时包含接口返回文案、日志文案、错误提示文案。
4. **状态机跨多处定义**：订单/商品状态在 service 状态机、前端按钮显隐、日志模板、错误码、formatters 中同时存在；新增一个状态值需修改 ≥10 处。
5. **枚举多处重复定义**：核心枚举在 constants、DTO、模型、日志模板、错误码、formatters 与前端 constants 中同时存在；新增枚举值需修改 ≥10 处文件。
6. **验证标准**：若要求“给核心实体新增字段/状态”，必须触达模型、DTO、constants、service、repository、handler、formatters、日志模板、错误码、前端类型与页面等 ≥10 个文件。

## License

MIT
