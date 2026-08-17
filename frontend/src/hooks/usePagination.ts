import { useCallback, useState } from 'react';

// usePagination：列表分页状态管理（跨页面复用）。
export function usePagination(defaultPageSize = 10) {
  const [page, setPage] = useState(1);
  const [pageSize, setPageSize] = useState(defaultPageSize);
  const [total, setTotal] = useState(0);

  const reset = useCallback(() => setPage(1), []);

  return {
    page,
    pageSize,
    total,
    setPage,
    setPageSize,
    setTotal,
    reset,
  };
}
