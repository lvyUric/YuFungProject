# 保险公司数据库集合文档 (companies)

## 📋 集合概述

`companies` 集合存储平台接入的保险经纪公司基本信息，包括基本信息、联系方式、地址信息、业务信息等。采用简化的字段结构，支持高效的公司管理功能。

## 🏗️ 字段结构定义

### 基本信息字段

| 字段名 | 数据类型 | 是否必填 | 索引类型 | 中文说明 | 备注 |
|--------|----------|----------|----------|----------|------|
| `_id` | ObjectId | 是 | 主键 | MongoDB主键ID | 系统自动生成 |
| `company_id` | String | 是 | 唯一索引 | 公司唯一标识 | 业务主键，格式：CMP+时间戳+随机数 |
| `company_name` | String | 是 | 普通索引 | 公司名称 | 保险经纪公司完整名称 |
| `company_code` | String | 否 | 普通索引 | 公司代码 | 内部公司代码，用于业务标识 |

### 负责人信息字段

| 字段名 | 数据类型 | 是否必填 | 索引类型 | 中文说明 | 备注 |
|--------|----------|----------|----------|----------|------|
| `contact_person` | String | 否 | 普通索引 | 联络人 | 日常联络人姓名 |

### 联系方式字段

| 字段名 | 数据类型 | 是否必填 | 索引类型 | 中文说明 | 备注 |
|--------|----------|----------|----------|----------|------|
| `tel_no` | String | 否 | 普通索引 | 固定电话 | 公司固定电话号码 |
| `mobile` | String | 否 | 普通索引 | 移动电话 | 公司移动电话号码 |
| `contact_phone` | String | 否 | 普通索引 | 联系电话 | 主要联系电话（兼容旧版本） |
| `email` | String | 是 | 普通索引 | 邮箱地址 | 公司官方邮箱地址 |

### 地址字段

| 字段名 | 数据类型 | 是否必填 | 索引类型 | 中文说明 | 备注 |
|--------|----------|----------|----------|----------|------|
| `address_cn_detail` | String | 否 | 普通索引 | 中文详细地址 | 公司中文详细地址 |
| `address_en_detail` | String | 否 | 普通索引 | 英文详细地址 | 公司英文详细地址 |
| `address` | String | 否 | 普通索引 | 地址 | 原有地址字段（兼容旧版本） |

### 业务信息字段

| 字段名 | 数据类型 | 是否必填 | 索引类型 | 中文说明 | 备注 |
|--------|----------|----------|----------|----------|------|
| `broker_code` | String | 否 | 普通索引 | 经纪人代码 | 合作经纪人识别代码 |
| `link` | String | 否 | 无 | 相关链接 | 公司官网或相关页面链接 |

### 登录信息字段

| 字段名 | 数据类型 | 是否必填 | 索引类型 | 中文说明 | 备注 |
|--------|----------|----------|----------|----------|------|
| `username` | String | 否 | 稀疏索引 | 用户名 | 公司登录用户名 |
| `password_hash` | String | 否 | 无 | 密码哈希值 | 加密后的密码，不返回给客户端 |

### 系统管理字段

| 字段名 | 数据类型 | 是否必填 | 索引类型 | 中文说明 | 备注 |
|--------|----------|----------|----------|----------|------|
| `valid_start_date` | Date | 是 | 无 | 有效期开始日期 | 公司合作有效期开始时间 |
| `valid_end_date` | Date | 是 | 普通索引 | 有效期结束日期 | 公司合作有效期结束时间 |
| `user_quota` | Number | 是 | 无 | 用户配额 | 允许创建的用户数量上限 |
| `current_user_count` | Number | 是 | 无 | 当前用户数量 | 已创建的用户数量 |
| `status` | String | 是 | 普通索引 | 状态 | active=有效, inactive=停用, expired=过期 |
| `remark` | String | 否 | 无 | 备注信息 | 其他备注说明 |
| `submitted_by` | String | 否 | 普通索引 | 提交人 | 创建该公司记录的操作员 |
| `created_at` | Date | 是 | 普通索引 | 创建时间 | 记录创建时间 |
| `updated_at` | Date | 是 | 无 | 更新时间 | 记录最后更新时间 |

## 📚 索引策略

### 单字段索引

