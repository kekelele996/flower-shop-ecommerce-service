// 格式化工具：价格、日期、数量。
export function formatPrice(price?: number): string {
  if (price === undefined || price === null) return '¥0.00';
  return `¥${price.toFixed(2)}`;
}

export function formatDateTime(value?: string): string {
  if (!value) return '-';
  return value;
}

export function formatDate(value?: string): string {
  if (!value) return '-';
  return value.slice(0, 10);
}

export function formatQuantity(q?: number): string {
  return q ? `${q}件` : '-';
}

export function formatPercent(rate?: number): string {
  if (!rate) return '';
  return `${Math.round((1 - rate) * 10)}折`;
}
