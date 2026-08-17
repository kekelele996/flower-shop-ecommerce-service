import { Card, Form, Input, Button, Typography, Space } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { Link } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

// 登录页。
export default function Login() {
  const { login } = useAuth();
  return (
    <div style={{ maxWidth: 420, margin: '60px auto' }}>
      <Card>
        <Typography.Title level={3} style={{ textAlign: 'center' }}>
          登录花语商城
        </Typography.Title>
        <Typography.Paragraph type="secondary" style={{ textAlign: 'center' }}>
          测试账号：user / user123（管理员 admin / admin123）
        </Typography.Paragraph>
        <Form onFinish={(values) => login(values.username, values.password)} initialValues={{ username: 'user', password: 'user123' }}>
          <Form.Item name="username" rules={[{ required: true, message: '请输入用户名' }]}>
            <Input prefix={<UserOutlined />} placeholder="用户名" size="large" />
          </Form.Item>
          <Form.Item name="password" rules={[{ required: true, message: '请输入密码' }]}>
            <Input.Password prefix={<LockOutlined />} placeholder="密码" size="large" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block size="large">
              登录
            </Button>
          </Form.Item>
          <Space style={{ justifyContent: 'center', width: '100%' }}>
            <Link to="/register">没有账号？立即注册</Link>
          </Space>
        </Form>
      </Card>
    </div>
  );
}
