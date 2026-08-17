import { useState } from 'react';
import { Button, Form, Input, InputNumber, Select, message, Card, Typography, Space, Tag } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { createCouponTemplate } from '../../api/coupon';

// 优惠券管理：创建满减/折扣券模板。
export default function CouponManage() {
  const [form] = Form.useForm();
  const [loading, setLoading] = useState(false);
  const [type, setType] = useState('FULL_REDUCTION');

  const handleCreate = async (values: { name: string; type: string; threshold: number; amount: number; discount_rate: number; total: number; valid_days: number }) => {
    setLoading(true);
    try {
      const res = await createCouponTemplate(values);
      message.success(`优惠券模板已创建，共发放 ${res.created} 张`);
      form.resetFields();
    } finally {
      setLoading(false);
    }
  };

  return (
    <Card style={{ maxWidth: 640 }}>
      <Typography.Title level={4}>优惠券管理</Typography.Title>
      <Typography.Paragraph type="secondary">
        创建满减券（满 X 减 Y）或折扣券（X 折），系统自动发放给券池。
      </Typography.Paragraph>
      <Form form={form} layout="vertical" onFinish={handleCreate} initialValues={{ type: 'FULL_REDUCTION', threshold: 100, amount: 20, total: 100, valid_days: 30, discount_rate: 0.9 }}>
        <Form.Item name="name" label="券名称" rules={[{ required: true, message: '请输入名称' }]}>
          <Input />
        </Form.Item>
        <Form.Item name="type" label="券类型" rules={[{ required: true }]}>
          <Select
            onChange={setType}
            options={[
              { label: <Tag color="orange">满减券</Tag>, value: 'FULL_REDUCTION' },
              { label: <Tag color="blue">折扣券</Tag>, value: 'DISCOUNT' },
            ]}
          />
        </Form.Item>
        <Space size="large" wrap>
          <Form.Item name="threshold" label="使用门槛（满 X 元）" rules={[{ required: true }]}>
            <InputNumber min={0} precision={2} style={{ width: 180 }} />
          </Form.Item>
          {type === 'FULL_REDUCTION' ? (
            <Form.Item name="amount" label="减免金额" rules={[{ required: true }]}>
              <InputNumber min={0.01} precision={2} style={{ width: 180 }} />
            </Form.Item>
          ) : (
            <Form.Item name="discount_rate" label="折扣率（0.9=9折）" rules={[{ required: true }]}>
              <InputNumber min={0.01} max={1} step={0.05} style={{ width: 180 }} />
            </Form.Item>
          )}
        </Space>
        <Space size="large" wrap>
          <Form.Item name="total" label="发放数量" rules={[{ required: true }]}>
            <InputNumber min={1} style={{ width: 180 }} />
          </Form.Item>
          <Form.Item name="valid_days" label="有效天数" rules={[{ required: true }]}>
            <InputNumber min={1} style={{ width: 180 }} />
          </Form.Item>
        </Space>
        <Button type="primary" htmlType="submit" icon={<PlusOutlined />} loading={loading}>
          创建并发放
        </Button>
      </Form>
    </Card>
  );
}
