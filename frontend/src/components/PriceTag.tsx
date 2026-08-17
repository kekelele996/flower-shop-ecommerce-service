import { formatPrice } from '../utils/format';

interface PriceTagProps {
  price: number;
  originalPrice?: number;
  size?: 'small' | 'large';
}

// PriceTag：价格展示（购物车/订单/详情复用）。
export default function PriceTag({ price, originalPrice, size = 'small' }: PriceTagProps) {
  const fontSize = size === 'large' ? 26 : 16;
  return (
    <span>
      <span style={{ color: '#e4393c', fontWeight: 600, fontSize }}>{formatPrice(price)}</span>
      {originalPrice && originalPrice > price ? (
        <span style={{ color: '#999', textDecoration: 'line-through', marginLeft: 8, fontSize: 12 }}>{formatPrice(originalPrice)}</span>
      ) : null}
    </span>
  );
}
