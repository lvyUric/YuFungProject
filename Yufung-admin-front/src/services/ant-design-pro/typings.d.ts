// @ts-ignore
/* eslint-disable */

declare namespace API {
  type CurrentUser = {
    id?: string;
    user_id?: string;
    username?: string;
    display_name?: string;
    company_id?: string;
    role_ids?: string[];
    status?: string;
    email?: string;
    phone?: string;
    last_login?: string;
    avatar?: string;
    signature?: string;
    title?: string;
    group?: string;
    tags?: { key?: string; label?: string }[];
    notifyCount?: number;
    unreadCount?: number;
    country?: string;
    access?: string;
    geographic?: {
      province?: { label?: string; key?: string };
      city?: { label?: string; key?: string };
    };
    address?: string;
  };

  type LoginResult = {
    code?: number;
    message?: string;
    data?: {
      token?: string;
      refresh_token?: string;
      expires_at?: string;
      user?: CurrentUser;
    };
    // 保留原有字段以兼容其他用法
    success?: boolean;
    token?: string;
    refresh_token?: string;
    expires_at?: string;
    user?: CurrentUser;
    status?: string;
    type?: string;
    currentAuthority?: string;
  };

  type LoginParams = {
    username?: string;
    password?: string;
    autoLogin?: boolean;
    type?: string;
  };

  type RegisterParams = {
    username: string;
    display_name: string;
    password: string;
    email?: string;
    phone?: string;
  };

  type ChangePasswordParams = {
    old_password: string;
    new_password: string;
  };

  type PageParams = {
    current?: number;
    pageSize?: number;
  };

  type RuleListItem = {
    key?: number;
    disabled?: boolean;
    href?: string;
    avatar?: string;
    name?: string;
    owner?: string;
    desc?: string;
    callNo?: number;
    status?: number;
    updatedAt?: string;
    createdAt?: string;
    progress?: number;
  };

  type RuleList = {
    data?: RuleListItem[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };

  type FakeCaptcha = {
    code?: number;
    status?: string;
  };

  type ErrorResponse = {
    /** 业务约定的错误码 */
    errorCode: string;
    /** 业务上的错误信息 */
    errorMessage?: string;
    /** 业务上的请求是否成功 */
    success?: boolean;
  };

  type NoticeIconList = {
    data?: NoticeIconItem[];
    /** 列表的内容总数 */
    total?: number;
    success?: boolean;
  };

  type NoticeIconItemType = 'notification' | 'message' | 'event';

  type NoticeIconItem = {
    id?: string;
    extra?: string;
    key?: string;
    read?: boolean;
    avatar?: string;
    title?: string;
    status?: string;
    datetime?: string;
    description?: string;
    type?: NoticeIconItemType;
  };

  // 通用响应类型
  type Response<T = any> = {
    success: boolean;
    message: string;
    code: number;
    data?: T;
  };

  // 公司管理相关类型
  type CompanyInfo = {
    id?: string;
    company_id?: string;
    company_name?: string;
    company_code?: string;
    
    // 负责人信息
    contact_person?: string;
    
    // 联系方式
    tel_no?: string;
    mobile?: string;
    contact_phone?: string;
    email?: string;
    
    // 中文地址信息 - 补全所有地址字段
    address_cn_province?: string;
    address_cn_city?: string;
    address_cn_district?: string;
    address_cn_detail?: string;
    
    // 英文地址信息 - 补全所有地址字段
    address_en_province?: string;
    address_en_city?: string;
    address_en_district?: string;
    address_en_detail?: string;
    
    address?: string; // 原有地址字段（保留兼容）
    
    // 业务信息
    broker_code?: string;
    link?: string;
    
    // 登录信息
    username?: string;
    
    // 系统字段
    valid_start_date?: string;
    valid_end_date?: string;
    user_quota?: number;
    current_user_count?: number;
    status?: string;
    status_text?: string;
    remark?: string; // 备注信息
    submitted_by?: string;
    created_at?: string;
    updated_at?: string;
  };

  type CreateCompanyRequest = {
    company_name: string;
    company_code?: string;
    
    // 负责人信息
    contact_person?: string;
    
    // 联系方式
    tel_no?: string;
    mobile?: string;
    contact_phone?: string;
    email: string;
    
    // 中文地址信息 - 补全所有地址字段
    address_cn_province?: string;
    address_cn_city?: string;
    address_cn_district?: string;
    address_cn_detail?: string;
    
    // 英文地址信息 - 补全所有地址字段
    address_en_province?: string;
    address_en_city?: string;
    address_en_district?: string;
    address_en_detail?: string;
    
    address?: string; // 原有地址字段（保留兼容）
    
    // 业务信息
    broker_code?: string;
    link?: string;
    
    // 登录信息
    username?: string;
    password?: string;
    
    // 系统字段
    valid_start_date?: string;
    valid_end_date?: string;
    user_quota?: number;
    remark?: string;
  };

  type UpdateCompanyRequest = {
    company_name?: string;
    company_code?: string;
    
    // 负责人信息
    contact_person?: string;
    
    // 联系方式
    tel_no?: string;
    mobile?: string;
    contact_phone?: string;
    email?: string;
    
    // 中文地址信息 - 补全所有地址字段
    address_cn_province?: string;
    address_cn_city?: string;
    address_cn_district?: string;
    address_cn_detail?: string;
    
    // 英文地址信息 - 补全所有地址字段
    address_en_province?: string;
    address_en_city?: string;
    address_en_district?: string;
    address_en_detail?: string;
    
    address?: string; // 原有地址字段（保留兼容）
    
    // 业务信息
    broker_code?: string;
    link?: string;
    
    // 登录信息
    username?: string;
    password?: string;
    
    // 系统字段
    valid_start_date?: string;
    valid_end_date?: string;
    user_quota?: number;
    status?: string;
    remark?: string;
  };

  type CompanyQueryParams = {
    page?: number;
    page_size?: number;
    status?: string;
    keyword?: string;
  };

  type CompanyListResponse = {
    companies?: CompanyInfo[];
    total?: number;
    page?: number;
    page_size?: number;
    total_pages?: number;
  };

  type CompanyStatsResponse = {
    total_companies?: number;
    active_companies?: number;
    expired_companies?: number;
    total_users?: number;
  };

  // 导入导出相关类型
  type CompanyImportRequest = {
    file: File;
    skip_header?: boolean;
    update_existing?: boolean;
  };

  type CompanyImportResponse = {
    success_count?: number;
    error_count?: number;
    total_count?: number;
    errors?: Array<{
      row: number;
      errors: string[];
      data?: any;
    }>;
    preview?: CompanyInfo[];
  };

  type CompanyExportRequest = {
    status?: string;
    keyword?: string;
    ids?: string[];
    export_type?: 'all' | 'selected' | 'filtered';
    format?: 'xlsx' | 'csv';
    template?: boolean;
  };

  type CompanyExportResponse = {
    file_url?: string;
    file_name?: string;
    download_token?: string;
  };
}
