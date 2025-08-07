import { request } from '@umijs/max';

// 保单数据类型定义
export interface PolicyInfo {
  id: string;
  policy_id: string;
  serial_number: number;
  account_number: string;
  customer_number: string;
  customer_name_cn: string;
  customer_name_en: string;
  proposal_number: string;
  policy_currency: 'USD' | 'HKD' | 'CNY';
  partner: string;
  referral_code: string;
  hk_manager: string;
  referral_pm: string;
  referral_branch: string;
  referral_sub_branch: string;
  referral_date?: string;
  is_surrendered: boolean;
  payment_date?: string;
  effective_date?: string;
  payment_method: '期缴' | '趸缴' | '预缴';
  payment_years: number;
  payment_periods: number;
  actual_premium: number;
  aum: number;
  past_cooling_period: boolean;
  is_paid_commission: boolean;
  is_employee: boolean;
  referral_rate: number;
  exchange_rate: number;
  expected_fee: number;
  payment_pay_date?: string;
  insurance_company: string;
  product_name: string;
  product_type: string;
  remark: string;
  company_id: string;
  created_by: string;
  updated_by: string;
  created_at: string;
  updated_at: string;
}

// 查询参数
export interface PolicyListParams {
  page?: number;
  page_size?: number;
  account_number?: string;
  customer_number?: string;
  customer_name_cn?: string;
  customer_name_en?: string;
  proposal_number?: string;
  policy_currency?: 'USD' | 'HKD' | 'CNY';
  partner?: string;
  referral_code?: string;
  hk_manager?: string;
  referral_pm?: string;
  referral_branch?: string;
  referral_sub_branch?: string;
  payment_method?: '期缴' | '趸缴' | '预缴';
  insurance_company?: string;
  product_name?: string;
  product_type?: string;
  is_surrendered?: boolean;
  past_cooling_period?: boolean;
  is_paid_commission?: boolean;
  is_employee?: boolean;
  referral_date_start?: string;
  referral_date_end?: string;
  payment_date_start?: string;
  payment_date_end?: string;
  effective_date_start?: string;
  effective_date_end?: string;
  sort_by?: string;
  sort_order?: 'asc' | 'desc';
}

// 列表响应
export interface PolicyListResponse {
  list: PolicyInfo[];
  total: number;
  page: number;
  page_size: number;
  total_pages: number;
}

// 创建保单请求
export interface PolicyCreateRequest {
  account_number?: string; // 改为非必填
  customer_number: string;
  customer_name_cn: string;
  customer_name_en?: string;
  proposal_number: string;
  policy_currency: 'USD' | 'HKD' | 'CNY';
  partner?: string;
  referral_code?: string;
  hk_manager?: string;
  referral_pm?: string;
  referral_branch?: string;
  referral_sub_branch?: string;
  referral_date?: string;
  is_surrendered?: boolean;
  payment_date?: string;
  effective_date?: string;
  payment_method?: '期缴' | '趸缴' | '预缴';
  payment_years?: number;
  payment_periods?: number;
  actual_premium?: number;
  aum?: number;
  past_cooling_period?: boolean;
  is_paid_commission?: boolean;
  is_employee?: boolean;
  referral_rate?: number;
  exchange_rate?: number; // 汇率字段，保留4位小数
  expected_fee?: number;
  payment_pay_date?: string;
  insurance_company: string;
  product_name: string;
  product_type: string;
  remark?: string;
}

// 更新保单请求
export interface PolicyUpdateRequest {
  customer_name_cn?: string;
  customer_name_en?: string;
  policy_currency?: 'USD' | 'HKD' | 'CNY';
  partner?: string;
  referral_code?: string;
  hk_manager?: string;
  referral_pm?: string;
  referral_branch?: string;
  referral_sub_branch?: string;
  referral_date?: string;
  is_surrendered?: boolean;
  payment_date?: string;
  effective_date?: string;
  payment_method?: '期缴' | '趸缴' | '预缴';
  payment_years?: number;
  payment_periods?: number;
  actual_premium?: number;
  aum?: number;
  past_cooling_period?: boolean;
  is_paid_commission?: boolean;
  is_employee?: boolean;
  referral_rate?: number;
  exchange_rate?: number;
  expected_fee?: number;
  payment_pay_date?: string;
  insurance_company?: string;
  product_name?: string;
  product_type?: string;
  remark?: string;
}

// 保单统计
export interface PolicyStatistics {
  total_policies: number;
  total_premium: number;
  total_aum: number;
  total_expected_fee: number;
  surrendered_count: number;
  employee_count: number;
  cooling_period_count: number;
  paid_commission_count: number;
}

// 批量更新状态请求
export interface BatchUpdatePolicyStatusRequest {
  policy_ids: string[];
  is_surrendered?: boolean;
  past_cooling_period?: boolean;
  is_paid_commission?: boolean;
}

