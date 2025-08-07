# MongoDB 数据库升级指南

本文档说明如何执行公司集合的字段结构升级，以支持完整的保险公司管理功能。

## 📋 升级内容

### 新增字段列表

| 字段分类 | 字段名 | 类型 | 说明 |
|---------|-------|------|------|
| **基本信息** | `company_code` | String | 保险公司代码 |
| **负责人信息** | `principles_cn` | String | 负责人中文名 |
| | `principles_en` | String | 负责人英文名 |
| | `contact_person` | String | 联络人 |
| **联系方式** | `tel_no` | String | 固定电话 |
| | `mobile` | String | 移动电话 |
| **中文地址** | `address_cn_province` | String | 省/自治区/直辖市 |
| | `address_cn_city` | String | 市 |
| | `address_cn_district` | String | 县/区 |
| | `address_cn_detail` | String | 详细地址 |
| **英文地址** | `address_en_province` | String | Province/State |
| | `address_en_city` | String | City |
| | `address_en_district` | String | District |
| | `address_en_detail` | String | Detailed Address |
| **业务信息** | `broker_code` | String | 经纪人代码 |
| | `link` | String | 相关链接 |
| **登录信息** | `username` | String | 用户名 |
| | `password_hash` | String | 密码哈希值 |
| **扩展信息** | `remark1` | String | 备注1 |
| | `remark2` | String | 备注2 |
| | `email_template` | String | 邮箱模版 |
| | `payment_notification` | String | 缴费账号通知 |
| **系统字段** | `submitted_by` | String | 提交人 |

### 新建索引

- 单字段索引：公司代码、负责人姓名、联络人、联系方式、经纪人代码等
- 复合索引：地址组合、状态+有效期、公司名称+代码等

## 🚀 执行步骤

### 前置条件

1. **MongoDB 连接**: 确保可以连接到 MongoDB 数据库
2. **权限验证**: 确保有读写数据库的权限
3. **备份策略**: 建议在生产环境执行前先做完整备份

### 1. 修改数据库名称

在执行脚本前，请根据实际情况修改脚本中的数据库名称：

```javascript
// 将 'yufung_admin' 替换为你的实际数据库名
use('your_database_name');
```

### 2. 执行升级脚本

#### 方法一: 使用 MongoDB Shell

```bash
# 连接到 MongoDB
mongosh "mongodb://localhost:27017"

# 或连接到远程数据库
mongosh "mongodb://username:password@your-server:27017/your_database"

# 执行升级脚本
load('scripts/mongodb/upgrade_company_schema.js')
```

#### 方法二: 直接执行文件

```bash
# 本地执行
mongosh your_database_name --file scripts/mongodb/upgrade_company_schema.js

# 远程执行
mongosh "mongodb://username:password@your-server:27017/your_database" --file scripts/mongodb/upgrade_company_schema.js
```

### 3. 验证升级结果

升级完成后，脚本会自动验证：

- ✅ 字段是否正确添加
- ✅ 索引是否成功创建
- ✅ 数据迁移是否完成
- ✅ 备份是否成功创建

### 4. 测试应用程序

升级完成后，启动应用程序进行测试：

```bash
# 启动后端服务
cd your-backend-project
go run main.go

# 启动前端服务
cd your-frontend-project
npm start
```

访问公司管理页面，测试新字段的显示和编辑功能。

## 🔄 回滚操作

如果升级后发现问题需要回滚：

### 执行回滚脚本

```bash
# 方法一: MongoDB Shell 中执行
load('scripts/mongodb/rollback_company_schema.js')

# 方法二: 直接执行文件
mongosh your_database_name --file scripts/mongodb/rollback_company_schema.js
```

### 回滚后验证

回滚脚本会：

- ✅ 删除所有新增字段
- ✅ 删除新建索引
- ✅ 创建回滚前状态备份
- ✅ 验证回滚结果

## 📊 执行示例

### 成功升级的输出示例

```
=== 开始升级公司集合字段结构 ===
1. 创建集合备份...
✓ 备份完成，备份集合名: companies_backup
2. 为现有公司文档添加新字段...
✓ 字段添加完成
3. 执行数据迁移...
✓ 数据迁移完成
4. 创建新字段索引...
✓ 索引创建完成
5. 验证升级结果...
总公司数量: 5
✓ 所有新字段已成功添加
6. 当前索引列表:
- _id_: {"_id":1}
- company_code_1: {"company_code":1}
- principles_cn_1: {"principles_cn":1}
...
=== 公司集合字段升级完成 ===
```

## ⚠️ 注意事项

### 生产环境执行前

1. **完整备份**: 建议使用 `mongodump` 创建完整数据库备份
2. **测试环境验证**: 先在测试环境执行并验证
3. **停机时间**: 评估执行时间，可能需要安排维护窗口
4. **回滚方案**: 准备回滚计划

### 数据迁移说明

- 原有 `address` 字段会保留，新增的中英文地址字段为空
- 脚本会将原有地址内容复制到 `address_cn_detail` 字段
- 如需要自定义地址解析逻辑，请修改脚本中的数据迁移部分

### 性能考虑

- 脚本使用 `background: true` 创建索引，不会阻塞其他操作
- 大数据量的情况下，索引创建可能需要较长时间
- 可以通过 `db.currentOp()` 监控索引创建进度

## 🛠️ 自定义配置

### 修改默认值

如果需要为新字段设置不同的默认值，可以修改脚本中的 `$set` 部分：

```javascript
$set: {
  "company_code": "DEFAULT_CODE",  // 设置默认代码
  "submitted_by": "SYSTEM",        // 设置默认提交人
  // ... 其他字段
}
```

### 自定义地址解析

如果原有地址数据有特定格式，可以在数据迁移部分添加解析逻辑：

```javascript
// 示例：解析 "北京市朝阳区三环路123号" 格式的地址
if (doc.address && doc.address.includes("市")) {
  let addressParts = parseChineseAddress(doc.address);
  updateDoc["address_cn_province"] = addressParts.province;
  updateDoc["address_cn_city"] = addressParts.city;
  // ...
}
```

## 📞 技术支持

如果在执行过程中遇到问题：

1. 检查 MongoDB 连接和权限
2. 查看脚本输出的错误信息
3. 确认数据库名称和集合名称是否正确
4. 如有疑问，请保存完整的执行日志以便排查 