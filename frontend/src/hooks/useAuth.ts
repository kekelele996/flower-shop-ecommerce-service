import { useNavigate } from 'react-router-dom';
import { message } from 'antd';
import { login as apiLogin, register as apiRegister } from '../api/auth';
import { useUserStore } from '../stores/userStore';

// useAuth：登录/注册/退出登录，路由守卫与按钮显隐共用。
export function useAuth() {
  const navigate = useNavigate();
  const { setToken, setUser, logout } = useUserStore();

  const login = async (username: string, password: string) => {
    const token = await apiLogin({ username, password });
    setToken(token.access_token);
    setUser(token.user);
    message.success('登录成功');
    navigate(token.user.role === 'ADMIN' ? '/admin' : '/');
    return token;
  };

  const register = async (data: { username: string; password: string; nickname?: string; email?: string }) => {
    await apiRegister(data);
    message.success('注册成功，请登录');
    navigate('/login');
  };

  const logoutFn = () => {
    logout();
    message.success('已退出登录');
    navigate('/login');
  };

  return { login, register, logout: logoutFn };
}
