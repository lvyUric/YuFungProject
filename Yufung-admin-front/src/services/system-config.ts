import { request } from '@umijs/max';

// 系统配置信息接口
export interface SystemConfigInfo {
  id: string;
  config_id: string;
  config_type: string;
  config_key: string;
  config_value: string;
  display_name: string;
  company_id: string;
  sort_order: number;
  status: string;
  status_text: string;
  remark: string;
  created_by: string;
  updated_by: string;
  created_at: string;
  updated_at: string;
}

// 系统配置列表查询参数
export interface SystemConfigListParams {
  current?: number;
  pageSize?: number;
  page?: number;
  page_size?: number;
  config_type?: string;
  config_key?: string;
  status?: string;
  keyword?: string;
}

// 创建系统配置请求
export interface SystemConfigCreateRequest {
  config_type: string;
  config_key: string;
  config_value: string;
  display_name: string;
  sort_order?: number;
  status?: string;
  remark?: string;
}

// 更新系统配置请求
export interface SystemConfigUpdateRequest {
  config_value?: string;
  display_name?: string;
  sort_order?: number;
  status?: string;
  remark?: string;
}

// 配置类型选项
export const configTypeOptions = [
  { label: '港分客户经理', value: 'hk_manager' },
  { label: '转介分行', value: 'referral_branch' },
  { label: '合作伙伴', value: 'partner' },
];

// 状态选项
export const statusOptions = [
  { label: '启用', value: 'enable' },
  { label: '禁用', value: 'disable' },
];

// 获取系统配置列表
export async function getSystemConfigList(params?: SystemConfigListParams): Promise<API.Response<{ list: SystemConfigInfo[]; total: number }>> {
  return request('/api/system-configs', {
    method: 'GET',
    params,
  });
}

// 创建系统配置
export async function createSystemConfig(data: SystemConfigCreateRequest): Promise<API.Response<SystemConfigInfo>> {
  return request('/api/system-configs', {
    method: 'POST',
    data,
  });
}

// 获取系统配置详情
export async function getSystemConfig(configId: string): Promise<API.Response<SystemConfigInfo>> {
  return request(`/api/system-configs/${configId}`, {
    method: 'GET',
  });
}

// 更新系统配置
export async function updateSystemConfig(configId: string, data: SystemConfigUpdateRequest): Promise<API.Response<SystemConfigInfo>> {
  return request(`/api/system-configs/${configId}`, {
    method: 'PUT',
    data,
  });
}

// 删除系统配置
export async function deleteSystemConfig(configId: string): Promise<API.Response<any>> {
  return request(`/api/system-configs/${configId}`, {
    method: 'DELETE',
  });
}

// 根据配置类型获取配置选项
export async function getSystemConfigOptions(configType: string): Promise<API.Response<SystemConfigInfo[]>> {
  return request(`/api/system-configs/options/${configType}`, {
    method: 'GET',
  });
} 