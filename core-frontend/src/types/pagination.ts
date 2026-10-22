export interface PaginationParams {
  page?: number;
  pageSize?: number;
  page_size?: number;
  sortBy?: string;
  sort_by?: string;
  order?: 'asc' | 'desc';
}

export interface PaginationMeta {
  page: number;
  pageSize: number;
  totalItems: number;
  totalPages: number;
  page_size?: number;
  total_items?: number;
  total_pages?: number;
}

export interface PaginatedData<T> {
  items: T[];
  pagination: PaginationMeta;
}
