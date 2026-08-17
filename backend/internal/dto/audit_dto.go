package dto

// AuditLogQuery 审计日志查询。
type AuditLogQuery struct {
	Page     int    `form:"page"`
	PageSize int    `form:"page_size"`
	Username string `form:"username"`
	Action   string `form:"action"`
	Entity   string `form:"entity"`
}

// AuditLogVO 审计日志视图对象。
type AuditLogVO struct {
	ID        uint   `json:"id"`
	UserID    uint   `json:"user_id"`
	Username  string `json:"username"`
	Action    string `json:"action"`
	Method    string `json:"method"`
	Path      string `json:"path"`
	Entity    string `json:"entity"`
	EntityID  string `json:"entity_id"`
	Detail    string `json:"detail"`
	RequestID string `json:"request_id"`
	IP        string `json:"ip"`
	CreatedAt string `json:"created_at"`
}

// StatsVO 销售统计视图对象。
type StatsVO struct {
	TotalSales       float64          `json:"total_sales"`
	OrderCount       int64            `json:"order_count"`
	ProductCount     int64            `json:"product_count"`
	UserCount        int64            `json:"user_count"`
	PendingShipment  int64            `json:"pending_shipment"`
	OrderStatusCount map[string]int64 `json:"order_status_count"`
	TopProducts      []TopProductVO   `json:"top_products"`
	DailySales       []DailySalesVO   `json:"daily_sales"`
}

// TopProductVO 热销商品。
type TopProductVO struct {
	ProductID   uint    `json:"product_id"`
	ProductName string  `json:"product_name"`
	Sales       int     `json:"sales"`
	Revenue     float64 `json:"revenue"`
}

// DailySalesVO 每日销售额。
type DailySalesVO struct {
	Date  string  `json:"date"`
	Sales float64 `json:"sales"`
}
