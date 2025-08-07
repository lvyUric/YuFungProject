import { request } from '@umijs/max';

export interface MenuInfo {
  id: string;
  menu_id: string;
  parent_id: string;
  menu_name: string;
  menu_type: string;
  route_path: string;
  component: string;
  permission_code: string;
  icon: string;
  sort_order: number;
  visible: boolean;
  status: string;
  created_at: string;
  updated_at: string;
  children: MenuInfo[];
}

export interface MenuQueryParams {
  menu_name?: string;
  menu_type?: string;
  status?: string;
  visible?: boolean;
  permission_code?: string;
}

export interface MenuListResponse {
  menus: MenuInfo[];
  total: number;
}

export interface MenuCreateRequest {
  parent_id?: string;
  menu_name: string;
  menu_type: string;
  route_path?: string;
  component?: string;
  permission_code?: string;
  icon?: string;
  sort_order: number;
  visible: boolean;
  status?: string;
}

export interface MenuUpdateRequest {
  parent_id?: string;
  menu_name?: string;
  menu_type?: string;
  route_path?: string;
  component?: string;
  permission_code?: string;
  icon?: string;
  sort_order?: number;
  visible?: boolean;
  status?: string;
}

export interface BatchUpdateMenuStatusRequest {
  menu_ids: string[];
  status: string;
}

export interface MenuStatsResponse {
  total_menus: number;
  enabled_menus: number;
  disabled_menus: number;
  directory_menus: number;
  page_menus: number;
  button_menus: number;
}

export interface UserMenuResponse {
  menu_id: string;
  menu_name: string;
  route_path: string;
  component: string;
  icon: string;
  sort_order: number;
  children: UserMenuResponse[];
}

export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data?: T;
}

// 获取菜单列表
export async function getMenuList(params: MenuQueryParams): Promise<ApiResponse<MenuListResponse>> {
  return request('/api/v1/menus', {
    method: 'GET',
    params,
  });
}

// 获取菜单树
export const getMenuTree = async (params: MenuQueryParams): Promise<ApiResponse<MenuInfo[]>> => {
  return request('/api/v1/menu/tree', {
    method: 'GET',
    params,
  });
};

// 根据ID获取菜单详情
export async function getMenuById(menuId: string): Promise<ApiResponse<MenuInfo>> {
  return request(`/api/v1/menus/${menuId}`, {
    method: 'GET',
  });
}

// 创建菜单
export async function createMenu(data: MenuCreateRequest): Promise<ApiResponse<MenuInfo>> {
  return request('/api/v1/menus', {
    method: 'POST',
    data,
  });
}

// 更新菜单
export async function updateMenu(menuId: string, data: MenuUpdateRequest): Promise<ApiResponse<MenuInfo>> {
  return request(`/api/v1/menus/${menuId}`, {
    method: 'PUT',
    data,
  });
}

// 删除菜单
export async function deleteMenu(menuId: string): Promise<ApiResponse> {
  return request(`/api/v1/menus/${menuId}`, {
    method: 'DELETE',
  });
}

// 批量更新菜单状态
export async function batchUpdateMenuStatus(data: BatchUpdateMenuStatusRequest): Promise<ApiResponse> {
  return request('/api/v1/menus/batch-status', {
    method: 'PUT',
    data,
  });
}

// 获取菜单统计信息
export async function getMenuStats(): Promise<ApiResponse<MenuStatsResponse>> {
  return request('/api/v1/menus/stats', {
    method: 'GET',
  });
}

// 获取用户菜单（前端菜单渲染）
export async function getUserMenus(menuIds?: string[]): Promise<ApiResponse<UserMenuResponse[]>> {
  return request('/api/v1/menus/user', {
    method: 'GET',
    params: {
      menu_ids: menuIds?.join(','),
    },
  });
} 