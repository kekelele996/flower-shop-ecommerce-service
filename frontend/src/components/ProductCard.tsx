import { Card, Tag } from 'antd';
import { useNavigate } from 'react-router-dom';
import type { ProductVO } from '../api/product';
import { formatPrice } from '../utils/format';

interface ProductCardProps {
  product: ProductVO;
}

// ProductCard：商品卡片（首页/推荐/搜索复用）。
export default function ProductCard({ product }: ProductCardProps) {
  const navigate = useNavigate();
  return (
    <Card
      hoverable
      cover={
        <img
          alt={product.name}
          src={product.cover_image || 'https://placehold.co/600x400?text=Flower'}
          style={{ height: 180, objectFit: 'cover' }}
          onClick={() => navigate(`/products/${product.id}`)}
        />
      }
    >
      <Card.Meta
        title={<span onClick={() => navigate(`/products/${product.id}`)}>{product.name}</span>}
        description={
          <div>
            <div style={{ color: '#999', fontSize: 12, marginBottom: 6 }}>{product.sub_title}</div>
            <div style={{ display: 'flex', justifyContent: 'space-between', alignItems: 'center' }}>
              <span style={{ color: '#e4393c', fontWeight: 600, fontSize: 18 }}>{formatPrice(product.price)}</span>
              {product.free_shipping && <Tag color="green">包邮</Tag>}
            </div>
            <div style={{ fontSize: 12, color: '#999', marginTop: 6 }}>
              已售 {product.sales} · 评分 {product.rating}
            </div>
          </div>
        }
      />
    </Card>
  );
}