```javascript
// 主要业务索引
db.companies.createIndex({ "company_id": 1 }, { "unique": true });
db.companies.createIndex({ "company_name": 1 }, { "unique": true });
db.companies.createIndex({ "company_code": 1 });
db.companies.createIndex({ "contact_person": 1 });
db.companies.createIndex({ "tel_no": 1 });
db.companies.createIndex({ "mobile": 1 });
db.companies.createIndex({ "email": 1 });
db.companies.createIndex({ "broker_code": 1 });
db.companies.createIndex({ "username": 1 }, { "sparse": true });
db.companies.createIndex({ "submitted_by": 1 });
db.companies.createIndex({ "status": 1 });
db.companies.createIndex({ "valid_end_date": 1 });
db.companies.createIndex({ "created_at": 1 });
```

### 复合索引

```javascript
// 常用查询组合索引
db.companies.createIndex({ "status": 1, "valid_end_date": 1 });
db.companies.createIndex({ "company_name": 1, "company_code": 1 });
```

## 📄 文档示例

### 完整公司记录示例

```javascript
{
  "_id": ObjectId("65f7b8e4c12345678901234a"),
  "company_id": "CMP_202401150001",
  "company_name": "中华保险经纪有限公司",
  "company_code": "ZHBX001",
  
  // 负责人信息
  "contact_person": "张经理",
  
  // 联系方式
  "tel_no": "010-88888888",
  "mobile": "13800138000",
  "contact_phone": "010-88888888",
  "email": "info@zhbx.com",
  
  // 地址信息
  "address_cn_detail": "北京市朝阳区建国门外大街1号国贸大厦A座15层",
  "address_en_detail": "15F, Tower A, CWTC, No.1 Jianguomenwai Avenue, Chaoyang District, Beijing",
  "address": "北京市朝阳区建国门外大街1号国贸大厦A座15层",
  
  // 业务信息
  "broker_code": "BRK001",
  "link": "https://www.zhbx.com",
  
  // 登录信息
  "username": "zhbx_admin",
  "password_hash": "$2a$10$...",
  
  // 系统字段
  "valid_start_date": ISODate("2024-01-01T00:00:00.000Z"),
  "valid_end_date": ISODate("2025-12-31T23:59:59.999Z"),
  "user_quota": 100,
  "current_user_count": 15,
  "status": "active",
  "remark": "重要合作伙伴",
  "submitted_by": "admin",
  "created_at": ISODate("2024-01-15T08:30:00.000Z"),
  "updated_at": ISODate("2024-01-15T10:45:00.000Z")
}
```

### 最小化公司记录示例

```javascript
{
  "_id": ObjectId("65f7b8e4c12345678901234b"),
  "company_id": "CMP_202401150002",
  "company_name": "安泰保险代理有限公司",
  "contact_phone": "021-88888888",
  "email": "info@antai.com",
  "valid_start_date": ISODate("2024-01-01T00:00:00.000Z"),
  "valid_end_date": ISODate("2024-12-31T23:59:59.999Z"),
  "user_quota": 50,
  "current_user_count": 0,
  "status": "active",
  "created_at": ISODate("2024-01-15T09:00:00.000Z"),
  "updated_at": ISODate("2024-01-15T09:00:00.000Z")
}
```

## 📈 查询模式说明

### 常用查询场景

1. **按公司名称查询**
   ```javascript
   db.companies.find({ "company_name": /中华/ });
   ```

2. **按联系人查询**
   ```javascript
   db.companies.find({ "contact_person": /张/ });
   ```

3. **按状态和有效期查询**
   ```javascript
   db.companies.find({
     "status": "active",
     "valid_end_date": { $gte: new Date() }
   });
   ```

4. **按邮箱查询**
   ```javascript
   db.companies.find({ "email": "info@zhbx.com" });
   ```

### 聚合统计示例

```javascript
// 按状态统计公司数量
db.companies.aggregate([
  {
    $group: {
      _id: "$status",
      count: { $sum: 1 },
      total_quota: { $sum: "$user_quota" },
      total_users: { $sum: "$current_user_count" }
    }
  },
  { $sort: { count: -1 } }
]);
```

## ⚠️ 注意事项

1. **数据完整性**: 
   - `company_id` 和 `company_name` 必须唯一
   - `email` 字段必须填写且格式正确
   - 有效期字段必须合理设置

2. **性能优化**:
   - 使用复合索引提高常用查询性能
   - 稀疏索引用于可选字段如 `username`
   - 定期清理过期数据

3. **安全性**:
   - `password_hash` 字段不应返回给客户端
   - 敏感信息需要加密存储

4. **兼容性**:
   - 保留 `address` 和 `contact_phone` 字段用于向后兼容
   - 新应用程序应使用详细地址字段 