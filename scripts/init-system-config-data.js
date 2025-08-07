// ========================================
// 系统配置测试数据初始化脚本
// ========================================

print("开始初始化系统配置测试数据...");

// 连接到数据库
db = db.getSiblingDB('yufung_admin');

// 获取当前时间
const now = new Date();

// 生成配置ID的函数
function generateConfigID() {
    const timestamp = Date.now().toString(36);
    const random = Math.random().toString(36).substr(2, 5);
    return 'CONFIG' + timestamp + random.toUpperCase();
}

// 测试公司ID（使用现有的公司ID）
const testCompanyId = "CMP1735967088DA82E1D9"; // 这个应该是现有的公司ID

// 清空现有的系统配置数据
try {
    db.system_configs.deleteMany({});
    print("已清空现有系统配置数据");
} catch (e) {
    print("清空数据时出错:", e.message);
}

// 初始化系统配置数据
const systemConfigs = [
    // 港分客户经理配置
    {
        config_id: generateConfigID(),
        config_type: "hk_manager",
        config_key: "manager_001",
        config_value: "张经理",
        display_name: "张经理",
        company_id: testCompanyId,
        sort_order: 1,
        status: "enable",
        remark: "港分资深客户经理",
        created_by: "admin",
        updated_by: "admin",
        created_at: now,
        updated_at: now
    },
    {
        config_id: generateConfigID(),
        config_type: "hk_manager",
        config_key: "manager_002",
        config_value: "李经理",
        display_name: "李经理",
        company_id: testCompanyId,
        sort_order: 2,
        status: "enable",
        remark: "港分高级客户经理",
        created_by: "admin",
        updated_by: "admin",
        created_at: now,
        updated_at: now
    },
    {
        config_id: generateConfigID(),
        config_type: "hk_manager",
        config_key: "manager_003",
        config_value: "王经理",
        display_name: "王经理",
        company_id: testCompanyId,
        sort_order: 3,
        status: "enable",
        remark: "港分专业客户经理",
        created_by: "admin",
        updated_by: "admin",
        created_at: now,
        updated_at: now
    },

    // 转介分行配置
    {
        config_id: generateConfigID(),
        config_type: "referral_branch",
        config_key: "branch_001",
        config_value: "中环分行",
        display_name: "中环分行",
        company_id: testCompanyId,
        sort_order: 1,
        status: "enable",
        remark: "香港中环核心商业区分行",
        created_by: "admin",
        updated_by: "admin",
        created_at: now,
        updated_at: now
    },
    {
        config_id: generateConfigID(),
        config_type: "referral_branch",
        config_key: "branch_002",
        config_value: "铜锣湾分行",
        display_name: "铜锣湾分行",
        company_id: testCompanyId,
        sort_order: 2,
        status: "enable",
        remark: "香港铜锣湾商业区分行",
        created_by: "admin",
        updated_by: "admin",
        created_at: now,
        updated_at: now
    },
    {
        config_id: generateConfigID(),
        config_type: "referral_branch",
        config_key: "branch_003",
        config_value: "尖沙咀分行",
        display_name: "尖沙咀分行",
        company_id: testCompanyId,
        sort_order: 3,
        status: "enable",
        remark: "香港尖沙咀旅游商业区分行",
        created_by: "admin",
        updated_by: "admin",
        created_at: now,
        updated_at: now
    },
    {
        config_id: generateConfigID(),
        config_type: "referral_branch",
        config_key: "branch_004",
        config_value: "深圳分行",
        display_name: "深圳分行",
        company_id: testCompanyId,
        sort_order: 4,
        status: "enable",
        remark: "深圳福田中心区分行",
        created_by: "admin",
        updated_by: "admin",
        created_at: now,
        updated_at: now
    },

    // 合作伙伴配置
    {
        config_id: generateConfigID(),
        config_type: "partner",
        config_key: "partner_001",
        config_value: "汇丰银行",
        display_name: "汇丰银行",
        company_id: testCompanyId,
        sort_order: 1,
        status: "enable",
        remark: "香港汇丰银行合作伙伴",
        created_by: "admin",
        updated_by: "admin",
        created_at: now,
        updated_at: now
    },
    {
        config_id: generateConfigID(),
        config_type: "partner",
        config_key: "partner_002",
        config_value: "渣打银行",
        display_name: "渣打银行",
        company_id: testCompanyId,
        sort_order: 2,
        status: "enable",
        remark: "香港渣打银行合作伙伴",
        created_by: "admin",
        updated_by: "admin",
        created_at: now,
        updated_at: now
    },
    {
        config_id: generateConfigID(),
        config_type: "partner",
        config_key: "partner_003",
        config_value: "恒生银行",
        display_name: "恒生银行",
        company_id: testCompanyId,
        sort_order: 3,
        status: "enable",
        remark: "香港恒生银行合作伙伴",
        created_by: "admin",
        updated_by: "admin",
        created_at: now,
        updated_at: now
    },
    {
        config_id: generateConfigID(),
        config_type: "partner",
        config_key: "partner_004",
        config_value: "中国银行",
        display_name: "中国银行（香港）",
        company_id: testCompanyId,
        sort_order: 4,
        status: "enable",
        remark: "中国银行香港分行合作伙伴",
        created_by: "admin",
        updated_by: "admin",
        created_at: now,
        updated_at: now
    },
    {
        config_id: generateConfigID(),
        config_type: "partner",
        config_key: "partner_005",
        config_value: "招商银行",
        display_name: "招商银行",
        company_id: testCompanyId,
        sort_order: 5,
        status: "enable",
        remark: "招商银行合作伙伴",
        created_by: "admin",
        updated_by: "admin",
        created_at: now,
        updated_at: now
    }
];

// 插入系统配置数据
try {
    const result = db.system_configs.insertMany(systemConfigs);
    print("✅ 成功插入 " + result.insertedIds.length + " 条系统配置记录");
    
    // 显示插入的数据统计
    const hkManagerCount = db.system_configs.countDocuments({ config_type: "hk_manager" });
    const referralBranchCount = db.system_configs.countDocuments({ config_type: "referral_branch" });
    const partnerCount = db.system_configs.countDocuments({ config_type: "partner" });
    
    print("\n📊 数据统计：");
    print("- 港分客户经理: " + hkManagerCount + " 条");
    print("- 转介分行: " + referralBranchCount + " 条");
    print("- 合作伙伴: " + partnerCount + " 条");
    print("- 总计: " + (hkManagerCount + referralBranchCount + partnerCount) + " 条");
    
} catch (e) {
    print("❌ 插入数据时出错:", e.message);
    throw e;
}

print("\n✅ 系统配置测试数据初始化完成！");
print("现在可以在前端系统配置管理页面中查看和管理这些配置项。"); 