import { Modal } from 'antd';
import { ExclamationCircleOutlined } from '@ant-design/icons';

interface ConfirmDialogOptions {
  title: string;
  content: string;
  onOk: () => void | Promise<void>;
  okText?: string;
}

// ConfirmDialog：确认弹窗（跨页面复用）。
export function showConfirm({ title, content, onOk, okText = '确定' }: ConfirmDialogOptions) {
  Modal.confirm({
    title,
    icon: <ExclamationCircleOutlined />,
    content,
    okText,
    cancelText: '取消',
    onOk,
  });
}
