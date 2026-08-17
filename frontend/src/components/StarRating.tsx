import { Rate } from 'antd';

interface StarRatingProps {
  value: number;
  disabled?: boolean;
  onChange?: (value: number) => void;
}

// StarRating：星级评分（评价展示/提交复用）。
export default function StarRating({ value, disabled = true, onChange }: StarRatingProps) {
  return <Rate disabled={disabled} value={value} onChange={onChange} allowHalf={false} />;
}
