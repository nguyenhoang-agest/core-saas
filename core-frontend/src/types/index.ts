export interface ApiResponse<T = unknown> {
  data: T;
  message?: string;
  status: number;
  success: boolean;
}

export interface UserProfile {
  id: string;
  email: string;
  name: string;
  image?: string;
}
