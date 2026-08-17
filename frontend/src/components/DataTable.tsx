import { Table } from 'antd';
import type { TableProps } from 'antd';

interface DataTableProps<T> extends Omit<TableProps<T>, 'pagination'> {
  total: number;
  page: number;
  pageSize: number;
  onPageChange: (page: number, pageSize: number) => void;
  rowKey?: string;
}

// DataTable：通用分页表格（跨页面复用）。
export default function DataTable<T extends object>({ total, page, pageSize, onPageChange, rowKey = 'id', ...rest }: DataTableProps<T>) {
  return (
    <Table<T>
      rowKey={rowKey}
      pagination={{
        current: page,
        pageSize,
        total,
        showSizeChanger: true,
        showTotal: (t) => `共 ${t} 条`,
        onChange: onPageChange,
      }}
      {...rest}
    />
  );
}
