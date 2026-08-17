import { useEffect, useState } from 'react';
import { Tabs, Table, Button, Space, message } from 'antd';
import { useNavigate, useSearchParams } from 'react-router-dom';
import { listOrders, cancelOrder, payOrder, completeOrder } from '../../api/order';
import { formatPrice } from '../../utils/format';
import StatusBadge from '../../components/StatusBadge';
import DataTable from '../../components/DataTable';
import { usePagination } from '../../hooks/usePagination';
import { showConfirm } from '../../components/ConfirmDialog';
import type { OrderVO } from '../../api/order';
import { ORDER_STATUS } from '../../constants/enums';

const TABS = [
  { key: '', label: '全部' },
  { key: ORDER_STATUS.PENDING_PAYMENT, label: '待付款' },
  { key: ORDER_STATUS.PENDING_SHIPMENT, label: '待发货' },
  { key: ORDER_STATUS.SHIPPED, label: '已发货' },
  { key: ORDER_STATUS.COMPLETED, label: '已完成' },
  { key: ORDER_STATUS.CANCELLED, label: '已取消' },
];

// 我的订单：状态筛选、支付/取消/确认收货。
export default function Orders() {
  const navigate = useNavigate();
  const [searchParams] = useSearchParams();
  const { page, pageSize, setPage, setPageSize, total, setTotal, reset } = usePagination(10);
  const [status, setStatus] = useState(searchParams.get('status') || '');
  const [orders, setOrders] = useState<OrderVO[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchOrders = async (p = page, ps = pageSize, st = status) => {
    setLoading(true);
    try {
      const data = await listOrders({ page: p, page_size: ps, status: st || undefined });
      setOrders(data.list);
      setTotal(data.total);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchOrders(1, pageSize, status);
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [status, pageSize]);

  const handlePay = async (order: OrderVO) => {
    await payOrder(order.id);
    message.success('支付成功（模拟支付宝沙箱）');
    fetchOrders();
  };

  const handleCancel = (order: OrderVO) => {
    showConfirm({
      title: '取消订单',
      content: `确定取消订单 ${order.order_no} 吗？`,
      onOk: async () => {
        await cancelOrder(order.id, '用户取消');
        message.success('订单已取消');
        fetchOrders();
      },
    });
  };

  const handleComplete = async (order: OrderVO) => {
    await completeOrder(order.id);
    message.success('已确认收货');
    fetchOrders();
  };

  const columns = [
    {
      title: '订单号',
      dataIndex: 'order_no',
      render: (v: string, r: OrderVO) => <a onClick={() => navigate(`/orders/${r.id}`)}>{v}</a>,
    },
    { title: '商品', dataIndex: 'items', render: (_: unknown, r: OrderVO) => <span>{r.items?.map((i) => i.product_name).join('、')}</span> },
    { title: '实付', dataIndex: 'pay_amount', render: (v: number) => <span style={{ color: '#e4393c', fontWeight: 600 }}>{formatPrice(v)}</span> },
    { title: '状态', dataIndex: 'status', render: (v: string) => <StatusBadge status={v} kind="order" /> },
    { title: '下单时间', dataIndex: 'created_at' },
    {
      title: '操作',
      render: (_: unknown, r: OrderVO) => (
        <Space>
          <Button size="small" onClick={() => navigate(`/orders/${r.id}`)}>
            详情
          </Button>
          {r.status === ORDER_STATUS.PENDING_PAYMENT && (
            <>
              <Button size="small" type="primary" onClick={() => handlePay(r)}>
                支付
              </Button>
              <Button size="small" danger onClick={() => handleCancel(r)}>
                取消
              </Button>
            </>
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
      <Tabs
        activeKey={status}
        items={TABS.map((t) => ({ key: t.key, label: t.label }))}
        onChange={(key) => {
          setStatus(key);
          reset();
        }}
      />
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
          fetchOrders(p, ps, status);
        }}
      />
    </div>
  );
}
