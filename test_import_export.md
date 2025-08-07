# 公司管理导入导出功能测试

## 功能已实现说明

✅ **后端API接口已完成**
- 导出公司数据：`POST /api/company/export`
- 下载导入模板：`GET /api/company/template?format=xlsx|csv`
- 预览导入数据：`POST /api/company/import/preview`
- 导入公司数据：`POST /api/company/import`

✅ **前端组件已完成**
- ImportModal：导入功能模态框
- ExportModal：导出功能模态框
- 集成到公司列表页面的"数据管理"按钮

✅ **功能特性**
- 支持 Excel(.xlsx, .xls) 和 CSV(.csv) 格式
- 数据预览和错误检查
- 批量导入和导出
- 多种导出选项（全部、选中、筛选）
- 完整的国际化支持（中英文）

## 测试步骤

### 1. 启动服务

#### 启动后端服务
```bash
cd D:\YuFungProject
go run ./cmd/main.go
```

#### 启动前端服务
在PowerShell中使用：
```powershell
cd D:\YuFungProject\Yufung-admin-front
npm run dev
```

或在CMD中使用：
```cmd
cd D:\YuFungProject\Yufung-admin-front && npm run dev
```

### 2. 测试导出功能

1. 登录系统，进入公司列表页面
2. 点击页面顶部的"数据管理"下拉按钮
3. 选择"导出数据"
4. 在弹出的对话框中：
   - 选择导出范围（全部数据/选中数据/筛选结果）
   - 选择导出格式（Excel/CSV）
   - 选择导出选项（数据文件/模板文件）
5. 点击"开始导出"

### 3. 测试导入功能

1. 在公司列表页面点击"数据管理" → "导入数据"
2. 在导入模态框中：
   - 下载Excel或CSV模板
   - 按模板格式填写公司数据
   - 上传填写好的文件
   - 配置导入选项（跳过表头、更新已存在公司）
3. 点击"预览数据"查看解析结果
4. 确认无误后点击"确认导入"

### 4. API测试（可选）

#### 下载模板
```bash
curl -H "Authorization: Bearer YOUR_TOKEN" \
     -o template.xlsx \
     "http://localhost:8080/api/company/template?format=xlsx"
```

#### 导出数据
```bash
curl -X POST \
     -H "Authorization: Bearer YOUR_TOKEN" \
     -H "Content-Type: application/json" \
     -d '{"export_type":"all","format":"xlsx","template":false}' \
     "http://localhost:8080/api/company/export"
```

## 已知问题和解决方案

### 问题1：PowerShell中 `&&` 语法错误
**解决方案：** 在PowerShell中使用分号`;`代替`&&`
```powershell
cd Yufung-admin-front; npm run dev
```

### 问题2：导出返回404错误
**原因：** 后端API接口未实现
**解决方案：** 已完成后端实现，确保后端服务正在运行

### 问题3：CORS跨域问题
**解决方案：** 确保后端已配置CORS中间件，允许前端域名访问

## 数据模板格式

导入Excel/CSV文件应包含以下列：

| 列名 | 必填 | 说明 |
|------|------|------|
| 公司名称 | ✓ | 公司的显示名称 |
| 公司代码 |  | 内部公司代码 |
| 负责人中文名 |  | 负责人中文姓名 |
| 负责人英文名 |  | 负责人英文姓名 |
| 联络人 |  | 主要联系人 |
| 固定电话 |  | 公司固定电话 |
| 移动电话 |  | 联系手机号 |
| 邮箱地址 | ✓ | 公司联系邮箱 |
| 中文地址省份 |  | 省/自治区/直辖市 |
| 中文地址城市 |  | 城市 |
| 中文地址区县 |  | 区/县 |
| 中文地址详细 |  | 详细地址 |
| 英文地址省份 |  | 英文省份 |
| 英文地址城市 |  | 英文城市 |
| 英文地址区县 |  | 英文区县 |
| 英文地址详细 |  | 英文详细地址 |
| 经纪人代码 |  | Broker Code |
| 相关链接 |  | 公司相关链接 |
| 用户名 |  | 登录用户名 |
| 备注信息 |  | 备注说明 |
| 有效期开始 |  | 格式：YYYY-MM-DD |
| 有效期结束 |  | 格式：YYYY-MM-DD |
| 用户配额 |  | 数字，默认为1 |

## 开发完成状态

### ✅ 已完成
- [x] 后端API接口实现
- [x] 前端导入导出组件
- [x] 路由配置
- [x] 类型定义
- [x] 国际化支持
- [x] 错误处理
- [x] 数据验证

### 🔧 可优化项
- [ ] 文件下载的临时URL处理
- [ ] 大文件导入的分片处理
- [ ] 导入进度条显示
- [ ] 文件存储到对象存储服务
- [ ] 导入历史记录
- [ ] 数据导入回滚功能

## 技术栈

- **后端**：Go + Gin + MongoDB + excelize
- **前端**：React + Ant Design Pro + TypeScript
- **文件处理**：支持Excel和CSV格式
- **国际化**：中英文双语支持 