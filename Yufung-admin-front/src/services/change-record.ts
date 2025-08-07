import { request } from '@umijs/max';

// 变更记录项目接口
export interface ChangeRecordItem {
  id: string;
  change_id: string;
  table_name: string;
  record_id: string;
  user_id: string;
  username: string;
  company_id: string;
  change_type: 'insert' | 'update' | 'delete';
  old_values?: Record<string, any>;
  new_values?: Record<string, any>;
  changed_fields: string[];
  change_time: string;
  change_reason?: string;
  ip_address?: string;
  user_agent?: string;
  change_time_formatted: string;
  change_details: ChangeDetail[];
}

// 变更详情接口
export interface ChangeDetail {
  field_name: string;
  field_label: string;
  old_value: any;
  new_value: any;
  old_value_text: string;
  new_value_text: string;
}

// 变更记录查询参数
export interface ChangeRecordParams {
  days?: number;
  page?: number;
  page_size?: number;
}

// 变更记录列表参数
export interface ChangeRecordListParams {
  table_name?: string;
  record_id?: string;
  user_id?: string;
  change_type?: string;
  start_time?: string;
  end_time?: string;
  page?: number;
  page_size?: number;
}

// API响应接口 - 修正为与后端一致的格式
export interface ChangeRecordResponse {
  code: number;
  message: string;
  data: {
    records: ChangeRecordItem[];
    total: number;
    page: number;
    page_size: number;
    has_more: boolean;
  };
}

/**
 * 获取保单变更记录
 * @param policyId 保单ID
 * @param params 查询参数
 */
export async function getPolicyChangeRecords(
  policyId: string,
  params?: ChangeRecordParams,
): Promise<ChangeRecordResponse> {
  return request(`/api/policies/${policyId}/change-records`, {
    method: 'GET',
    params: {
      days: 10, // 默认查询10天
      page: 1,
      page_size: 10,
      ...params,
    },
  });
}

/**
 * 获取变更记录列表
 * @param params 查询参数
 */
export async function getChangeRecordsList(
  params?: ChangeRecordListParams,
): Promise<ChangeRecordResponse> {
  return request('/api/change-records', {
    method: 'GET',
    params: {
      page: 1,
      page_size: 10,
      ...params,
    },
  });
} 