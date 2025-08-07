import { request } from '@umijs/max';

// 用户相关类型定义
export interface UserItem {
  id: string;
  user_id: string;
  username: string;
  display_name: string;
  company_id: string;
  role_ids: string[];
  status: 'active' | 'inactive' | 'locked';
  email: string;
  phone: string;
  last_login?: string;
  created_at: string;
  updated_at: string;
  remark?: string;
}

export interface UserCreateRequest {
  username: string;
  display_name: string;
  password: string;
  company_id: string;
  role_ids: string[];
  email?: string;
  phone?: string;
  remark?: string;
}

export interface UserUpdateRequest {
  display_name?: string;
  role_ids?: string[];
  email?: string;
  phone?: string;
  remark?: string;
  status?: 'active' | 'inactive' | 'locked';
}

export interface UserListParams {
  page?: number;
  page_size?: number;
  username?: string;
  display_name?: string;
  company_id?: string;
  status?: string;
}

export interface UserListResponse {
  users: UserItem[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

export interface ResetPasswordRequest {
  new_password: string;
}

export interface ApiResponse<T = any> {
  code: number;
  message: string;
  data?: T;
}

// 用户API接口
export const userService = {
  // 获取用户列表
  async getUserList(params: UserListParams): Promise<ApiResponse<UserListResponse>> {
    return request('/api/v1/users', {
      method: 'GET',
      params,
    });
  },

  // 获取用户详情
  async getUserById(userId: string): Promise<ApiResponse<UserItem>> {
    return request(`/api/v1/users/${userId}`, {
      method: 'GET',
    });
  },

  // 创建用户
  async createUser(userData: UserCreateRequest): Promise<ApiResponse<UserItem>> {
    return request('/api/v1/users', {
      method: 'POST',
      data: userData,
    });
  },

  // 更新用户
  async updateUser(userId: string, userData: UserUpdateRequest): Promise<ApiResponse> {
    return request(`/api/v1/users/${userId}`, {
      method: 'PUT',
      data: userData,
    });
  },

  // 删除用户
  async deleteUser(userId: string): Promise<ApiResponse> {
    return request(`/api/v1/users/${userId}`, {
      method: 'DELETE',
    });
  },

  // 重置密码
  async resetPassword(userId: string, passwordData: ResetPasswordRequest): Promise<ApiResponse> {
    return request(`/api/v1/users/${userId}/reset-password`, {
      method: 'PUT',
      data: passwordData,
    });
  },

  // 批量更新用户状态
  async batchUpdateStatus(userIds: string[], status: string): Promise<ApiResponse> {
    return request('/api/v1/users/batch-status', {
      method: 'PUT',
      data: {
        user_ids: userIds,
        status,
      },
    });
  },

  // 快捷停用用户
  async quickDisableUser(userId: string): Promise<ApiResponse> {
    return request(`/api/v1/users/${userId}/quick-disable`, {
      method: 'PUT',
    });
  },

  // 导出用户数据
  async exportUsers(params: any): Promise<Blob> {
    return request('/api/v1/users/export', {
      method: 'GET',
      params,
      responseType: 'blob',
    }) as Promise<Blob>;
  },

  // 高级导出用户数据
  async exportUsersAdvanced(data: any): Promise<ApiResponse> {
    return request('/api/v1/users/export-advanced', {
      method: 'POST',
      data,
    });
  },

  // 下载用户导入模板
  async downloadUserTemplate(format: string = 'xlsx'): Promise<Blob> {
    return request('/api/v1/users/template', {
      method: 'GET',
      params: { format },
      responseType: 'blob',
    }) as Promise<Blob>;
  },

  // 预览用户导入
  async previewUserImport(formData: FormData): Promise<ApiResponse> {
    return request('/api/v1/users/import/preview', {
      method: 'POST',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },

  // 导入用户数据
  async importUsers(formData: FormData): Promise<ApiResponse> {
    return request('/api/v1/users/import', {
      method: 'POST',
      data: formData,
      headers: {
        'Content-Type': 'multipart/form-data',
      },
    });
  },
};

// 导出单独的函数（与组件中的导入保持一致）
export const getUserList = userService.getUserList;
export const getUserById = userService.getUserById;
export const createUser = userService.createUser;
export const updateUser = userService.updateUser;
export const deleteUser = userService.deleteUser;
export const resetPassword = userService.resetPassword;
export const batchUpdateUserStatus = userService.batchUpdateStatus;
export const quickDisableUser = userService.quickDisableUser;
export const exportUsers = userService.exportUsers;
export const exportUsersAdvanced = userService.exportUsersAdvanced;
export const downloadUserTemplate = userService.downloadUserTemplate;
export const previewUserImport = userService.previewUserImport;
export const importUsers = userService.importUsers; 