import { apiClient } from '@/apis/axios-client';

export interface TemplateData {
  id: string;
  name: string;
}

export async function fetchTemplateData(tenantId: string): Promise<TemplateData[]> {
  try {
    const response = await apiClient.get<TemplateData[]>('/__MODULE_NAME__');
    return response.data;
  } catch (_error) {
    return [];
  }
}
