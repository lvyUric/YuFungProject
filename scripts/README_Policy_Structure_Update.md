# 保单管理表结构修改脚本使用说明

## 📋 概述

本脚本 `update-policy-structure.js` 用于修改保单管理系统的数据库表结构，主要实现以下功能：

1. **汇率字段精度控制**：确保汇率字段保留4位小数
2. **投保单号全局唯一性**：设置投保单号为全局唯一主键
3. **数据验证规则**：添加字段验证和约束
4. **索引优化**：创建必要的数据库索引

## 🎯 修改内容

### 1. 汇率字段优化
- ✅ 设置汇率字段最多保留4位小数
- ✅ 更新现有数据的汇率精度
- ✅ 添加汇率字段验证规则
- ✅ 创建汇率字段索引

### 2. 投保单号唯一性
- ✅ 删除旧的公司级投保单号唯一约束
- ✅ 创建全局唯一投保单号索引
- ✅ 添加重复数据检测功能

### 3. 数据验证
- ✅ 添加MongoDB文档验证规则
- ✅ 设置必填字段约束
- ✅ 类型验证和格式校验

## 🚀 使用方法

### 方法一：MongoDB Shell 执行（推荐）

```bash
# 1. 确保MongoDB服务正在运行
sudo systemctl status mongod

# 2. 连接到MongoDB
mongo

# 3. 执行脚本
load("scripts/update-policy-structure.js")
```

### 方法二：命令行直接执行

```bash
# Windows PowerShell（需要注意路径分隔符）
mongo yufung_admin scripts/update-policy-structure.js

# Linux/macOS
mongo yufung_admin scripts/update-policy-structure.js
```

### 方法三：MongoDB Compass / Studio 3T

1. 连接到数据库 `yufung_admin`
2. 打开 JavaScript/MongoDB Shell 窗口
3. 复制 `update-policy-structure.js` 文件内容
4. 粘贴并执行

## 📊 执行前准备

### 1. 数据库备份（重要！）

```bash
# 备份整个数据库
mongodump --db yufung_admin --out ./backup/$(date +%Y%m%d_%H%M%S)

# 仅备份保单集合
mongodump --db yufung_admin --collection policies --out ./backup/policies_$(date +%Y%m%d_%H%M%S)
```

### 2. 检查现有数据

```javascript
// 连接到数据库
use yufung_admin;

// 检查保单数量
db.policies.countDocuments();

// 检查重复投保单号
db.policies.aggregate([
    { $match: { "proposal_number": { $exists: true, $ne: "", $ne: null } } },
    { $group: { _id: "$proposal_number", count: { $sum: 1 } } },
    { $match: { count: { $gt: 1 } } }
]);

// 检查汇率字段
db.policies.find({ "exchange_rate": { $exists: true } }).limit(5);
```

## ⚠️ 注意事项

### 执行前注意
1. **务必备份数据库**，避免数据丢失
2. **检查重复投保单号**，如有重复需先清理
3. **确认数据库名称**，脚本中默认使用 `yufung_admin`
4. **评估执行时间**，大量数据可能需要较长时间

### 执行后验证
1. **检查索引创建**：确认所有索引都成功创建
2. **验证数据完整性**：确认数据没有丢失
3. **测试应用功能**：确保前后端功能正常
4. **监控性能**：观察查询性能是否改善

## 🛠️ 故障排除

### 常见问题

#### 1. 投保单号重复错误
```
❌ 创建全局唯一投保单号索引失败: E11000 duplicate key error
```

**解决方案：**
```javascript
// 查找重复的投保单号
db.policies.aggregate([
    { $match: { "proposal_number": { $exists: true, $ne: "", $ne: null } } },
    { $group: { _id: "$proposal_number", count: { $sum: 1 }, docs: { $push: "$_id" } } },
    { $match: { count: { $gt: 1 } } }
]);

// 手动清理重复数据（根据具体情况选择保留哪条记录）
```

#### 2. 权限不足错误
```
❌ not authorized to execute command
```

**解决方案：**
```bash
# 使用管理员账户连接
mongo -u admin -p --authenticationDatabase admin yufung_admin
```

#### 3. 数据库连接失败
```
❌ couldn't connect to server
```

**解决方案：**
```bash
# 检查MongoDB服务状态
sudo systemctl status mongod

# 启动MongoDB服务
sudo systemctl start mongod
```

## 📈 性能优化建议

### 1. 索引使用建议
- 投保单号查询使用新的全局唯一索引
- 汇率相关统计使用汇率字段索引
- 复合查询考虑创建复合索引

### 2. 查询优化
```javascript
// 高效的投保单号查询
db.policies.findOne({ "proposal_number": "P123456" });

// 汇率范围查询
db.policies.find({ 
    "exchange_rate": { $gte: 1.0, $lte: 2.0 } 
});
```

## 🔧 应用程序更新建议

### 1. 后端API修改

```go
// 投保单号唯一性验证
func validateProposalNumber(proposalNumber string) error {
    count, err := policyCollection.CountDocuments(context.TODO(), 
        bson.M{"proposal_number": proposalNumber})
    if err != nil {
        return err
    }
    if count > 0 {
        return errors.New("投保单号已存在")
    }
    return nil
}

// 汇率精度控制
func formatExchangeRate(rate float64) float64 {
    return math.Round(rate*10000) / 10000
}
```

### 2. 前端表单验证

```javascript
// 投保单号重复检查
const checkProposalNumber = async (proposalNumber) => {
    const response = await api.get(`/policies/check-proposal/${proposalNumber}`);
    return !response.data.exists;
};

// 汇率格式化
const formatExchangeRate = (rate) => {
    return Math.round(rate * 10000) / 10000;
};
```

## 📞 技术支持

如遇到问题，请检查：
1. MongoDB版本兼容性（建议4.4+）
2. 用户权限设置
3. 数据库连接配置
4. 日志文件错误信息

建议在测试环境先执行，确认无误后再在生产环境使用。 