import { Layout, Menu, Button } from 'antd';
import { DashboardOutlined, ShoppingOutlined, FileTextOutlined, CommentOutlined, GiftOutlined, SafetyCertificateOutlined, RollbackOutlined } from '@ant-design/icons';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';

const { Sider, Content } = Layout;

const items = [
  { key: '/admin', icon: <DashboardOutlined />, label: '数据看板' },
  { key: '/admin/products', icon: <ShoppingOutlined />, label: '商品管理' },
  { key: '/admin/orders', icon: <FileTextOutlined />, label: '订单管理' },
  { key: '/admin/reviews', icon: <CommentOutlined />, label: '评价管理' },
  { key: '/admin/coupons', icon: <GiftOutlined />, label: '优惠券管理' },
  { key: '/admin/audit-logs', icon: <SafetyCertificateOutlined />, label: '审计日志' },
];

export default function AdminLayout() {
  const navigate = useNavigate();
  const location = useLocation();

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Sider theme="dark">
        <div style={{ color: '#fff', fontSize: 18, fontWeight: 700, padding: 16, textAlign: 'center' }}>
          🌸 花语商城后台
        </div>
        <Menu theme="dark" mode="inline" selectedKeys={[location.pathname]} items={items} onClick={({ key }) => navigate(key)} />
        <div style={{ padding: 16, textAlign: 'center' }}>
          <Button icon={<RollbackOutlined />} onClick={() => navigate('/')} ghost size="small">
            返回前台
          </Button>
        </div>
      </Sider>
      <Content style={{ padding: 24, background: '#f5f5f5' }}>
        <Outlet />
      </Content>
    </Layout>
  );
}
