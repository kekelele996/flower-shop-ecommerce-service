import { Layout, Badge, Input, Button, Space, Dropdown } from 'antd';
import { ShoppingCartOutlined, UserOutlined, HomeOutlined, LogoutOutlined, GiftOutlined, FileTextOutlined } from '@ant-design/icons';
import { Outlet, useNavigate, useLocation } from 'react-router-dom';
import { useEffect, useMemo, useState } from 'react';
import { useUserStore } from '../stores/userStore';
import { useCartStore } from '../stores/cartStore';
import { useAuth } from '../hooks/useAuth';
import { listCategories } from '../api/category';
import type { CategoryVO } from '../api/category';

const { Header, Content, Footer } = Layout;

export default function MainLayout() {
  const navigate = useNavigate();
  const location = useLocation();
  const { user, isLogin, isAdmin } = useUserStore();
  const { logout } = useAuth();
  const { summary, fetchCart } = useCartStore();
  const [categories, setCategories] = useState<CategoryVO[]>([]);
  const [keyword, setKeyword] = useState('');

  useEffect(() => {
    listCategories().then(setCategories).catch(() => setCategories([]));
    if (isLogin()) {
      fetchCart();
    }
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [location.pathname]);

  const topCategories = useMemo(() => categories.filter((c) => c.level === 1).slice(0, 6), [categories]);

  const userMenu = {
    items: [
      { key: 'profile', icon: <UserOutlined />, label: '个人中心' },
      { key: 'orders', icon: <FileTextOutlined />, label: '我的订单' },
      { key: 'coupons', icon: <GiftOutlined />, label: '我的优惠券' },
      ...(isAdmin() ? [{ key: 'admin', icon: <HomeOutlined />, label: '管理后台' }] : []),
      { type: 'divider' as const },
      { key: 'logout', icon: <LogoutOutlined />, label: '退出登录' },
    ],
    onClick: ({ key }: { key: string }) => {
      if (key === 'logout') {
        logout();
      } else if (key === 'admin') {
        navigate('/admin');
      } else if (key === 'orders') {
        navigate('/orders');
      } else if (key === 'coupons') {
        navigate('/coupons');
      } else {
        navigate('/profile');
      }
    },
  };

  return (
    <Layout style={{ minHeight: '100vh' }}>
      <Header style={{ background: '#2f5d3a', display: 'flex', alignItems: 'center', gap: 16, padding: '0 24px' }}>
        <div
          style={{ color: '#fff', fontSize: 20, fontWeight: 700, cursor: 'pointer', whiteSpace: 'nowrap' }}
          onClick={() => navigate('/')}
        >
          🌸 花语商城
        </div>
        <Input.Search
          placeholder="搜索花卉/绿植/花盆/工具"
          allowClear
          style={{ maxWidth: 420, flex: 1 }}
          value={keyword}
          onChange={(e) => setKeyword(e.target.value)}
          onSearch={(value) => {
            navigate(`/?keyword=${encodeURIComponent(value)}`);
          }}
        />
        <div style={{ flex: 1 }} />
        <Space size="middle">
          {isLogin() ? (
            <Dropdown menu={userMenu}>
              <span style={{ color: '#fff', cursor: 'pointer' }}>
                <UserOutlined /> {user?.nickname || user?.username}
              </span>
            </Dropdown>
          ) : (
            <>
              <Button type="primary" ghost onClick={() => navigate('/login')}>
                登录
              </Button>
              <Button onClick={() => navigate('/register')}>注册</Button>
            </>
          )}
          <Badge count={summary?.total_quantity || 0} size="small">
            <Button icon={<ShoppingCartOutlined />} onClick={() => navigate(isLogin() ? '/cart' : '/login')}>
              购物车
            </Button>
          </Badge>
        </Space>
      </Header>
      <div style={{ background: '#fff', borderBottom: '1px solid #f0f0f0', padding: '8px 24px' }}>
        <Space size="large" wrap>
          {topCategories.map((c) => (
            <a key={c.id} onClick={() => navigate(`/?category_id=${c.id}`)} style={{ color: '#333', fontSize: 14 }}>
              {c.name}
            </a>
          ))}
        </Space>
      </div>
      <Content style={{ padding: '24px 24px 0', maxWidth: 1200, width: '100%', margin: '0 auto' }}>
        <Outlet />
      </Content>
      <Footer style={{ textAlign: 'center', color: '#999' }}>
        花语商城 · 花卉B2C电商平台 · Go 1.22 + Gin + GORM / React 18 + Ant Design
      </Footer>
    </Layout>
  );
}
