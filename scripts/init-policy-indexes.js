// MongoDB保单集合索引初始化脚本

// 切换到项目数据库
db = db.getSiblingDB('insurance_db');

// 创建保单集合索引
print('正在创建保单集合索引...');

// 1. 业务主键索引
db.policies.createIndex({ "policy_id": 1 }, { unique: true, name: "idx_policy_id" });
print('创建保单ID唯一索引: idx_policy_id');

// 2. 公司隔离索引
db.policies.createIndex({ "company_id": 1 }, { name: "idx_company_id" });
print('创建公司ID索引: idx_company_id');

// 3. 序号索引（按公司分组）
db.policies.createIndex({ "company_id": 1, "serial_number": 1 }, { name: "idx_company_serial" });
print('创建公司序号复合索引: idx_company_serial');

// 4. 账户号索引（业务查询）
db.policies.createIndex({ "account_number": 1 }, { name: "idx_account_number" });
print('创建账户号索引: idx_account_number');

// 5. 客户号索引（业务查询）
db.policies.createIndex({ "customer_number": 1 }, { name: "idx_customer_number" });
print('创建客户号索引: idx_customer_number');

// 6. 投保单号索引（业务查询）
db.policies.createIndex({ "proposal_number": 1 }, { name: "idx_proposal_number" });
print('创建投保单号索引: idx_proposal_number');

// 7. 客户中文名索引（支持模糊查询）
db.policies.createIndex({ "customer_name_cn": "text" }, { name: "idx_customer_name_cn_text" });
print('创建客户中文名文本索引: idx_customer_name_cn_text');

// 8. 承保公司索引（业务筛选）
db.policies.createIndex({ "insurance_company": 1 }, { name: "idx_insurance_company" });
print('创建承保公司索引: idx_insurance_company');

// 9. 保单币种索引（业务筛选）
db.policies.createIndex({ "policy_currency": 1 }, { name: "idx_policy_currency" });
print('创建保单币种索引: idx_policy_currency');

// 10. 状态字段复合索引（业务统计）
db.policies.createIndex({ 
  "company_id": 1, 
  "is_surrendered": 1, 
  "past_cooling_period": 1, 
  "is_paid_commission": 1, 
  "is_employee": 1 
}, { name: "idx_status_fields" });
print('创建状态字段复合索引: idx_status_fields');

// 11. 日期范围查询索引
db.policies.createIndex({ "referral_date": 1 }, { name: "idx_referral_date" });
db.policies.createIndex({ "payment_date": 1 }, { name: "idx_payment_date" });
db.policies.createIndex({ "effective_date": 1 }, { name: "idx_effective_date" });
print('创建日期字段索引');

// 12. 创建时间索引（排序）
db.policies.createIndex({ "created_at": -1 }, { name: "idx_created_at_desc" });
print('创建创建时间索引: idx_created_at_desc');

// 13. 更新时间索引（排序）
db.policies.createIndex({ "updated_at": -1 }, { name: "idx_updated_at_desc" });
print('创建更新时间索引: idx_updated_at_desc');

// 14. 复合唯一索引：防止同一公司内重复的账户号和投保单号
db.policies.createIndex({ 
  "company_id": 1, 
  "account_number": 1 
}, { unique: true, name: "idx_company_account_unique" });
print('创建公司账户号唯一索引: idx_company_account_unique');

db.policies.createIndex({ 
  "company_id": 1, 
  "proposal_number": 1 
}, { unique: true, name: "idx_company_proposal_unique" });
print('创建公司投保单号唯一索引: idx_company_proposal_unique');

// 15. 金额字段索引（用于统计查询）
db.policies.createIndex({ "actual_premium": 1 }, { name: "idx_actual_premium" });
db.policies.createIndex({ "aum": 1 }, { name: "idx_aum" });
db.policies.createIndex({ "expected_fee": 1 }, { name: "idx_expected_fee" });
print('创建金额字段索引');

// 16. 创建人和更新人索引（审计查询）
db.policies.createIndex({ "created_by": 1 }, { name: "idx_created_by" });
db.policies.createIndex({ "updated_by": 1 }, { name: "idx_updated_by" });
print('创建操作人索引');

// 检查索引创建结果
print('\n=== 保单集合索引列表 ===');
var indexes = db.policies.getIndexes();
indexes.forEach(function(index) {
  print('索引名: ' + index.name + ', 字段: ' + JSON.stringify(index.key));
});

print('\n保单集合索引创建完成！');
print('总计创建索引数量: ' + indexes.length); 