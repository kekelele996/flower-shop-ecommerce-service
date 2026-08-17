import { useEffect, useState } from 'react';
import { useParams, useNavigate } from 'react-router-dom';
import { Row, Col, Image, Button, InputNumber, Space, Descriptions, Tag, Divider, Card, List, Rate, message } from 'antd';
import { ShoppingCartOutlined, ThunderboltOutlined } from '@ant-design/icons';
import { getProduct, recordView } from '../../api/product';
import { addCart } from '../../api/cart';
import { listReviews } from '../../api/review';
import { formatPrice } from '../../utils/format';
import StarRating from '../../components/StarRating';
import EmptyState from '../../components/EmptyState';
import type { ProductVO } from '../../api/product';
import type { ReviewVO } from '../../api/review';
import { useUserStore } from '../../stores/userStore';

// 商品详情：多图轮播、加购、立即购买、评价列表。
export default function ProductDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const { isLogin } = useUserStore();
  const [product, setProduct] = useState<ProductVO | null>(null);
  const [reviews, setReviews] = useState<ReviewVO[]>([]);
  const [quantity, setQuantity] = useState(1);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    if (!id) return;
    getProduct(Number(id))
      .then(setProduct)
      .catch(() => setProduct(null))
      .finally(() => setLoading(false));
    recordView(Number(id)).catch(() => undefined);
    listReviews({ product_id: Number(id), page: 1, page_size: 5 })
      .then((d) => setReviews(d.list))
      .catch(() => setReviews([]));
  }, [id]);

  const handleAddCart = async () => {
    if (!isLogin()) {
      navigate('/login');
      return;
    }
    await addCart({ product_id: Number(id), quantity });
    message.success('已加入购物车');
  };

  const handleBuyNow = async () => {
    if (!isLogin()) {
      navigate('/login');
      return;
    }
    await addCart({ product_id: Number(id), quantity });
    navigate('/checkout');
  };

  if (loading) return <div style={{ padding: 64, textAlign: 'center' }}>加载中...</div>;
  if (!product) return <EmptyState description="商品不存在或已下架" />;

  return (
    <div>
      <Row gutter={32}>
        <Col xs={24} md={10}>
          <Image
            src={product.cover_image}
            alt={product.name}
            width="100%"
            height={380}
            style={{ objectFit: 'cover', borderRadius: 8 }}
            fallback="https://placehold.co/600x400?text=Flower"
          />
          <Space style={{ marginTop: 8 }} wrap>
            {(product.images || []).map((img, idx) => (
              <Image key={idx} src={img} width={72} height={72} style={{ objectFit: 'cover', borderRadius: 4 }} />
            ))}
          </Space>
        </Col>
        <Col xs={24} md={14}>
          <h2 style={{ fontSize: 24, marginBottom: 8 }}>{product.name}</h2>
          <div style={{ color: '#999', marginBottom: 12 }}>{product.sub_title}</div>
          <div style={{ background: '#fff7e6', padding: 16, borderRadius: 8, marginBottom: 16 }}>
            <div style={{ fontSize: 28, color: '#e4393c', fontWeight: 700 }}>
              {formatPrice(product.price)}
              {product.original_price > product.price && (
                <span style={{ fontSize: 14, color: '#999', textDecoration: 'line-through', marginLeft: 12, fontWeight: 400 }}>
                  {formatPrice(product.original_price)}
                </span>
              )}
            </div>
            <Space size="large" style={{ marginTop: 8 }}>
              <span>已售 {product.sales}</span>
              <span>
                <StarRating value={Math.round(product.rating)} /> {product.rating}
              </span>
              <span>库存 {product.stock}</span>
            </Space>
          </div>
          <Descriptions column={2} size="small">
            <Descriptions.Item label="发货地">{product.shipping_from || '-'}</Descriptions.Item>
            <Descriptions.Item label="配送">
              {product.free_shipping ? <Tag color="green">包邮</Tag> : '运费 ¥8'}
            </Descriptions.Item>
            <Descriptions.Item label="分类">{product.category_name || '-'}</Descriptions.Item>
            <Descriptions.Item label="状态">
              {product.status === 'ON_SALE' ? <Tag color="green">在售</Tag> : <Tag color="red">下架</Tag>}
            </Descriptions.Item>
          </Descriptions>
          <Space style={{ marginTop: 24 }} size="large">
            <Space>
              <span>数量</span>
              <InputNumber min={1} max={Math.max(product.stock, 1)} value={quantity} onChange={(v) => setQuantity(v || 1)} />
            </Space>
            <Button type="primary" size="large" icon={<ShoppingCartOutlined />} onClick={handleAddCart} disabled={product.stock <= 0}>
              加入购物车
            </Button>
            <Button type="primary" danger size="large" icon={<ThunderboltOutlined />} onClick={handleBuyNow} disabled={product.stock <= 0}>
              立即购买
            </Button>
          </Space>
        </Col>
      </Row>

      <Divider orientation="left">商品详情</Divider>
      <Card>
        <div dangerouslySetInnerHTML={{ __html: product.detail || product.description || '<p>暂无详情</p>' }} />
      </Card>

      <Divider orientation="left">用户评价（{reviews.length}）</Divider>
      {reviews.length === 0 ? (
        <EmptyState description="暂无评价" />
      ) : (
        <List
          dataSource={reviews}
          renderItem={(r) => (
            <List.Item>
              <List.Item.Meta
                title={
                  <Space>
                    <span>{r.username}</span>
                    <Rate disabled value={r.rating} style={{ fontSize: 14 }} />
                  </Space>
                }
                description={
                  <div>
                    <div>{r.content}</div>
                    {r.reply && (
                      <div style={{ background: '#f6ffed', padding: 8, marginTop: 8, borderRadius: 4 }}>
                        商家回复：{r.reply}
                      </div>
                    )}
                    <div style={{ color: '#999', fontSize: 12, marginTop: 8 }}>{r.created_at}</div>
                  </div>
                }
              />
            </List.Item>
          )}
        />
      )}
    </div>
  );
}
