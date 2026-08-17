import { useEffect, useState } from 'react';
import { Card, Typography, Tag, Space } from 'antd';
import { listAuditLogs } from '../../api/admin';
import DataTable from '../../components/DataTable';
import { usePagination } from '../../hooks/usePagination';
import type { AuditLogVO } from '../../api/admin';

const ACTION_COLOR: Record<string, string> = {
  CREATE: 'green',
  UPDATE: 'blue',
  DELETE: 'red',
  LOGIN: 'purple',
};

// 审计日志：管理员查看全站操作记录（横切关注点）。
export default function AuditLogs() {
  const { page, pageSize, setPage, setPageSize, total, setTotal } = usePagination(10);
  const [logs, setLogs] = useState<AuditLogVO[]>([]);
  const [loading, setLoading] = useState(false);

  const fetchData = async (p = page, ps = pageSize) => {
    setLoading(true);
    try {
      const data = await listAuditLogs({ page: p, page_size: ps });
      setLogs(data.list);
      setTotal(data.total);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pageSize]);

  const columns = [
    { title: 'ID', dataIndex: 'id', width: 70 },
    { title: '用户', dataIndex: 'username', render: (v: string) => v || '-' },
    {
      title: '动作',
      dataIndex: 'action',
      render: (v: string) => <Tag color={ACTION_COLOR[v] || 'default'}>{v}</Tag>,
    },
    { title: '方法', dataIndex: 'method', render: (v: string) => <Tag>{v}</Tag> },
    { title: '实体', dataIndex: 'entity' },
    { title: '实体ID', dataIndex: 'entity_id', render: (v: string) => v || '-' },
    { title: '路径', dataIndex: 'path' },
    { title: '详情', dataIndex: 'detail', ellipsis: true },
    { title: 'IP', dataIndex: 'ip' },
    { title: '时间', dataIndex: 'created_at' },
  ];

  return (
    <Card>
      <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
        <Typography.Title level={4} style={{ margin: 0 }}>
          审计日志
        </Typography.Title>
        <Space>
          <Typography.Text type="secondary">操作审计：写入 audit_logs 表，覆盖写操作</Typography.Text>
        </Space>
      </div>
      <DataTable
        rowKey="id"
        loading={loading}
        dataSource={logs}
        columns={columns}
        total={total}
        page={page}
        pageSize={pageSize}
        onPageChange={(p, ps) => {
          setPage(p);
          setPageSize(ps);
          fetchData(p, ps);
        }}
        scroll={{ x: 1200 }}
      />
    </Card>
  );
}
