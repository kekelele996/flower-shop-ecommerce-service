import { Empty } from 'antd';
import type { ReactNode } from 'react';

interface EmptyStateProps {
  description?: string;
  children?: ReactNode;
}

// EmptyState：空状态组件（跨页面复用）。
export default function EmptyState({ description = '暂无数据', children }: EmptyStateProps) {
  return (
    <div style={{ padding: 48 }}>
      <Empty description={description}>
        {children}
      </Empty>
    </div>
  );
}
