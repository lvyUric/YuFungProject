// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

const API_BASE_URL = 'http://106.52.172.124:8088';

/** 获取当前的用户 GET /api/auth/user-info */
export async function currentUser(options?: { [key: string]: any }) {
  return request<API.Response<API.CurrentUser>>(`${API_BASE_URL}/api/auth/user-info`, {
    method: 'GET',
    ...(options || {}),
  });
}

/** 退出登录接口 POST /api/auth/logout */
export async function outLogin(options?: { [key: string]: any }) {
  return request<API.Response<any>>(`${API_BASE_URL}/api/auth/logout`, {
    method: 'POST',
    ...(options || {}),
  });
}

/** 登录接口 POST /api/auth/login */
export async function login(body: API.LoginParams, options?: { [key: string]: any }) {
  return request<API.LoginResult>(`${API_BASE_URL}/api/auth/login`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 用户注册接口 POST /api/auth/register */
export async function register(body: API.RegisterParams, options?: { [key: string]: any }) {
  return request<API.Response<API.CurrentUser>>(`${API_BASE_URL}/api/auth/register`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 修改密码接口 POST /api/auth/change-password */
export async function changePassword(body: API.ChangePasswordParams, options?: { [key: string]: any }) {
  return request<API.Response<any>>(`${API_BASE_URL}/api/auth/change-password`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 刷新Token接口 POST /api/auth/refresh */
export async function refreshToken(body: { refresh_token: string }, options?: { [key: string]: any }) {
  return request<API.LoginResult>(`${API_BASE_URL}/api/auth/refresh`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 此处后端没有提供注释 GET /api/notices */
export async function getNotices(options?: { [key: string]: any }) {
  return request<API.NoticeIconList>(`${API_BASE_URL}/api/notices`, {
    method: 'GET',
    ...(options || {}),
  });
}

/** 获取规则列表 GET /api/rule */
export async function rule(
  params: {
    // query
    /** 当前的页码 */
    current?: number;
    /** 页面的容量 */
    pageSize?: number;
  },
  options?: { [key: string]: any },
) {
  return request<API.RuleList>(`${API_BASE_URL}/api/rule`, {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 更新规则 PUT /api/rule */
export async function updateRule(options?: { [key: string]: any }) {
  return request<API.RuleListItem>(`${API_BASE_URL}/api/rule`, {
    method: 'POST',
    data: {
      method: 'update',
      ...(options || {}),
    },
  });
}

/** 新建规则 POST /api/rule */
export async function addRule(options?: { [key: string]: any }) {
  return request<API.RuleListItem>(`${API_BASE_URL}/api/rule`, {
    method: 'POST',
    data: {
      method: 'post',
      ...(options || {}),
    },
  });
}

/** 删除规则 DELETE /api/rule */
export async function removeRule(options?: { [key: string]: any }) {
  return request<Record<string, any>>(`${API_BASE_URL}/api/rule`, {
    method: 'POST',
    data: {
      method: 'delete',
      ...(options || {}),
    },
  });
}
