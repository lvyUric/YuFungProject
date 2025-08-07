// @ts-ignore
/* eslint-disable */
import { request } from '@umijs/max';

/** 获取公司列表 GET /api/company */
export async function getCompanyList(
  params: API.CompanyQueryParams,
  options?: { [key: string]: any },
) {
  return request<API.Response<API.CompanyListResponse>>('/api/company', {
    method: 'GET',
    params: {
      ...params,
    },
    ...(options || {}),
  });
}

/** 获取公司详情 GET /api/company/:id */
export async function getCompanyById(id: string, options?: { [key: string]: any }) {
  return request<API.Response<API.CompanyInfo>>(`/api/company/${id}`, {
    method: 'GET',
    ...(options || {}),
  });
}

/** 创建公司 POST /api/company */
export async function createCompany(body: API.CreateCompanyRequest, options?: { [key: string]: any }) {
  return request<API.Response<API.CompanyInfo>>('/api/company', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 更新公司 PUT /api/company/:id */
export async function updateCompany(
  id: string,
  body: API.UpdateCompanyRequest,
  options?: { [key: string]: any },
) {
  return request<API.Response<API.CompanyInfo>>(`/api/company/${id}`, {
    method: 'PUT',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 删除公司 DELETE /api/company/:id */
export async function deleteCompany(id: string, options?: { [key: string]: any }) {
  return request<API.Response<any>>(`/api/company/${id}`, {
    method: 'DELETE',
    ...(options || {}),
  });
}

/** 获取公司统计 GET /api/company/stats */
export async function getCompanyStats(options?: { [key: string]: any }) {
  return request<API.Response<API.CompanyStatsResponse>>('/api/company/stats', {
    method: 'GET',
    ...(options || {}),
  });
}

/** 导出公司数据 POST /api/company/export */
export async function exportCompany(body: API.CompanyExportRequest, options?: { [key: string]: any }) {
  return request<API.Response<API.CompanyExportResponse>>('/api/company/export', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    data: body,
    ...(options || {}),
  });
}

/** 下载导出模板 GET /api/company/template */
export async function downloadCompanyTemplate(format: 'xlsx' | 'csv' = 'xlsx', options?: { [key: string]: any }) {
  return request<Blob>('/api/company/template', {
    method: 'GET',
    params: { format },
    responseType: 'blob',
    ...(options || {}),
  });
}

/** 导入公司数据 POST /api/company/import */
export async function importCompany(formData: FormData, options?: { [key: string]: any }) {
  return request<API.Response<API.CompanyImportResponse>>('/api/company/import', {
    method: 'POST',
    data: formData,
    requestType: 'form',
    ...(options || {}),
  });
}

/** 预览导入数据 POST /api/company/import/preview */
export async function previewCompanyImport(formData: FormData, options?: { [key: string]: any }) {
  return request<API.Response<API.CompanyImportResponse>>('/api/company/import/preview', {
    method: 'POST',
    data: formData,
    requestType: 'form',
    ...(options || {}),
  });
} 