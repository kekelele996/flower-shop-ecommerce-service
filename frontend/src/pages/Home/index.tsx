import { useEffect, useState } from 'react';
import { Row, Col, Spin, Select, Space, Button, Radio, Divider, Typography } from 'antd';
import { useSearchParams } from 'react-router-dom';
import ProductCard from '../../components/ProductCard';
import EmptyState from '../../components/EmptyState';
import { useProductStore } from '../../stores/productStore';

const { Title } = Typography;

const SORTS = [
  { label: '综合', value: 'default' },
  { label: '销量优先', value: 'sales' },
  { label: '评分优先', value: 'rating' },
  { label: '价格从低到高', value: 'price_asc' },
  { label: '价格从高到低', value: 'price_desc' },
];

// 首页：商品浏览 + 分类 + 搜索 + 推荐。
export default function Home() {
  const [searchParams, setSearchParams] = useSearchParams();
  const { products, total, loading, recommendations, fetchProducts, fetchRecommendations } = useProductStore();
  const [sort, setSort] = useState('default');
  const [freeShipping, setFreeShipping] = useState<boolean | undefined>(undefined);

  const keyword = searchParams.get('keyword') || '';
  const categoryId = searchParams.get('category_id') || '';

  useEffect(() => {
    fetchProducts({
      page: 1,
      page_size: 12,
      keyword: keyword || undefined,
      category_id: categoryId ? Number(categoryId) : undefined,
      sort: sort === 'default' ? undefined : sort,
      free_shipping: freeShipping,
    });
  }, [keyword, categoryId, sort, freeShipping, fetchProducts]);

  useEffect(() => {
    fetchRecommendations();
  }, [fetchRecommendations]);

  return (
    <div>
      <Space wrap style={{ marginBottom: 16 }}>
        <Select value={sort} onChange={setSort} options={SORTS} style={{ width: 160 }} placeholder="排序" />
        <Radio.Group
          value={freeShipping === undefined ? 'all' : freeShipping ? 'yes' : 'no'}
          onChange={(e) => {
            const v = e.target.value;
            setFreeShipping(v === 'all' ? undefined : v === 'yes');
          }}
          optionType="button"
          buttonStyle="solid"
        >
          <Radio.Button value="all">全部</Radio.Button>
          <Radio.Button value="yes">包邮</Radio.Button>
          <Radio.Button value="no">不包邮</Radio.Button>
        </Radio.Group>
        {(keyword || categoryId) && <Button onClick={() => setSearchParams({})}>清除筛选</Button>}
      </Space>

      {keyword && <Typography.Paragraph type="secondary">搜索 “{keyword}” 共 {total} 条结果</Typography.Paragraph>}

      <Spin spinning={loading}>
        {products.length === 0 ? (
          <EmptyState description="没有找到相关商品" />
        ) : (
          <Row gutter={[16, 16]}>
            {products.map((p) => (
              <Col key={p.id} xs={12} sm={8} md={6} lg={6}>
                <ProductCard product={p} />
              </Col>
            ))}
          </Row>
        )}
      </Spin>

      {recommendations.length > 0 && (
        <div style={{ marginTop: 32 }}>
          <Divider>
            <Title level={4} style={{ margin: 0 }}>
              猜你喜欢
            </Title>
          </Divider>
          <Row gutter={[16, 16]}>
            {recommendations.slice(0, 4).map((p) => (
              <Col key={p.id} xs={12} sm={8} md={6} lg={6}>
                <ProductCard product={p} />
              </Col>
            ))}
          </Row>
        </div>
      )}
    </div>
  );
}
