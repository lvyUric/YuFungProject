# 保险经纪公司表创建脚本

基于Go模型 `internal/model/user.go` 中的 `Company` 结构创建的MongoDB表结构。

## 📁 文件说明

### 1. `create_company_table.js` - 完整版脚本
- **功能**: 创建完整的公司表结构、索引和示例数据
- **包含**: 集合创建、索引设置、数据验证、示例数据插入
- **适合**: 开发环境、测试环境初始化

### 2. `create_company_table_simple.js` - 简化版脚本  
- **功能**: 仅创建表结构和核心索引
- **包含**: 集合创建、字段验证、必要索引
- **适合**: 生产环境部署

## 🚀 使用方法

### 方法一：MongoDB Shell 执行
```bash
# 连接到MongoDB数据库
mongo your_database_name

# 执行完整版脚本（包含示例数据）
load("create_company_table.js")

# 或执行简化版脚本（仅表结构）
load("create_company_table_simple.js")
```

### 方法二：命令行直接执行
```bash
# 执行完整版脚本
mongo your_database_name create_company_table.js

# 执行简化版脚本
mongo your_database_name create_company_table_simple.js
```

### 方法三：MongoDB Compass / Studio 3T
1. 连接到数据库
2. 打开JavaScript执行窗口
3. 复制脚本内容并执行

## 📋 表结构说明

### 集合名称
```
companies
```

### 字段结构
根据Go模型映射的完整字段：

```javascript
{
  // MongoDB主键
  _id: ObjectId,
  
  // 基本信息
  company_id: String,        // 公司唯一标识（业务主键）
  company_name: String,      // 公司名称
  company_code: String,      // 内部公司代码
  
  // 负责人信息
  contact_person: String,    // 联络人
  
  // 联系方式
  tel_no: String,           // 固定电话
  mobile: String,           // 移动电话
  contact_phone: String,    // 联系电话（兼容字段）
  email: String,            // 邮箱地址
  
  // 中文地址信息
  address_cn_province: String,  // 省/自治区/直辖市
  address_cn_city: String,      // 市
  address_cn_district: String,  // 县/区
  address_cn_detail: String,    // 详细地址
  
  // 英文地址信息
  address_en_province: String,  // Province/State
  address_en_city: String,      // City
  address_en_district: String,  // District
  address_en_detail: String,    // Detailed Address
  
  address: String,              // 原有地址字段（兼容）
  
  // 业务信息
  broker_code: String,      // 经纪人代码
  link: String,            // 相关链接
  
  // 登录信息
  username: String,         // 用户名
  password_hash: String,    // 密码哈希值
  
  // 系统字段
  valid_start_date: Date,   // 有效期开始日期
  valid_end_date: Date,     // 有效期结束日期
  user_quota: Number,       // 用户配额
  current_user_count: Number, // 当前用户数量
  status: String,           // 状态：active/inactive/expired
  remark: String,           // 备注信息
  submitted_by: String,     // 提交人
  created_at: Date,         // 创建时间
  updated_at: Date          // 更新时间
}
```

## 🔍 索引说明

### 唯一索引
- `company_id` - 公司唯一标识
- `company_name` - 公司名称 
- `email` - 邮箱地址（sparse，允许空值）
- `username` - 用户名（sparse，允许空值）

### 查询索引
- `status` - 状态查询
- `valid_start_date + valid_end_date` - 有效期查询
- `created_at` - 创建时间排序
- `address_cn_province + address_cn_city` - 地址查询

### 文本搜索索引
- `company_name + company_code + contact_person + email` - 全文搜索

## 📝 常用查询示例

### 基本查询
```javascript
// 查询所有有效公司
db.companies.find({ status: "active" });

// 按公司名称查询
db.companies.findOne({ company_name: "中国平安保险经纪有限公司" });

// 按地区查询
db.companies.find({ 
  address_cn_province: "北京市", 
  address_cn_city: "北京市" 
});
```

### 文本搜索
```javascript
// 搜索包含"平安"的公司
db.companies.find({ $text: { $search: "平安" } });

// 搜索多个关键词
db.companies.find({ $text: { $search: "平安 保险" } });
```

### 有效期查询
```javascript
// 查询当前有效的公司
const now = new Date();
db.companies.find({
  valid_start_date: { $lte: now },
  valid_end_date: { $gte: now },
  status: "active"
});

// 查询即将过期的公司（30天内）
const thirtyDaysLater = new Date(Date.now() + 30 * 24 * 60 * 60 * 1000);
db.companies.find({
  valid_end_date: { $lte: thirtyDaysLater },
  status: "active"
});
```

### 聚合查询
```javascript
// 按省份统计公司数量
db.companies.aggregate([
  { $group: { 
    _id: "$address_cn_province", 
    count: { $sum: 1 } 
  }},
  { $sort: { count: -1 }}
]);

// 统计不同状态的公司数量
db.companies.aggregate([
  { $group: { 
    _id: "$status", 
    count: { $sum: 1 } 
  }}
]);
```

## 🔒 安全注意事项

### 密码安全
```javascript
// 查询时排除密码哈希字段
db.companies.find({}, { password_hash: 0 });

// 更新时不要直接操作密码哈希
// 应该通过应用程序的密码加密逻辑处理
```

### 数据验证
- 脚本包含字段验证规则，确保数据完整性
- 必填字段：`company_id`, `company_name`, `email`, `status`, `created_at`, `updated_at`
- 状态字段限制：只能是 `active`, `inactive`, `expired`

## 🛠️ 维护操作

### 重建索引
```javascript
// 重建所有索引
db.companies.reIndex();

// 查看索引使用情况
db.companies.getIndexes();
```

### 数据备份
```bash
# 备份公司表
mongodump --db your_database_name --collection companies --out ./backup

# 恢复数据
mongorestore --db your_database_name --collection companies ./backup/your_database_name/companies.bson
```

### 性能优化
```javascript
// 查看集合统计信息
db.companies.stats();

// 分析查询性能
db.companies.find({ status: "active" }).explain("executionStats");
```

## ⚡ 故障排除

### 常见错误

1. **重复键错误 (E11000)**
   - 原因：违反唯一索引约束
   - 解决：检查 company_id、company_name、email、username 是否重复

2. **字段验证错误**
   - 原因：必填字段缺失或数据类型不匹配
   - 解决：确保必填字段完整，日期字段使用 Date 类型

3. **索引创建失败**
   - 原因：现有数据不符合索引要求
   - 解决：清理不符合要求的数据后重新创建索引

### 删除重建
```javascript
// 如果需要完全重建表
db.companies.drop();

// 然后重新执行创建脚本
load("create_company_table.js");
```

## 📞 技术支持

如果在使用过程中遇到问题，请检查：
1. MongoDB版本兼容性（建议4.0+）
2. 数据库连接权限
3. 脚本执行日志
4. Go应用程序的模型定义是否与脚本一致 