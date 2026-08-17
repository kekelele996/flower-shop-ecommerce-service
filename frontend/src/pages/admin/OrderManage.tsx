import { useEffect, useState } from 'react';
import { Button, Space, Modal, Descriptions, Timeline, message, Select, Card, Typography } from 'antd';
import { adminListOrders, shipOrder, adminGetLogistics, adminCompleteOrder } from '../../api/admin';
import DataTable from '../../components/DataTable';
import StatusBadge from '../../components/StatusBadge';
import { usePagination } from '../../hooks/usePagination';
import { formatPrice } from '../../utils/format';
import type { OrderVO, LogisticsVO } from '../../api/order';
import { ORDER_STATUS } from '../../constants/enums';

// 订单管理：发货、确认收货、物流查看。
export default function OrderManage() {
  const { page, pageSize, setPage, setPageSize, total, setTotal } = usePagination(10);
  const [orders, setOrders] = useState<OrderVO[]>([]);
  const [loading, setLoading] = useState(false);
  const [status, setStatus] = useState<string | undefined>(undefined);
  const [logisticsOrder, setLogisticsOrder] = useState<OrderVO | null>(null);
  const [logistics, setLogistics] = useState<LogisticsVO | null>(null);

  const fetchData = async (p = page, ps = pageSize, st = status) => {
    setLoading(true);
    try {
      const data = await adminListOrders({ page: p, page_size: ps, status: st });
      setOrders(data.list);
      setTotal(data.total);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pageSize, status]);

  const handleShip = async (order: OrderVO) => {
    await shipOrder(order.id, 'SF_EXPRESS');
    message.success(`订单 ${order.order_no} 已发货`);
    fetchData();
  };

  const handleComplete = async (order: OrderVO) => {
    await adminCompleteOrder(order.id);
    message.success('已确认收货');
    fetchData();
  };

  const showLogistics = async (order: OrderVO) => {
    setLogisticsOrder(order);
    try {
      const data = await adminGetLogistics(order.id);
      setLogistics(data);
    } catch {
      setLogistics(null);
    }
  };

  const columns = [
    { title: '订单号', dataIndex: 'order_no' },
    { title: '用户ID', dataIndex: 'user_id' },
    { title: '商品', dataIndex: 'items', render: (_: unknown, r: OrderVO) => <span>{r.items?.map((i) => i.product_name).join('、')}</span> },
    { title: '实付', dataIndex: 'pay_amount', render: (v: number) => formatPrice(v) },
    { title: '状态', dataIndex: 'status', render: (v: string) => <StatusBadge status={v} kind="order" /> },
    { title: '下单时间', dataIndex: 'created_at' },
    {
      title: '操作',
      render: (_: unknown, r: OrderVO) => (
        <Space>
          <Button size="small" onClick={() => showLogistics(r)}>
            物流
          </Button>
          {r.status === ORDER_STATUS.PENDING_SHIPMENT && (
            <Button size="small" type="primary" onClick={() => handleShip(r)}>
              发货
            </Button>
          )}
          {r.status === ORDER_STATUS.SHIPPED && (
            <Button size="small" type="primary" onClick={() => handleComplete(r)}>
              确认收货
            </Button>
          )}
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Card>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
          <Typography.Title level={4} style={{ margin: 0 }}>
            订单管理
          </Typography.Title>
          <Select
            allowClear
            placeholder="按状态筛选"
            style={{ width: 160 }}
            value={status}
            onChange={(v) => {
              setStatus(v);
              setPage(1);
            }}
            options={Object.entries(ORDER_STATUS).map(([label, value]) => ({ label, value }))}
          />
        </div>
        <DataTable
          rowKey="id"
          loading={loading}
          dataSource={orders}
          columns={columns}
          total={total}
          page={page}
          pageSize={pageSize}
          onPageChange={(p, ps) => {
            setPage(p);
            setPageSize(ps);
            fetchData(p, ps);
          }}
        />
      </Card>

      <Modal
        title={logisticsOrder ? `物流信息 ${logisticsOrder.order_no}` : '物流信息'}
        open={!!logisticsOrder}
        onCancel={() => setLogisticsOrder(null)}
        footer={null}
      >
        {logistics ? (
          <div>
            <Descriptions column={1} size="small">
              <Descriptions.Item label="运单号">{logistics.tracking_no}</Descriptions.Item>
              <Descriptions.Item label="承运商">{logistics.carrier}</Descriptions.Item>
              <Descriptions.Item label="状态">
                <StatusBadge status={logistics.status} />
              </Descriptions.Item>
            </Descriptions>
            <Timeline
              style={{ marginTop: 12 }}
              items={[...(logistics.events || [])].reverse().map((e) => ({
                children: (
                  <div>
                    <div>{e.desc}</div>
                    <div style={{ color: '#999', fontSize: 12 }}>{e.time}</div>
                  </div>
                ),
              }))}
            />
          </div>
        ) : (
          <Typography.Text type="secondary">该订单暂无物流信息</Typography.Text>
        )}
      </Modal>
    </div>
  );
}
