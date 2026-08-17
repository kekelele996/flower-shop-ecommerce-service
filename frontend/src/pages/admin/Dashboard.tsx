import { useEffect, useState } from 'react';
import { Row, Col, Card, Statistic, Table, Typography } from 'antd';
import { getStats } from '../../api/admin';
import { formatPrice } from '../../utils/format';
import type { StatsVO } from '../../api/admin';

// 数据看板：销售统计。
export default function Dashboard() {
  const [stats, setStats] = useState<StatsVO | null>(null);

  useEffect(() => {
    getStats().then(setStats).catch(() => undefined);
  }, []);

  return (
    <div>
      <Typography.Title level={4}>数据看板</Typography.Title>
      <Row gutter={16}>
        <Col span={6}>
          <Card>
            <Statistic title="总销售额" value={stats?.total_sales || 0} precision={2} prefix="¥" valueStyle={{ color: '#cf1322' }} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="订单总数" value={stats?.order_count || 0} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="商品数" value={stats?.product_count || 0} />
          </Card>
        </Col>
        <Col span={6}>
          <Card>
            <Statistic title="用户数" value={stats?.user_count || 0} />
          </Card>
        </Col>
      </Row>

      <Row gutter={16} style={{ marginTop: 16 }}>
        <Col span={12}>
          <Card title="订单状态分布">
            {Object.entries(stats?.order_status_count || {}).map(([k, v]) => (
              <div key={k} style={{ display: 'flex', justifyContent: 'space-between', padding: '8px 0' }}>
                <span>{k}</span>
                <span>{v}</span>
              </div>
            ))}
          </Card>
        </Col>
        <Col span={12}>
          <Card title="热销商品 Top5">
            <Table
              rowKey="product_id"
              size="small"
              pagination={false}
              dataSource={stats?.top_products || []}
              columns={[
                { title: '商品', dataIndex: 'product_name' },
                { title: '销量', dataIndex: 'sales' },
                { title: '成交额', dataIndex: 'revenue', render: (v: number) => formatPrice(v) },
              ]}
            />
          </Card>
        </Col>
      </Row>
    </div>
  );
}
