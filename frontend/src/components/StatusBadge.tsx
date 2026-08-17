import { Tag } from 'antd';
import { ORDER_STATUS, PRODUCT_STATUS, COUPON_STATUS, ORDER_STATUS_TEXT, PRODUCT_STATUS_TEXT, COUPON_STATUS_TEXT } from '../constants/enums';

const COLOR_MAP: Record<string, string> = {
  [ORDER_STATUS.PENDING_PAYMENT]: 'orange',
  [ORDER_STATUS.PENDING_SHIPMENT]: 'blue',
  [ORDER_STATUS.SHIPPED]: 'cyan',
  [ORDER_STATUS.COMPLETED]: 'green',
  [ORDER_STATUS.CANCELLED]: 'default',
  [PRODUCT_STATUS.ON_SALE]: 'green',
  [PRODUCT_STATUS.OFF_SALE]: 'red',
  [COUPON_STATUS.UNUSED]: 'blue',
  [COUPON_STATUS.USED]: 'default',
  [COUPON_STATUS.EXPIRED]: 'red',
};

const TEXT_MAP: Record<string, Record<string, string>> = {
  order: ORDER_STATUS_TEXT,
  product: PRODUCT_STATUS_TEXT,
  coupon: COUPON_STATUS_TEXT,
};

interface StatusBadgeProps {
  status: string;
  kind?: 'order' | 'product' | 'coupon' | 'default';
}

// StatusBadge：跨页面状态徽标（订单/商品/优惠券复用）。
export default function StatusBadge({ status, kind = 'default' }: StatusBadgeProps) {
  const textMap = TEXT_MAP[kind] || {};
  const text = textMap[status] || status;
  return <Tag color={COLOR_MAP[status] || 'default'}>{text}</Tag>;
}
