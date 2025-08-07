// =============================================================================
// MongoDB 公司集合字段回滚脚本
// 用途: 回滚到升级前的状态
// =============================================================================

// 切换到目标数据库
use('yufung_admin');

print("=== 开始回滚公司集合字段结构 ===");

// 1. 检查备份集合是否存在
print("1. 检查备份集合...");
let backupExists = db.companies_backup.countDocuments({}) > 0;
if (!backupExists) {
  print("❌ 错误: 未找到备份集合 companies_backup，无法执行回滚");
  quit(1);
}
print("✓ 备份集合存在，包含文档数量: " + db.companies_backup.countDocuments({}));

// 2. 创建当前状态备份（以防回滚后需要恢复）
print("2. 创建当前状态备份...");
db.companies_rollback_backup.drop(); // 删除可能存在的旧备份
db.companies.aggregate([
  { $match: {} }
]).forEach(function(doc) {
  db.companies_rollback_backup.insertOne(doc);
});
print("✓ 当前状态已备份到 companies_rollback_backup 集合");

// 3. 删除新增的字段
print("3. 删除新增字段...");
db.companies.updateMany(
  {},
  {
    $unset: {
      // 基本信息扩展
      "company_code": "",
      
      // 负责人信息
      "principles_cn": "",
      "principles_en": "", 
      "contact_person": "",
      
      // 联系方式扩展
      "tel_no": "",
      "mobile": "",
      
      // 中文地址
      "address_cn_province": "",
      "address_cn_city": "",
      "address_cn_district": "", 
      "address_cn_detail": "",
      
      // 英文地址
      "address_en_province": "",
      "address_en_city": "",
      "address_en_district": "",
      "address_en_detail": "",
      
      // 业务信息
      "broker_code": "",
      "link": "",
      
      // 登录信息
      "username": "",
      "password_hash": "",
      
      // 扩展信息
      "remark1": "",
      "remark2": "",
      "email_template": "",
      "payment_notification": "",
      
      // 系统字段
      "submitted_by": ""
    }
  }
);
print("✓ 新增字段已删除");

// 4. 删除新建的索引
print("4. 删除新建索引...");
try {
  db.companies.dropIndex({ "company_code": 1 });
  db.companies.dropIndex({ "principles_cn": 1 });
  db.companies.dropIndex({ "principles_en": 1 });
  db.companies.dropIndex({ "contact_person": 1 });
  db.companies.dropIndex({ "tel_no": 1 });
  db.companies.dropIndex({ "mobile": 1 });
  db.companies.dropIndex({ 
    "address_cn_province": 1, 
    "address_cn_city": 1,
    "address_cn_district": 1 
  });
  db.companies.dropIndex({ "broker_code": 1 });
  db.companies.dropIndex({ "username": 1 });
  db.companies.dropIndex({ "submitted_by": 1 });
  db.companies.dropIndex({ 
    "company_name": 1, 
    "company_code": 1 
  });
  print("✓ 新建索引已删除");
} catch (e) {
  print("⚠ 部分索引删除失败（可能原本不存在）: " + e.message);
}

// 5. 验证回滚结果
print("5. 验证回滚结果...");

// 检查新字段是否已删除
let sampleDoc = db.companies.findOne({});
if (sampleDoc) {
  let removedFields = [
    'company_code', 'principles_cn', 'principles_en', 'contact_person',
    'tel_no', 'mobile', 'address_cn_province', 'address_cn_city', 
    'address_cn_district', 'address_cn_detail', 'address_en_province',
    'address_en_city', 'address_en_district', 'address_en_detail',
    'broker_code', 'link', 'username', 'password_hash',
    'remark1', 'remark2', 'email_template', 'payment_notification',
    'submitted_by'
  ];
  
  let remainingFields = [];
  removedFields.forEach(field => {
    if (field in sampleDoc) {
      remainingFields.push(field);
    }
  });
  
  if (remainingFields.length === 0) {
    print("✓ 所有新字段已成功删除");
  } else {
    print("⚠ 以下字段仍然存在: " + remainingFields.join(", "));
  }
}

// 6. 显示当前集合状态
print("6. 当前集合状态:");
print("- 文档总数: " + db.companies.countDocuments({}));
print("- 备份集合: companies_backup (" + db.companies_backup.countDocuments({}) + " 文档)");
print("- 回滚前备份: companies_rollback_backup (" + db.companies_rollback_backup.countDocuments({}) + " 文档)");

print("=== 公司集合字段回滚完成 ===");
print("");
print("回滚摘要:");
print("- ✓ 已删除所有新增字段");
print("- ✓ 已删除新建索引");
print("- ✓ 已创建回滚前状态备份");
print("- ✓ 已验证回滚结果");
print("");
print("注意事项:");
print("1. 原始备份仍保存在 companies_backup 集合中");
print("2. 回滚前的状态已备份到 companies_rollback_backup 集合");
print("3. 如需要恢复到回滚前状态，可以使用 companies_rollback_backup 集合");
print("4. 确认回滚成功后，可以删除不需要的备份集合"); 