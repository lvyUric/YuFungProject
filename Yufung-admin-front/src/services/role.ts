import { request } from '@umijs/max';

export interface RoleInfo {
  id: string;
  role_id: string;
  role_name: string;
  role_key: string;
  company_id: string;
  sort_order: number;
  data_scope: string;
  menu_ids: string[];
  status: string;
  remark: string;
  created_at: string;
  updated_at: string;
}

export interface RoleListParams {
  page?: number;
  page_size?: number;
  role_name?: string;
  role_key?: string;
  company_id?: string;
  data_scope?: string;
  status?: string;
}

export interface RoleListResponse {
  roles: RoleInfo[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface RoleCreateRequest {
  role_name: string;
  role_key: string;
  company_id?: string;
  sort_order: number;
  data_scope: string;
  menu_ids?: string[];
  status?: string;
  remark?: string;
}

export interface RoleUpdateRequest {
  role_name?: string;
  role_key?: string;
  sort_order?: number;
  data_scope?: string;
  menu_ids?: string[];
  status?: string;
  remark?: string;
}

export interface BatchUpdateRoleStatusRequest {
  role_ids: string[];
  status: string;
}

export interface RoleStatsResponse {
  total_roles: number;
  enabled_roles: number;
  disabled_roles: number;
  platform_roles: number;
  company_roles: number;
}

export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data?: T;
}

// 获取角色列表
export async function getRoleList(params: RoleListParams): Promise<ApiResponse<RoleListResponse>> {
  return request('/api/v1/roles', {
    method: 'GET',
    params,
  });
}

// 根据ID获取角色详情
export async function getRoleById(roleId: string): Promise<ApiResponse<RoleInfo>> {
  return request(`/api/v1/roles/${roleId}`, {
    method: 'GET',
  });
}

// 创建角色
export async function createRole(data: RoleCreateRequest): Promise<ApiResponse<RoleInfo>> {
  return request('/api/v1/roles', {
    method: 'POST',
    data,
  });
}

// 更新角色
export async function updateRole(roleId: string, data: RoleUpdateRequest): Promise<ApiResponse<RoleInfo>> {
  return request(`/api/v1/roles/${roleId}`, {
    method: 'PUT',
    data,
  });
}

// 删除角色
export async function deleteRole(roleId: string): Promise<ApiResponse> {
  return request(`/api/v1/roles/${roleId}`, {
    method: 'DELETE',
  });
}

// 批量更新角色状态
export async function batchUpdateRoleStatus(data: BatchUpdateRoleStatusRequest): Promise<ApiResponse> {
  return request('/api/v1/roles/batch-status', {
    method: 'PUT',
    data,
  });
}

// 获取角色统计信息
export async function getRoleStats(params?: { company_id?: string }): Promise<ApiResponse<RoleStatsResponse>> {
  return request('/api/v1/roles/stats', {
    method: 'GET',
    params,
  });
}

// 根据公司ID获取角色列表
export async function getRolesByCompanyId(companyId: string): Promise<ApiResponse<RoleInfo[]>> {
  return request(`/api/v1/roles/company/${companyId}`, {
    method: 'GET',
  });
} 