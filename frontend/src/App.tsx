import { BrowserRouter, Routes, Route } from 'react-router-dom';
import MainLayout from './layouts/MainLayout';
import AdminLayout from './layouts/AdminLayout';
import AuthGuard from './components/AuthGuard';
import Home from './pages/Home';
import ProductDetail from './pages/ProductDetail';
import Login from './pages/Login';
import Register from './pages/Register';
import Cart from './pages/Cart';
import Checkout from './pages/Checkout';
import Orders from './pages/Orders';
import OrderDetail from './pages/OrderDetail';
import Profile from './pages/Profile';
import Coupons from './pages/Coupons';
import AdminDashboard from './pages/admin/Dashboard';
import AdminProducts from './pages/admin/ProductManage';
import AdminOrders from './pages/admin/OrderManage';
import AdminReviews from './pages/admin/ReviewManage';
import AdminCoupons from './pages/admin/CouponManage';
import AdminAuditLogs from './pages/admin/AuditLogs';

// 路由配置：前台 + 管理后台。
export default function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<MainLayout />}>
          <Route index element={<Home />} />
          <Route path="products/:id" element={<ProductDetail />} />
          <Route path="login" element={<Login />} />
          <Route path="register" element={<Register />} />
          <Route
            path="cart"
            element={
              <AuthGuard>
                <Cart />
              </AuthGuard>
            }
          />
          <Route
            path="checkout"
            element={
              <AuthGuard>
                <Checkout />
              </AuthGuard>
            }
          />
          <Route
            path="orders"
            element={
              <AuthGuard>
                <Orders />
              </AuthGuard>
            }
          />
          <Route
            path="orders/:id"
            element={
              <AuthGuard>
                <OrderDetail />
              </AuthGuard>
            }
          />
          <Route
            path="profile"
            element={
              <AuthGuard>
                <Profile />
              </AuthGuard>
            }
          />
          <Route
            path="coupons"
            element={
              <AuthGuard>
                <Coupons />
              </AuthGuard>
            }
          />
        </Route>
        <Route
          path="/admin"
          element={
            <AuthGuard adminOnly>
              <AdminLayout />
            </AuthGuard>
          }
        >
          <Route index element={<AdminDashboard />} />
          <Route path="products" element={<AdminProducts />} />
          <Route path="orders" element={<AdminOrders />} />
          <Route path="reviews" element={<AdminReviews />} />
          <Route path="coupons" element={<AdminCoupons />} />
          <Route path="audit-logs" element={<AdminAuditLogs />} />
        </Route>
      </Routes>
    </BrowserRouter>
  );
}
