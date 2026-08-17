import { useEffect, useState } from 'react';
import { Button, Space, Modal, Form, Input, InputNumber, Select, Switch, message, Popconfirm, Tag, Card, Typography } from 'antd';
import { PlusOutlined } from '@ant-design/icons';
import { adminListProducts, createProduct, updateProduct, changeProductStatus } from '../../api/admin';
import { listCategories } from '../../api/category';
import DataTable from '../../components/DataTable';
import StatusBadge from '../../components/StatusBadge';
import { usePagination } from '../../hooks/usePagination';
import { formatPrice } from '../../utils/format';
import type { ProductVO } from '../../api/product';
import type { CategoryVO } from '../../api/category';

// 商品管理：上下架、价格、库存、编辑。
export default function ProductManage() {
  const { page, pageSize, setPage, setPageSize, total, setTotal } = usePagination(10);
  const [products, setProducts] = useState<ProductVO[]>([]);
  const [loading, setLoading] = useState(false);
  const [categories, setCategories] = useState<CategoryVO[]>([]);
  const [modalOpen, setModalOpen] = useState(false);
  const [editing, setEditing] = useState<ProductVO | null>(null);
  const [form] = Form.useForm();

  const fetchData = async (p = page, ps = pageSize) => {
    setLoading(true);
    try {
      const data = await adminListProducts({ page: p, page_size: ps });
      setProducts(data.list);
      setTotal(data.total);
    } finally {
      setLoading(false);
    }
  };

  useEffect(() => {
    fetchData();
    listCategories().then(setCategories).catch(() => setCategories([]));
    // eslint-disable-next-line react-hooks/exhaustive-deps
  }, [pageSize]);

  const openCreate = () => {
    setEditing(null);
    form.resetFields();
    form.setFieldsValue({ status: 'ON_SALE', free_shipping: true, stock: 100, price: 99 });
    setModalOpen(true);
  };

  const openEdit = (p: ProductVO) => {
    setEditing(p);
    form.setFieldsValue({
      category_id: p.category_id,
      name: p.name,
      sub_title: p.sub_title,
      description: p.description,
      price: p.price,
      original_price: p.original_price,
      stock: p.stock,
      shipping_from: p.shipping_from,
      free_shipping: p.free_shipping,
      status: p.status,
    });
    setModalOpen(true);
  };

  const handleSave = async () => {
    const values = await form.validateFields();
    if (editing) {
      await updateProduct(editing.id, values);
      message.success('商品已更新');
    } else {
      await createProduct(values);
      message.success('商品已创建');
    }
    setModalOpen(false);
    fetchData();
  };

  const toggleStatus = async (p: ProductVO) => {
    const next = p.status === 'ON_SALE' ? 'OFF_SALE' : 'ON_SALE';
    await changeProductStatus(p.id, next);
    message.success(`已${next === 'ON_SALE' ? '上架' : '下架'}`);
    fetchData();
  };

  const columns = [
    {
      title: '商品',
      dataIndex: 'name',
      render: (v: string, r: ProductVO) => (
        <Space>
          <img src={r.cover_image} alt={v} style={{ width: 48, height: 48, objectFit: 'cover', borderRadius: 4 }} />
          <div>
            <div>{v}</div>
            <div style={{ color: '#999', fontSize: 12 }}>{r.sub_title}</div>
          </div>
        </Space>
      ),
    },
    { title: '价格', dataIndex: 'price', render: (v: number) => formatPrice(v) },
    { title: '库存', dataIndex: 'stock' },
    { title: '销量', dataIndex: 'sales' },
    { title: '状态', dataIndex: 'status', render: (v: string) => <StatusBadge status={v} kind="product" /> },
    {
      title: '操作',
      render: (_: unknown, r: ProductVO) => (
        <Space>
          <Button size="small" onClick={() => openEdit(r)}>
            编辑
          </Button>
          <Popconfirm title={`确定${r.status === 'ON_SALE' ? '下架' : '上架'}？`} onConfirm={() => toggleStatus(r)}>
            <Button size="small" danger={r.status === 'ON_SALE'}>
              {r.status === 'ON_SALE' ? '下架' : '上架'}
            </Button>
          </Popconfirm>
        </Space>
      ),
    },
  ];

  return (
    <div>
      <Card>
        <div style={{ display: 'flex', justifyContent: 'space-between', marginBottom: 16 }}>
          <Typography.Title level={4} style={{ margin: 0 }}>
            商品管理
          </Typography.Title>
          <Button type="primary" icon={<PlusOutlined />} onClick={openCreate}>
            新增商品
          </Button>
        </div>
        <DataTable
          rowKey="id"
          loading={loading}
          dataSource={products}
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

      <Modal title={editing ? '编辑商品' : '新增商品'} open={modalOpen} onCancel={() => setModalOpen(false)} onOk={handleSave} width={640}>
        <Form form={form} layout="vertical">
          <Form.Item name="category_id" label="分类" rules={[{ required: true, message: '请选择分类' }]}>
            <Select
              placeholder="选择分类"
              options={categories.filter((c) => c.level === 1).map((c) => ({ label: c.name, value: c.id }))}
            />
          </Form.Item>
          <Form.Item name="name" label="商品名称" rules={[{ required: true, message: '请输入名称' }]}>
            <Input />
          </Form.Item>
          <Form.Item name="sub_title" label="副标题">
            <Input />
          </Form.Item>
          <Form.Item name="description" label="描述">
            <Input.TextArea rows={2} />
          </Form.Item>
          <Space size="large" wrap>
            <Form.Item name="price" label="价格" rules={[{ required: true, message: '请输入价格' }]}>
              <InputNumber min={0.01} precision={2} style={{ width: 160 }} />
            </Form.Item>
            <Form.Item name="original_price" label="原价">
              <InputNumber min={0} precision={2} style={{ width: 160 }} />
            </Form.Item>
            <Form.Item name="stock" label="库存" rules={[{ required: true, message: '请输入库存' }]}>
              <InputNumber min={0} style={{ width: 160 }} />
            </Form.Item>
          </Space>
          <Space size="large" wrap>
            <Form.Item name="shipping_from" label="发货地">
              <Input style={{ width: 200 }} />
            </Form.Item>
            <Form.Item name="free_shipping" label="包邮" valuePropName="checked">
              <Switch />
            </Form.Item>
            <Form.Item name="status" label="状态">
              <Select
                options={[
                  { label: <Tag color="green">在售</Tag>, value: 'ON_SALE' },
                  { label: <Tag color="red">下架</Tag>, value: 'OFF_SALE' },
                ]}
              />
            </Form.Item>
          </Space>
        </Form>
      </Modal>
    </div>
  );
}
