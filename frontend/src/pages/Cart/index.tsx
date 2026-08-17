import { useEffect } from 'react';
import { Table, Button, InputNumber, Space, Card, message, Typography, Checkbox } from 'antd';
import { useNavigate } from 'react-router-dom';
import { getCart, updateCart, removeCart } from '../../api/cart';
import { useCartStore } from '../../stores/cartStore';
import { formatPrice } from '../../utils/format';
import EmptyState from '../../components/EmptyState';
import PriceTag from '../../components/PriceTag';
import { showConfirm } from '../../components/ConfirmDialog';
import type { CartItemVO } from '../../api/cart';

// 购物车：数量修改、删除、选中结算（复用 CartStore + PriceTag）。
export default function Cart() {
  const navigate = useNavigate();
  const { summary, fetchCart } = useCartStore();

  useEffect(() => {
    fetchCart();
  }, [fetchCart]);

  const handleQuantity = async (item: CartItemVO, quantity: number) => {
    await updateCart(item.id, { quantity });
    message.success('数量已更新');
    fetchCart();
  };

  const handleSelected = async (item: CartItemVO, selected: boolean) => {
    await updateCart(item.id, { selected });
    fetchCart();
  };

  const handleRemove = (item: CartItemVO) => {
    showConfirm({
      title: '移除商品',
      content: `确定从购物车移除「${item.product.name}」吗？`,
      onOk: async () => {
        await removeCart(item.id);
        message.success('已移除');
        fetchCart();
      },
    });
  };

  const selectedItems = (summary?.items || []).filter((i) => i.selected);
  const selectedTotal = selectedItems.reduce((sum, i) => sum + i.subtotal, 0);

  const columns = [
    {
      title: '商品',
      dataIndex: 'product',
      render: (_: unknown, item: CartItemVO) => (
        <Space>
          <img src={item.product.cover_image} alt={item.product.name} style={{ width: 64, height: 64, objectFit: 'cover', borderRadius: 4 }} />
          <div>
            <Typography.Text strong>{item.product.name}</Typography.Text>
            <div style={{ color: '#999', fontSize: 12 }}>{item.product.shipping_from}</div>
          </div>
        </Space>
      ),
    },
    {
      title: '单价',
      dataIndex: 'product',
      render: (_: unknown, item: CartItemVO) => <PriceTag price={item.product.price} />,
    },
    {
      title: '数量',
      dataIndex: 'quantity',
      render: (_: unknown, item: CartItemVO) => (
        <InputNumber min={1} max={99} value={item.quantity} onChange={(v) => handleQuantity(item, v || 1)} />
      ),
    },
    {
      title: '小计',
      dataIndex: 'subtotal',
      render: (v: number) => <span style={{ color: '#e4393c', fontWeight: 600 }}>{formatPrice(v)}</span>,
    },
    {
      title: '操作',
      render: (_: unknown, item: CartItemVO) => <Button danger type="link" onClick={() => handleRemove(item)}>删除</Button>,
    },
    {
      title: '选择',
      dataIndex: 'selected',
      width: 60,
      render: (_: unknown, item: CartItemVO) => (
        <Checkbox checked={item.selected} onChange={(e) => handleSelected(item, e.target.checked)} />
      ),
    },
  ];

  if (!summary || summary.items.length === 0) {
    return (
      <EmptyState description="购物车是空的">
        <Button type="primary" onClick={() => navigate('/')}>
          去逛逛
        </Button>
      </EmptyState>
    );
  }

  return (
    <div>
      <Typography.Title level={4}>我的购物车</Typography.Title>
      <Table rowKey="id" dataSource={summary.items} columns={columns} pagination={false} />
      <Card style={{ marginTop: 16, textAlign: 'right' }}>
        <Space size="large">
          <span>
            已选 {selectedItems.length} 件，合计：
            <span style={{ color: '#e4393c', fontWeight: 700, fontSize: 22 }}>{formatPrice(selectedTotal)}</span>
          </span>
          <Button type="primary" size="large" disabled={selectedItems.length === 0} onClick={() => navigate('/checkout')}>
            去结算
          </Button>
        </Space>
      </Card>
    </div>
  );
}
