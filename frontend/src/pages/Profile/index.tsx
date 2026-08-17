import { Card, Form, Input, Button, Descriptions, message, Divider, Space, Avatar } from 'antd';
import { UserOutlined } from '@ant-design/icons';
import { useUserStore } from '../../stores/userStore';
import { updateProfile, changePassword } from '../../api/auth';

// 个人中心：资料编辑 + 修改密码。
export default function Profile() {
  const { user, setUser } = useUserStore();
  const [profileForm] = Form.useForm();
  const [pwdForm] = Form.useForm();

  const saveProfile = async (values: { nickname?: string; email?: string; phone?: string }) => {
    const updated = await updateProfile(values);
    setUser(updated);
    message.success('资料已更新');
  };

  const savePassword = async (values: { old_password: string; new_password: string }) => {
    await changePassword(values);
    message.success('密码已修改');
    pwdForm.resetFields();
  };

  return (
    <div style={{ maxWidth: 720, margin: '0 auto' }}>
      <Card style={{ marginBottom: 16 }}>
        <Space size="large">
          <Avatar size={64} icon={<UserOutlined />} />
          <div>
            <div style={{ fontSize: 20, fontWeight: 700 }}>{user?.nickname || user?.username}</div>
            <div style={{ color: '#999' }}>@{user?.username} · {user?.role === 'ADMIN' ? '管理员' : '普通用户'}</div>
          </div>
        </Space>
      </Card>

      <Card title="基本资料" style={{ marginBottom: 16 }}>
        <Descriptions column={2} style={{ marginBottom: 16 }}>
          <Descriptions.Item label="用户名">{user?.username}</Descriptions.Item>
          <Descriptions.Item label="角色">{user?.role}</Descriptions.Item>
          <Descriptions.Item label="邮箱">{user?.email || '-'}</Descriptions.Item>
          <Descriptions.Item label="手机号">{user?.phone || '-'}</Descriptions.Item>
          <Descriptions.Item label="注册时间">{user?.created_at}</Descriptions.Item>
        </Descriptions>
        <Divider />
        <Form
          form={profileForm}
          layout="vertical"
          onFinish={saveProfile}
          initialValues={{ nickname: user?.nickname, email: user?.email, phone: user?.phone }}
        >
          <Form.Item name="nickname" label="昵称">
            <Input />
          </Form.Item>
          <Form.Item name="email" label="邮箱" rules={[{ type: 'email', message: '邮箱格式不正确' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="phone" label="手机号">
            <Input />
          </Form.Item>
          <Button type="primary" htmlType="submit">
            保存资料
          </Button>
        </Form>
      </Card>

      <Card title="修改密码">
        <Form form={pwdForm} layout="vertical" onFinish={savePassword}>
          <Form.Item name="old_password" label="原密码" rules={[{ required: true, min: 6, message: '请输入原密码' }]}>
            <Input.Password />
          </Form.Item>
          <Form.Item name="new_password" label="新密码" rules={[{ required: true, min: 6, message: '至少 6 位' }]}>
            <Input.Password />
          </Form.Item>
          <Button type="primary" htmlType="submit">
            修改密码
          </Button>
        </Form>
      </Card>
    </div>
  );
}
