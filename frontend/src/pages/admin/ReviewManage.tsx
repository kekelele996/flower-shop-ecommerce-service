import { useEffect, useState } from 'react';
import { Button, Space, Modal, Input, message, Rate, Card, Typography } from 'antd';
import { adminListReviews } from '../../api/admin';
import { replyReview } from '../../api/review';
import DataTable from '../../components/DataTable';
import { usePagination } from '../../hooks/usePagination';
import type { ReviewVO } from '../../api/review';

// 评价管理：查看并回复评价。
export default function ReviewManage() {
  const { page, pageSize, setPage, setPageSize, total, setTotal } = usePagination(10);
  const [reviews, setReviews] = useState<ReviewVO[]>([]);
  const [loading, setLoading] = useState(false);
  const [replyTarget, setReplyTarget] = useState<ReviewVO | null>(null);
  const [reply, setReply] = useState('');

  const fetchData = async (p = page, ps = pageSize) => {
    setLoading(true);
    try {
      const data = await adminListReviews({ page: p, page_size: ps });
      setReviews(data.list);
      setTotal(data.total);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pageSize]);

  const handleReply = async () => {
    if (!replyTarget || !reply.trim()) return;
    await replyReview(replyTarget.id, reply.trim());
    message.success('回复成功');
    setReplyTarget(null);
    setReply('');
    fetchData();
  };

  const columns = [
    { title: '用户', dataIndex: 'username' },
    { title: '商品', dataIndex: 'product_name' },
    {
      title: '评分',
      dataIndex: 'rating',
      render: (v: number) => <Rate disabled value={v} style={{ fontSize: 14 }} />,
    },
    { title: '内容', dataIndex: 'content' },
    { title: '商家回复', dataIndex: 'reply', render: (v: string) => v || '-' },
    { title: '时间', dataIndex: 'created_at' },
    {
      title: '操作',
      render: (_: unknown, r: ReviewVO) => (
        <Button
          size="small"
          type="primary"
          disabled={!!r.reply}
          onClick={() => {
            setReplyTarget(r);
            setReply('');
          }}
        >
          {r.reply ? '已回复' : '回复'}
        </Button>
      ),
    },
  ];

  return (
    <div>
      <Card>
        <Typography.Title level={4} style={{ marginBottom: 16 }}>
          评价管理
        </Typography.Title>
        <DataTable
          rowKey="id"
          loading={loading}
          dataSource={reviews}
          columns={columns}
          total={total}
          page={page}
          pageSize={pageSize}
          onPageChange={(p, ps) => {
            setPage(p);
            setPageSize(ps);
            fetchData(p, ps);
          }}
        />
      </Card>

      <Modal
        title={`回复 ${replyTarget?.username || ''} 的评价`}
        open={!!replyTarget}
        onCancel={() => setReplyTarget(null)}
        onOk={handleReply}
        okText="提交回复"
      >
        <Input.TextArea rows={4} value={reply} onChange={(e) => setReply(e.target.value)} placeholder="输入回复内容" />
      </Modal>
    </div>
  );
}
