// 与后端 internal/constants/enums.go 对应的业务枚举（新增枚举值需同步 ≥10 处文件）。
export const ROLE = {
  USER: 'USER',
  ADMIN: 'ADMIN',
} as const;

export const USER_STATUS = {
  ACTIVE: 'ACTIVE',
  DISABLED: 'DISABLED',
} as const;

export const PRODUCT_STATUS = {
  ON_SALE: 'ON_SALE',
  OFF_SALE: 'OFF_SALE',
} as const;

export const ORDER_STATUS = {
  PENDING_PAYMENT: 'PENDING_PAYMENT',
  PENDING_SHIPMENT: 'PENDING_SHIPMENT',
  SHIPPED: 'SHIPPED',
  COMPLETED: 'COMPLETED',
  CANCELLED: 'CANCELLED',
} as const;

export const PAYMENT_STATUS = {
  PENDING: 'PENDING',
  SUCCESS: 'SUCCESS',
  FAILED: 'FAILED',
} as const;

export const LOGISTICS_STATUS = {
  PENDING: 'PENDING',
  PICKED_UP: 'PICKED_UP',
  IN_TRANSIT: 'IN_TRANSIT',
  OUT_FOR_DELIVERY: 'OUT_FOR_DELIVERY',
  DELIVERED: 'DELIVERED',
} as const;

export const COUPON_TYPE = {
  FULL_REDUCTION: 'FULL_REDUCTION',
  DISCOUNT: 'DISCOUNT',
} as const;

export const COUPON_STATUS = {
  UNUSED: 'UNUSED',
  USED: 'USED',
  EXPIRED: 'EXPIRED',
} as const;

// 中文文案（与后端 formatters / messages 一致）
export const ORDER_STATUS_TEXT: Record<string, string> = {
  [ORDER_STATUS.PENDING_PAYMENT]: '待付款',
  [ORDER_STATUS.PENDING_SHIPMENT]: '待发货',
  [ORDER_STATUS.SHIPPED]: '已发货',
  [ORDER_STATUS.COMPLETED]: '已完成',
  [ORDER_STATUS.CANCELLED]: '已取消',
};

export const PRODUCT_STATUS_TEXT: Record<string, string> = {
  [PRODUCT_STATUS.ON_SALE]: '在售',
  [PRODUCT_STATUS.OFF_SALE]: '下架',
};

export const COUPON_TYPE_TEXT: Record<string, string> = {
  [COUPON_TYPE.FULL_REDUCTION]: '满减券',
  [COUPON_TYPE.DISCOUNT]: '折扣券',
};

export const COUPON_STATUS_TEXT: Record<string, string> = {
  [COUPON_STATUS.UNUSED]: '未使用',
  [COUPON_STATUS.USED]: '已使用',
  [COUPON_STATUS.EXPIRED]: '已过期',
};

export const LOGISTICS_STATUS_TEXT: Record<string, string> = {
  [LOGISTICS_STATUS.PENDING]: '待揽收',
  [LOGISTICS_STATUS.PICKED_UP]: '已揽收',
  [LOGISTICS_STATUS.IN_TRANSIT]: '运输中',
  [LOGISTICS_STATUS.OUT_FOR_DELIVERY]: '派送中',
  [LOGISTICS_STATUS.DELIVERED]: '已签收',
};
