// 初始化系统配置数据
use yufung_db;

// 删除现有的系统配置数据
db.system_configs.deleteMany({});

// 港分客户经理配置
db.system_configs.insertMany([
  {
    config_id: "CONFIG001",
    config_type: "hk_manager",
    config_key: "manager_zhang",
    config_value: "张经理",
    display_name: "张经理",
    company_id: "COMP001",
    sort_order: 1,
    status: "enable",
    remark: "港分客户经理",
    created_by: "system",
    updated_by: "system",
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    config_id: "CONFIG002",
    config_type: "hk_manager",
    config_key: "manager_li",
    config_value: "李经理",
    display_name: "李经理",
    company_id: "COMP001",
    sort_order: 2,
    status: "enable",
    remark: "港分客户经理",
    created_by: "system",
    updated_by: "system",
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    config_id: "CONFIG003",
    config_type: "hk_manager",
    config_key: "manager_wang",
    config_value: "王经理",
    display_name: "王经理",
    company_id: "COMP001",
    sort_order: 3,
    status: "enable",
    remark: "港分客户经理",
    created_by: "system",
    updated_by: "system",
    created_at: new Date(),
    updated_at: new Date()
  }
]);

// 转介分行配置
db.system_configs.insertMany([
  {
    config_id: "CONFIG004",
    config_type: "referral_branch",
    config_key: "branch_central",
    config_value: "中环分行",
    display_name: "中环分行",
    company_id: "COMP001",
    sort_order: 1,
    status: "enable",
    remark: "转介分行",
    created_by: "system",
    updated_by: "system",
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    config_id: "CONFIG005",
    config_type: "referral_branch",
    config_key: "branch_tsim",
    config_value: "尖沙咀分行",
    display_name: "尖沙咀分行",
    company_id: "COMP001",
    sort_order: 2,
    status: "enable",
    remark: "转介分行",
    created_by: "system",
    updated_by: "system",
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    config_id: "CONFIG006",
    config_type: "referral_branch",
    config_key: "branch_causeway",
    config_value: "铜锣湾分行",
    display_name: "铜锣湾分行",
    company_id: "COMP001",
    sort_order: 3,
    status: "enable",
    remark: "转介分行",
    created_by: "system",
    updated_by: "system",
    created_at: new Date(),
    updated_at: new Date()
  }
]);

// 合作伙伴配置
db.system_configs.insertMany([
  {
    config_id: "CONFIG007",
    config_type: "partner",
    config_key: "partner_bank_a",
    config_value: "银行A",
    display_name: "银行A",
    company_id: "COMP001",
    sort_order: 1,
    status: "enable",
    remark: "合作伙伴",
    created_by: "system",
    updated_by: "system",
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    config_id: "CONFIG008",
    config_type: "partner",
    config_key: "partner_insurance_b",
    config_value: "保险公司B",
    display_name: "保险公司B",
    company_id: "COMP001",
    sort_order: 2,
    status: "enable",
    remark: "合作伙伴",
    created_by: "system",
    updated_by: "system",
    created_at: new Date(),
    updated_at: new Date()
  },
  {
    config_id: "CONFIG009",
    config_type: "partner",
    config_key: "partner_financial_c",
    config_value: "金融机构C",
    display_name: "金融机构C",
    company_id: "COMP001",
    sort_order: 3,
    status: "enable",
    remark: "合作伙伴",
    created_by: "system",
    updated_by: "system",
    created_at: new Date(),
    updated_at: new Date()
  }
]);

print("系统配置初始化完成！");
print("港分客户经理配置: " + db.system_configs.count({config_type: "hk_manager"}) + " 条");
print("转介分行配置: " + db.system_configs.count({config_type: "referral_branch"}) + " 条");
print("合作伙伴配置: " + db.system_configs.count({config_type: "partner"}) + " 条"); 