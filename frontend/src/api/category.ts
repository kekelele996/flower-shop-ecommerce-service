import { get } from '../utils/request';

export interface CategoryVO {
  id: number;
  parent_id: number;
  name: string;
  level: number;
  sort: number;
  children?: CategoryVO[];
}

export const listCategories = () => get<CategoryVO[]>('/categories');
