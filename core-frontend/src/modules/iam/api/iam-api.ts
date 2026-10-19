import { apiClient } from '@/apis/axios-client';

export interface IamUser {
  id: string;
  email: string;
  name: string;
  role: string;
}

export async function fetchIamUsers(tenantId: string): Promise<IamUser[]> {
  try {
    const res = await apiClient.get<IamUser[]>('/iam/users');
    return res.data;
  } catch (_e) {
    return [];
  }
}
