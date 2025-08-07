# 公司管理导入导出功能

## 功能概述

本系统为公司管理模块提供了完整的数据导入导出功能，支持批量操作，提高数据管理效率。

## 主要功能

### 1. 数据导入功能

- **文件格式支持**: Excel (.xlsx, .xls) 和 CSV (.csv) 格式
- **文件大小限制**: 最大 10MB
- **预览功能**: 导入前可预览数据，检查格式和内容
- **错误检查**: 自动验证数据格式，显示详细错误信息
- **更新模式**: 支持更新已存在的公司数据

#### 导入流程
1. 下载数据模板（Excel 或 CSV 格式）
2. 按照模板格式填写公司数据
3. 上传填写好的文件
4. 预览数据并检查错误
5. 确认导入，查看导入结果

#### 支持的字段
- 公司名称 (必填)
- 公司代码
- 负责人信息（中文、英文）
- 联系人
- 联系方式（电话、手机、邮箱）
- 地址信息（中文、英文地址）
- 业务信息（Broker Code、链接）
- 系统字段（有效期、用户配额、状态、备注）

### 2. 数据导出功能

- **导出范围**: 
  - 全部数据：导出系统中所有公司数据
  - 选中数据：导出当前选中的公司数据
  - 筛选结果：导出当前筛选条件下的数据
  
- **导出格式**: 
  - Excel 格式 (.xlsx) - 推荐格式，支持丰富的数据格式和样式
  - CSV 格式 (.csv) - 通用格式，可用于各种数据处理工具
  
- **导出选项**:
  - 数据文件：导出包含实际数据的文件
  - 模板文件：导出空白模板，用于批量导入时的格式参考

## 技术实现

### 前端组件

1. **ImportModal 组件**
   - 位置: `src/pages/company-list/components/ImportModal.tsx`
   - 功能: 处理文件上传、数据预览、导入确认等流程
   - 特性: 多步骤向导、实时错误检查、进度显示

2. **ExportModal 组件**
   - 位置: `src/pages/company-list/components/ExportModal.tsx`
   - 功能: 配置导出选项、格式选择、范围设定
   - 特性: 灵活的导出配置、实时预览选择

### API 接口

1. **导入相关接口**
   ```typescript
   // 预览导入数据
   POST /api/company/import/preview
   
   // 导入公司数据
   POST /api/company/import
   ```

2. **导出相关接口**
   ```typescript
   // 导出公司数据
   POST /api/company/export
   
   // 下载模板
   GET /api/company/template
   ```

### 类型定义

```typescript
// 导入请求类型
type CompanyImportRequest = {
  file: File;
  skip_header?: boolean;
  update_existing?: boolean;
};

// 导入响应类型
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

// 导出请求类型
type CompanyExportRequest = {
  status?: string;
  keyword?: string;
  ids?: string[];
  export_type?: 'all' | 'selected' | 'filtered';
  format?: 'xlsx' | 'csv';
  template?: boolean;
};
```

## 国际化支持

系统支持中文和英文界面，所有提示信息和界面文本都已本地化。

- 中文翻译: `src/locales/zh-CN/pages.ts`
- 英文翻译: `src/locales/en-US/pages.ts`

## 使用说明

### 导入数据
1. 在公司列表页面点击"数据管理"下拉菜单
2. 选择"导入数据"
3. 按照向导提示完成导入操作

### 导出数据
1. 在公司列表页面点击"数据管理"下拉菜单
2. 选择"导出数据"
3. 配置导出选项后点击"开始导出"

## 注意事项

1. **文件格式**: 确保上传的文件格式正确（Excel 或 CSV）
2. **数据验证**: 导入前会进行数据格式验证，请根据错误提示修正数据
3. **权限控制**: 导入导出功能需要相应的权限
4. **数据备份**: 建议在大批量导入前先备份现有数据
5. **网络稳定**: 大文件导入导出时请确保网络连接稳定

## 开发扩展

如需添加新的导入导出字段或修改验证规则：

1. 更新类型定义 (`src/services/ant-design-pro/typings.d.ts`)
2. 修改组件逻辑 (ImportModal.tsx, ExportModal.tsx)
3. 更新国际化文件
4. 调整后端API接口（如需要）

## 版本历史

- v1.0.0 - 初始版本，支持基本的导入导出功能
  - 支持 Excel 和 CSV 格式
  - 数据预览和错误检查
  - 多种导出选项
  - 完整的国际化支持 