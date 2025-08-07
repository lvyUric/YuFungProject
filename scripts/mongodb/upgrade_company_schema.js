// =============================================================================
// MongoDB 公司集合字段升级脚本
// 版本: v2.1 - 简化版公司管理字段支持
// 创建时间: 2024
// =============================================================================

// 切换到目标数据库（请根据实际情况修改数据库名）
use('insurance_db');

print("=== 开始升级公司集合字段结构 ===");

// 1. 备份现有数据（可选，建议先执行）
print("1. 创建集合备份...");
db.companies.aggregate([
  { $match: {} }
]).forEach(function(doc) {
  db.companies_backup.insertOne(doc);
});
print("✓ 备份完成，备份集合名: companies_backup");

// 2. 为现有文档添加新字段
print("2. 为现有公司文档添加新字段...");

db.companies.updateMany(
  {}, // 匹配所有文档
  {
    $set: {
      // 基本信息扩展
      "company_code": "",

      // 负责人信息
      "contact_person": "",

      // 联系方式扩展
      "tel_no": "",
      "mobile": "",
      "email": "",
      // contact_phone 保留原有字段

      // 中文地址
      "address_cn_detail": "",

      // 英文地址
      "address_en_detail": "",
      // address 保留原有字段

      // 业务信息
      "broker_code": "",
      "link": "",

      // 登录信息
      "username": "",
      "password_hash": "",

      // 系统字段
      "submitted_by": ""
      // valid_start_date, valid_end_date, user_quota, current_user_count 等保留原有字段
      // status, remark, created_at, updated_at 等保留原有字段
    }
  }
);

print("✓ 字段添加完成");

// 3. 数据迁移和清理（将原有address拆分到新字段，如果需要的话）
print("3. 执行数据迁移...");

// 示例：如果原有address字段包含完整地址，可以尝试解析
db.companies.find({}).forEach(function(doc) {
  if (doc.address && doc.address.trim() !== "") {
    // 如果需要，可以在这里添加地址解析逻辑
    // 例如：将 "北京市朝阳区三环路123号" 拆分为详细地址

    // 简单示例（根据实际数据格式调整）
    let updateDoc = {};

    // 如果原有地址不为空，并且新的中文详细地址为空，则将原地址放到详细地址中
    if (!doc.address_cn_detail || doc.address_cn_detail === "") {
      updateDoc["address_cn_detail"] = doc.address;
    }

    // 只有在有更新内容时才执行更新
    if (Object.keys(updateDoc).length > 0) {
      db.companies.updateOne(
        { _id: doc._id },
        { $set: updateDoc }
      );
    }
  }
});

print("✓ 数据迁移完成");

// 4. 创建必要的索引
print("4. 创建新字段索引...");

// 公司代码索引
db.companies.createIndex({ "company_code": 1 }, { "background": true });

// 联络人索引
db.companies.createIndex({ "contact_person": 1 }, { "background": true });

// 联系方式索引
db.companies.createIndex({ "tel_no": 1 }, { "background": true });
db.companies.createIndex({ "mobile": 1 }, { "background": true });
db.companies.createIndex({ "email": 1 }, { "background": true });

// 经纪人代码索引
db.companies.createIndex({ "broker_code": 1 }, { "background": true });

// 用户名索引（如果作为登录使用）
db.companies.createIndex({ "username": 1 }, { "background": true, "sparse": true });

// 提交人索引
db.companies.createIndex({ "submitted_by": 1 }, { "background": true });

// 复合索引（常用查询组合）
db.companies.createIndex({
  "status": 1,
  "valid_end_date": 1
}, { "background": true });

db.companies.createIndex({
  "company_name": 1,
  "company_code": 1
}, { "background": true });

print("✓ 索引创建完成");

// 5. 验证升级结果
print("5. 验证升级结果...");

// 统计文档数量
let totalCount = db.companies.countDocuments({});
print(`总公司数量: ${totalCount}`);

// 检查新字段是否存在
let sampleDoc = db.companies.findOne({});
if (sampleDoc) {
  let newFields = [
    'company_code', 'contact_person', 'tel_no', 'mobile', 'email',
    'address_cn_detail', 'address_en_detail', 'broker_code', 'link',
    'username', 'password_hash', 'submitted_by'
  ];

  let missingFields = [];
  newFields.forEach(field => {
    if (!(field in sampleDoc)) {
      missingFields.push(field);
    }
  });

  if (missingFields.length === 0) {
    print("✓ 所有新字段已成功添加");
  } else {
    print("⚠ 以下字段未找到: " + missingFields.join(", "));
  }
}

// 6. 显示索引信息
print("6. 当前索引列表:");
db.companies.getIndexes().forEach(function(index) {
  print(`- ${index.name}: ${JSON.stringify(index.key)}`);
});

print("=== 公司集合字段升级完成 ===");
print("");
print("升级摘要:");
print("- ✓ 已备份原始数据到 companies_backup 集合");
print("- ✓ 已为所有现有文档添加新字段");
print("- ✓ 已执行数据迁移");
print("- ✓ 已创建必要的索引");
print("- ✓ 已验证升级结果");
print("");
print("注意事项:");
print("1. 如果需要回滚，可以使用 companies_backup 集合恢复数据");
print("2. 新字段已设置为空字符串默认值，可通过应用程序界面填充");
print("3. 建议在生产环境执行前先在测试环境验证");
print("4. 如有自定义的地址解析需求，请修改步骤3中的数据迁移逻辑");
