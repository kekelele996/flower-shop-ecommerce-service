import { useCallback, useState } from 'react';

// useRequest：统一 loading 状态封装。
export function useRequest<T extends unknown[]>(fn: (...args: T) => Promise<void>) {
  const [loading, setLoading] = useState(false);
  const run = useCallback(
    async (...args: T) => {
      setLoading(true);
      try {
        await fn(...args);
      } finally {
        setLoading(false);
      }
    },
    [fn],
  );
  return { loading, run };
}
