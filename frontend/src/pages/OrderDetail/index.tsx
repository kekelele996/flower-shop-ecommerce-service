import { useEffect, useState } from 'react';
import { useParams, useNavigate, useSearchParams } from 'react-router-dom';
import { Card, Descriptions, Button, Space, Steps, Timeline, Rate, Input, Modal, message, Divider, Typography } from 'antd';
import { getOrder, payOrder, cancelOrder, completeOrder, getLogistics } from '../../api/order';
import { createReview } from '../../api/review';
import { formatPrice } from '../../utils/format';
import StatusBadge from '../../components/StatusBadge';
import EmptyState from '../../components/EmptyState';
import { showConfirm } from '../../components/ConfirmDialog';
import type { OrderVO, LogisticsVO } from '../../api/order';
import { ORDER_STATUS } from '../../constants/enums';

// 订单详情：状态流转 + 物流轨迹 + 评价。
export default function OrderDetail() {
  const { id } = useParams();
  const navigate = useNavigate();
  const [searchParams, setSearchParams] = useSearchParams();
  const [order, setOrder] = useState<OrderVO | null>(null);
  const [logistics, setLogistics] = useState<LogisticsVO | null>(null);
  const [reviewOpen, setReviewOpen] = useState(false);
  const [reviewTarget, setReviewTarget] = useState<{ orderItemId: number; productName: string } | null>(null);
  const [rating, setRating] = useState(5);
  const [content, setContent] = useState('');

  useEffect(() => {
    if (!id) return;
    getOrder(Number(id)).then(setOrder).catch(() => setOrder(null));
    getLogistics(Number(id)).then(setLogistics).catch(() => setLogistics(null));
  }, [id]);

  const reload = async () => {
    if (!id) return;
    const o = await getOrder(Number(id));
    setOrder(o);
    if (o.status !== ORDER_STATUS.PENDING_PAYMENT) {
      getLogistics(Number(id)).then(setLogistics).catch(() => setLogistics(null));
    }
  };

  const handlePay = async () => {
    if (!id) return;
    await payOrder(Number(id));
    message.success('支付成功（模拟支付宝沙箱）');
    setSearchParams({});
    reload();
  };

  const handleCancel = () => {
    if (!id) return;
    showConfirm({
      title: '取消订单',
      content: '确定取消该订单吗？取消后库存将回补。',
      onOk: async () => {
        await cancelOrder(Number(id), '用户取消');
        message.success('订单已取消');
        reload();
      },
    });
  };

  const handleComplete = async () => {
    if (!id) return;
    await completeOrder(Number(id));
    message.success('已确认收货');
    reload();
  };

  const openReview = (orderItemId: number, productName: string) => {
    setReviewTarget({ orderItemId, productName });
    setRating(5);
    setContent('');
    setReviewOpen(true);
  };

  const submitReview = async () => {
    if (!reviewTarget) return;
    await createReview({ order_item_id: reviewTarget.orderItemId, rating, content });
    message.success('评价成功');
    setReviewOpen(false);
    reload();
  };

  if (!order) return <EmptyState description="订单不存在" />;

  const stepIndex =
    order.status === ORDER_STATUS.PENDING_PAYMENT
      ? 0
      : order.status === ORDER_STATUS.PENDING_SHIPMENT
        ? 1
        : order.status === ORDER_STATUS.SHIPPED
          ? 2
          : order.status === ORDER_STATUS.COMPLETED
            ? 3
            : 0;

  return (
    <div style={{ maxWidth: 960, margin: '0 auto' }}>
      <Space style={{ marginBottom: 16 }}>
        <Button onClick={() => navigate('/orders')}>返回订单列表</Button>
        <Typography.Title level={4} style={{ margin: 0 }}>
          订单详情 {order.order_no}
        </Typography.Title>
        <StatusBadge status={order.status} kind="order" />
      </Space>

      <Card title="订单进度" style={{ marginBottom: 16 }}>
        {order.status === ORDER_STATUS.CANCELLED ? (
          <Typography.Text type="danger">订单已取消</Typography.Text>
        ) : (
          <Steps
            current={stepIndex}
            items={[
              { title: '提交订单' },
              { title: '付款' },
              { title: '发货' },
              { title: '收货完成' },
            ]}
          />
        )}
        <Divider style={{ margin: '16px 0' }} />
        <Space>
          {order.status === ORDER_STATUS.PENDING_PAYMENT && (
            <>
              <Button type="primary" onClick={handlePay}>
                立即支付（模拟支付宝沙箱）
              </Button>
              <Button danger onClick={handleCancel}>
                取消订单
              </Button>
            </>
          )}
          {order.status === ORDER_STATUS.SHIPPED && (
            <Button type="primary" onClick={handleComplete}>
              确认收货
            </Button>
          )}
        </Space>
      </Card>

      <Card title="商品清单" style={{ marginBottom: 16 }}>
        {order.items?.map((item) => (
          <div key={item.id} style={{ display: 'flex', justifyContent: 'space-between', padding: '8px 0', borderBottom: '1px solid #f0f0f0' }}>
            <Space>
              <img src={item.product_image} alt={item.product_name} style={{ width: 48, height: 48, objectFit: 'cover', borderRadius: 4 }} />
              <div>
                <div>{item.product_name}</div>
                <div style={{ color: '#999', fontSize: 12 }}>
                  ¥{item.price} × {item.quantity}
                </div>
              </div>
            </Space>
            <Space>
              <span style={{ fontWeight: 600 }}>{formatPrice(item.total_price)}</span>
              {order.status === ORDER_STATUS.COMPLETED && !item.reviewed && (
                <Button size="small" type="primary" onClick={() => openReview(item.id, item.product_name)}>
                  评价
                </Button>
              )}
            </Space>
          </div>
        ))}
        <div style={{ textAlign: 'right', marginTop: 12 }}>
          <div>商品总额：{formatPrice(order.total_amount)}</div>
          <div>优惠券抵扣：-{formatPrice(order.discount_amount)}</div>
          <div>运费：{order.shipping_fee === 0 ? '免运费' : formatPrice(order.shipping_fee)}</div>
          <div style={{ fontSize: 18, fontWeight: 700, color: '#e4393c' }}>实付：{formatPrice(order.pay_amount)}</div>
        </div>
      </Card>

      <Card title="收货信息" style={{ marginBottom: 16 }}>
        <Descriptions column={2}>
          <Descriptions.Item label="收货人">{order.receiver_name}</Descriptions.Item>
          <Descriptions.Item label="手机号">{order.receiver_phone}</Descriptions.Item>
          <Descriptions.Item label="地址" span={2}>
            {order.receiver_addr}
          </Descriptions.Item>
          <Descriptions.Item label="备注" span={2}>
            {order.remark || '-'}
          </Descriptions.Item>
        </Descriptions>
      </Card>

      {logistics && (
        <Card title="物流跟踪" style={{ marginBottom: 16 }}>
          <div style={{ marginBottom: 12 }}>
            运单号：{logistics.tracking_no}（{logistics.carrier}） <StatusBadge status={logistics.status} />
          </div>
          <Timeline
            items={[...(logistics.events || [])].reverse().map((e) => ({
              children: (
                <div>
                  <div>{e.desc}</div>
                  <div style={{ color: '#999', fontSize: 12 }}>{e.time}</div>
                </div>
              ),
            }))}
          />
        </Card>
      )}

      <Modal
        title={`评价「${reviewTarget?.productName || ''}」`}
        open={reviewOpen}
        onCancel={() => setReviewOpen(false)}
        onOk={submitReview}
        okText="提交评价"
      >
        <Space direction="vertical" style={{ width: '100%' }}>
          <div>
            评分：
            <Rate value={rating} onChange={setRating} />
          </div>
          <Input.TextArea rows={4} value={content} onChange={(e) => setContent(e.target.value)} placeholder="写下你的使用感受（至少 2 字）" />
        </Space>
      </Modal>
    </div>
  );
}