// API调用方法

/** 获取保单列表 */
export async function getPolicyList(params: PolicyListParams) {
  return request<{
    code: number;
    data: PolicyListResponse;
    message: string;
  }>('/api/policies', {
    method: 'GET',
    params,
  });
}

/** 获取保单详情 */
export async function getPolicyDetail(policyId: string) {
  return request<{
    code: number;
    data: PolicyInfo;
    message: string;
  }>(`/api/policies/${policyId}`, {
    method: 'GET',
  });
}

/** 创建保单 */
export async function createPolicy(data: PolicyCreateRequest) {
  return request<{
    code: number;
    data: PolicyInfo;
    message: string;
  }>('/api/policies', {
    method: 'POST',
    data,
  });
}

/** 更新保单 */
export async function updatePolicy(policyId: string, data: PolicyUpdateRequest) {
  return request<{
    code: number;
    data: PolicyInfo;
    message: string;
  }>(`/api/policies/${policyId}`, {
    method: 'PUT',
    data,
  });
}

/** 删除保单 */
export async function deletePolicy(policyId: string) {
  return request<{
    code: number;
    message: string;
  }>(`/api/policies/${policyId}`, {
    method: 'DELETE',
  });
}

/** 获取保单统计 */
export async function getPolicyStatistics() {
  return request<{
    code: number;
    data: PolicyStatistics;
    message: string;
  }>('/api/policies/statistics', {
    method: 'GET',
  });
}

/** 批量更新保单状态 */
export async function batchUpdatePolicyStatus(data: BatchUpdatePolicyStatusRequest) {
  return request<{
    code: number;
    message: string;
  }>('/api/policies/batch-update', {
    method: 'POST',
    data,
  });
}

/** 导入保单 */
export async function importPolicies(data: { data: PolicyCreateRequest[] }) {
  return request<{
    code: number;
    data: {
      success_count: number;
      error_count: number;
      success_ids: string[];
      errors: string[];
    };
    message: string;
  }>('/api/policies/import', {
    method: 'POST',
    data,
  });
}

/** 导出保单 */
export async function exportPolicies(data: { policy_ids?: string[]; export_type: 'xlsx' | 'csv' }) {
  return request<{
    code: number;
    data: PolicyInfo[];
    message: string;
  }>('/api/policies/export', {
    method: 'POST',
    data,
  });
}

/** 下载保单模板 */
export async function downloadPolicyTemplate(format: 'xlsx' | 'csv' = 'xlsx') {
  return request<Blob>('/api/policies/template', {
    method: 'GET',
    params: { type: format },
    responseType: 'blob',
  });
}

/** 获取字段验证规则 */
export async function getPolicyValidationRules() {
  return request<{
    code: number;
    data: Record<string, any>;
    message: string;
  }>('/api/policies/validation-rules', {
    method: 'GET',
  });
}

// 系统配置相关接口
export interface SystemConfigOption {
  config_id: string;
  config_key: string;
  config_value: string;
  display_name: string;
}

// 获取系统配置选项
export async function getSystemConfigOptions(configType: string): Promise<API.Response<SystemConfigOption[]>> {
  return request(`/api/system-configs/options/${configType}`, {
    method: 'GET',
  });
}

// 获取公司列表（用于承保公司下拉框）
export async function getCompanyOptions(): Promise<API.Response<any[]>> {
  return request('/api/company', {
    method: 'GET',
    params: {
      page: 1,
      page_size: 1000, // 获取所有公司
      status: 'active', // 只获取激活的公司
    },
  });
}

// 导入导出相关类型定义
export interface PolicyImportError {
  row: number;
  errors: string[];
  data: any;
}

export interface PolicyImportResponse {
  success_count: number;
  error_count: number;
  total_count: number;
  errors: PolicyImportError[];
  preview?: PolicyCreateRequest[];
}

export interface PolicyExportRequest {
  policy_ids?: string[];
  export_type: 'xlsx' | 'csv';
}

// 保单导入导出API
export async function previewPolicyImport(formData: FormData) {
  return request<API.Response<PolicyImportResponse>>('/api/policies/import/preview', {
    method: 'POST',
    data: formData,
  });
}

export async function importPoliciesFromFile(formData: FormData) {
  return request<API.Response<PolicyImportResponse>>('/api/policies/import', {
    method: 'POST',
    data: formData,
  });
}

export async function exportPoliciesToFile(
  params: PolicyExportRequest & { format?: 'xlsx' | 'csv' }
) {
  return request(`/api/policies/export?format=${params.format || 'xlsx'}`, {
    method: 'POST',
    data: {
      policy_ids: params.policy_ids,
      export_type: params.export_type,
    },
    responseType: 'blob',
  });
} 