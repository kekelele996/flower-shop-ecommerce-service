import { Card, Form, Input, Button, Typography } from 'antd';
import { Link } from 'react-router-dom';
import { useAuth } from '../../hooks/useAuth';

// 注册页。
export default function Register() {
  const { register } = useAuth();
  return (
    <div style={{ maxWidth: 420, margin: '60px auto' }}>
      <Card>
        <Typography.Title level={3} style={{ textAlign: 'center' }}>
          注册花语商城
        </Typography.Title>
        <Form onFinish={(values) => register(values)}>
          <Form.Item name="username" rules={[{ required: true, min: 3, message: '用户名至少 3 位' }]}>
            <Input placeholder="用户名" size="large" />
          </Form.Item>
          <Form.Item name="password" rules={[{ required: true, min: 6, message: '密码至少 6 位' }]}>
            <Input.Password placeholder="密码" size="large" />
          </Form.Item>
          <Form.Item name="nickname">
            <Input placeholder="昵称（可选）" size="large" />
          </Form.Item>
          <Form.Item name="email" rules={[{ type: 'email', message: '邮箱格式不正确' }]}>
            <Input placeholder="邮箱（可选）" size="large" />
          </Form.Item>
          <Form.Item>
            <Button type="primary" htmlType="submit" block size="large">
              注册
            </Button>
          </Form.Item>
          <div style={{ textAlign: 'center' }}>
            <Link to="/login">已有账号？去登录</Link>
          </div>
        </Form>
      </Card>
    </div>
  );
}
