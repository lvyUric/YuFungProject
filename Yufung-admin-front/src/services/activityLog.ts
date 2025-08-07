import { request } from '@umijs/max';

export interface ActivityLog {
  id?: string;
  log_id?: string;
  user_id: string;
  username: string;
  company_id: string;
  company_name: string;
  operation_type: string;
  module_name: string;
  operation_desc: string;
  request_url: string;
  request_method: string;
  request_params: any;
  ip_address: string;
  user_agent: string;
  operation_time: string;
  execution_time: number;
  result_status: string;
  target_id?: string;
  target_name?: string;
}

export interface ActivityLogQuery {
  page?: number;
  page_size?: number;
  company_id?: string;
  user_id?: string;
  operation_type?: string;
  module_name?: string;
  result_status?: string;
  start_time?: string;
  end_time?: string;
  sort_by?: string;
  sort_order?: string;
}

export interface ActivityLogResponse {
  code: number;
  message: string;
  data: {
    total: number;
    list: ActivityLog[];
  };
}

export interface ActivityLogStatistics {
  total: number;
  operation_types: Array<{ _id: string; count: number }>;
  modules: Array<{ _id: string; count: number }>;
}

/**
 * 获取活动记录列表
 */
export async function getActivityLogs(params: ActivityLogQuery): Promise<ActivityLogResponse> {
  return request('/api/v1/activity-logs', {
    method: 'GET',
    params,
  });
}

/**
 * 获取最近的活动记录（用于仪表盘）
 */
export async function getRecentActivityLogs(limit: number = 5): Promise<ActivityLog[]> {
  return request('/api/v1/activity-logs/recent', {
    method: 'GET',
    params: { limit },
  });
}

/**
 * 获取活动记录详情
 */
export async function getActivityLogById(id: string): Promise<ActivityLog> {
  return request(`/api/v1/activity-logs/${id}`, {
    method: 'GET',
  });
}

/**
 * 获取活动记录统计
 */
export async function getActivityLogStatistics(days: number = 7): Promise<ActivityLogStatistics> {
  return request('/api/v1/activity-logs/statistics', {
    method: 'GET',
    params: { days },
  });
}

/**
 * 删除指定公司的活动记录（仅平台管理员）
 */
export async function deleteActivityLogsByCompany(companyId: string): Promise<void> {
  return request(`/api/v1/activity-logs/company/${companyId}`, {
    method: 'DELETE',
  });
}

// 操作类型常量
export const OPERATION_TYPES = {
  CREATE: 'create',
  UPDATE: 'update',
  DELETE: 'delete',
  VIEW: 'view',
  EXPORT: 'export',
  IMPORT: 'import',
  LOGIN: 'login',
  LOGOUT: 'logout',
} as const;

// 模块名称常量
export const MODULE_NAMES = {
  USER: '用户管理',
  ROLE: '角色管理',
  MENU: '菜单管理',
  COMPANY: '公司管理',
  POLICY: '保单管理',
  CUSTOMER: '客户管理',
  SYSTEM: '系统管理',
  AUTH: '认证授权',
} as const;

// 操作类型显示文本
export const OPERATION_TYPE_LABELS = {
  [OPERATION_TYPES.CREATE]: '新增',
  [OPERATION_TYPES.UPDATE]: '更新',
  [OPERATION_TYPES.DELETE]: '删除',
  [OPERATION_TYPES.VIEW]: '查看',
  [OPERATION_TYPES.EXPORT]: '导出',
  [OPERATION_TYPES.IMPORT]: '导入',
  [OPERATION_TYPES.LOGIN]: '登录',
  [OPERATION_TYPES.LOGOUT]: '登出',
} as const;

// 结果状态显示文本
export const RESULT_STATUS_LABELS = {
  success: '成功',
  failure: '失败',
} as const; 