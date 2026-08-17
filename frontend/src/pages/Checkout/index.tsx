import { useEffect, useMemo, useState } from 'react';
import { Card, Form, Input, Button, Space, Radio, Divider, Typography, message, List } from 'antd';
import { useNavigate } from 'react-router-dom';
import { getCart } from '../../api/cart';
import { checkout } from '../../api/order';
import { getAvailableCoupons } from '../../api/coupon';
import { formatPrice } from '../../utils/format';
import { useCartStore } from '../../stores/cartStore';
import EmptyState from '../../components/EmptyState';
import type { CouponVO } from '../../api/coupon';

// 结算页：选择地址、优惠券、模拟支付。
export default function Checkout() {
  const navigate = useNavigate();
  const { summary, fetchCart, clear } = useCartStore();
  const [available, setAvailable] = useState<CouponVO[]>([]);
  const [couponId, setCouponId] = useState<number | undefined>(undefined);
  const [submitting, setSubmitting] = useState(false);

  useEffect(() => {
    fetchCart();
    getAvailableCoupons().then(setAvailable).catch(() => setAvailable([]));
  }, [fetchCart]);

  const selectedItems = useMemo(() => (summary?.items || []).filter((i) => i.selected), [summary]);

  const calc = useMemo(() => {
    const total = selectedItems.reduce((s, i) => s + i.subtotal, 0);
    const hasFreeShip = selectedItems.some((i) => i.product.free_shipping);
    const shipping = selectedItems.length === 0 || hasFreeShip ? 0 : 8;
    const coupon = available.find((c) => c.id === couponId);
    let discount = 0;
    if (coupon && total >= coupon.threshold) {
      discount = coupon.type === 'FULL_REDUCTION' ? coupon.amount : total * (1 - coupon.discount_rate);
      if (discount > total) discount = total;
    }
    return { total, shipping, discount, pay: total - discount + shipping };
  }, [selectedItems, available, couponId]);

  const handleSubmit = async (values: { receiver_name: string; receiver_phone: string; receiver_addr: string; remark?: string }) => {
    if (selectedItems.length === 0) {
      message.warning('请先选择商品');
      return;
    }
    setSubmitting(true);
    try {
      const res = await checkout({
        cart_item_ids: selectedItems.map((i) => i.id),
        coupon_id: couponId,
        ...values,
      });
      message.success('下单成功，请完成支付');
      clear();
      navigate(`/orders/${res.order.id}?action=pay`);
    } finally {
      setSubmitting(false);
    }
  };

  if (!summary || selectedItems.length === 0) {
    return (
      <EmptyState description="没有可结算的商品">
        <Button type="primary" onClick={() => navigate('/cart')}>
          返回购物车
        </Button>
      </EmptyState>
    );
  }

  return (
    <div style={{ maxWidth: 900, margin: '0 auto' }}>
      <Typography.Title level={4}>确认订单</Typography.Title>
      <Card title="商品清单">
        <List
          dataSource={selectedItems}
          renderItem={(i) => (
            <List.Item>
              <Space>
                <img src={i.product.cover_image} alt={i.product.name} style={{ width: 48, height: 48, objectFit: 'cover', borderRadius: 4 }} />
                <div>
                  <div>{i.product.name}</div>
                  <div style={{ color: '#999', fontSize: 12 }}>{i.product.shipping_from}</div>
                </div>
              </Space>
              <Space>
                <span>¥{i.product.price} × {i.quantity}</span>
                <span style={{ color: '#e4393c', fontWeight: 600 }}>{formatPrice(i.subtotal)}</span>
              </Space>
            </List.Item>
          )}
        />
        <Divider style={{ margin: '12px 0' }} />
        <Space direction="vertical" style={{ width: '100%' }} size={4}>
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <span>商品总额</span>
            <span>{formatPrice(calc.total)}</span>
          </div>
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <span>运费</span>
            <span>{calc.shipping === 0 ? '免运费' : formatPrice(calc.shipping)}</span>
          </div>
          <div style={{ display: 'flex', justifyContent: 'space-between' }}>
            <span>优惠券抵扣</span>
            <span style={{ color: '#e4393c' }}>-{formatPrice(calc.discount)}</span>
          </div>
          <div style={{ display: 'flex', justifyContent: 'space-between', fontSize: 18, fontWeight: 700 }}>
            <span>应付金额</span>
            <span style={{ color: '#e4393c' }}>{formatPrice(calc.pay)}</span>
          </div>
        </Space>
      </Card>

      <Card title="优惠券" style={{ marginTop: 16 }}>
        <Radio.Group value={couponId} onChange={(e) => setCouponId(e.target.value)}>
          <Space direction="vertical">
            <Radio value={undefined}>不使用优惠券</Radio>
            {available.map((c) => (
              <Radio key={c.id} value={c.id}>
                {c.name}（满{c.threshold}减{c.amount || `${Math.round((1 - c.discount_rate) * 10)}折`}）
              </Radio>
            ))}
          </Space>
        </Radio.Group>
      </Card>

      <Card title="收货信息" style={{ marginTop: 16 }}>
        <Form layout="vertical" onFinish={handleSubmit} initialValues={{ receiver_name: '张三', receiver_phone: '13800138000', receiver_addr: '上海市浦东新区花木路 100 号' }}>
          <Space size="large" wrap>
            <Form.Item name="receiver_name" label="收货人" rules={[{ required: true, message: '请输入收货人' }]}>
              <Input style={{ width: 200 }} />
            </Form.Item>
            <Form.Item name="receiver_phone" label="手机号" rules={[{ required: true, message: '请输入手机号' }]}>
              <Input style={{ width: 200 }} />
            </Form.Item>
          </Space>
          <Form.Item name="receiver_addr" label="收货地址" rules={[{ required: true, message: '请输入地址' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="remark" label="订单备注">
            <Input placeholder="选填" />
          </Form.Item>
          <div style={{ textAlign: 'right' }}>
            <Button type="primary" size="large" loading={submitting} htmlType="submit">
              提交订单并去支付
            </Button>
          </div>
        </Form>
      </Card>
    </div>
  );
}
