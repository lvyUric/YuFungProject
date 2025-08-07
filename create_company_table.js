// MongoDB JavaScript脚本 - 创建保险经纪公司表
// 使用方法: mongo your_database_name create_company_table.js

// ===== 公司表创建脚本 =====
// 基于Go模型: internal/model/user.go 中的Company结构

print("🚀 开始创建保险经纪公司表...");

// 1. 创建集合（如果不存在）
const collectionName = "companies";

// 检查集合是否存在
const existingCollections = db.getCollectionNames();
if (existingCollections.includes(collectionName)) {
    print(`⚠️  集合 '${collectionName}' 已存在，将在现有集合上操作`);
} else {
    db.createCollection(collectionName);
    print(`✅ 成功创建集合 '${collectionName}'`);
}

// 2. 创建索引
print("📋 创建索引...");

try {
    // 业务主键唯一索引
    db.companies.createIndex(
        { "company_id": 1 }, 
        { 
            unique: true, 
            name: "idx_company_id_unique",
            background: true 
        }
    );
    print("✅ 创建公司ID唯一索引");

    // 公司名称唯一索引
    db.companies.createIndex(
        { "company_name": 1 }, 
        { 
            unique: true, 
            name: "idx_company_name_unique",
            background: true 
        }
    );
    print("✅ 创建公司名称唯一索引");

    // 邮箱唯一索引
    db.companies.createIndex(
        { "email": 1 }, 
        { 
            unique: true, 
            sparse: true,  // 允许空值
            name: "idx_email_unique",
            background: true 
        }
    );
    print("✅ 创建邮箱唯一索引");

    // 用户名唯一索引（如果存在）
    db.companies.createIndex(
        { "username": 1 }, 
        { 
            unique: true, 
            sparse: true,  // 允许空值
            name: "idx_username_unique",
            background: true 
        }
    );
    print("✅ 创建用户名唯一索引");

    // 状态查询索引
    db.companies.createIndex(
        { "status": 1 }, 
        { 
            name: "idx_status",
            background: true 
        }
    );
    print("✅ 创建状态查询索引");

    // 有效期查询复合索引
    db.companies.createIndex(
        { 
            "valid_start_date": 1, 
            "valid_end_date": 1 
        }, 
        { 
            name: "idx_valid_period",
            background: true 
        }
    );
    print("✅ 创建有效期查询索引");

    // 创建时间索引
    db.companies.createIndex(
        { "created_at": -1 }, 
        { 
            name: "idx_created_at_desc",
            background: true 
        }
    );
    print("✅ 创建时间排序索引");

    // 地址查询索引
    db.companies.createIndex(
        { 
            "address_cn_province": 1, 
            "address_cn_city": 1 
        }, 
        { 
            name: "idx_address_cn",
            background: true 
        }
    );
    print("✅ 创建中文地址查询索引");

    // 文本搜索索引（用于搜索公司名称、联系人等）
    db.companies.createIndex(
        {
            "company_name": "text",
            "company_code": "text",
            "contact_person": "text",
            "email": "text"
        },
        {
            name: "idx_text_search",
            background: true,
            default_language: "none"  // 支持中文搜索
        }
    );
    print("✅ 创建文本搜索索引");

} catch (error) {
    print("❌ 创建索引时出错:", error.message);
}

// 3. 插入示例数据
print("📝 插入示例数据...");

const currentTime = new Date();
const sampleCompanies = [
    {
        company_id: "COMP001",
        company_name: "中国平安保险经纪有限公司",
        company_code: "PA001",
        
        // 负责人信息
        contact_person: "张三",
        
        // 联系方式
        tel_no: "010-12345678",
        mobile: "13800138000",
        contact_phone: "13800138000", // 兼容字段
        email: "contact@pingan-broker.com",
        
        // 中文地址信息
        address_cn_province: "北京市",
        address_cn_city: "北京市",
        address_cn_district: "朝阳区",
        address_cn_detail: "建国门外大街88号",
        
        // 英文地址信息
        address_en_province: "Beijing",
        address_en_city: "Beijing",
        address_en_district: "Chaoyang District",
        address_en_detail: "88 Jianguomenwai Avenue",
        
        address: "北京市朝阳区建国门外大街88号", // 兼容字段
        
        // 业务信息
        broker_code: "PA-BROKER-001",
        link: "https://www.pingan-broker.com",
        
        // 登录信息
        username: "pingan_admin",
        password_hash: "$2a$10$example_hash_value_for_password", // 示例哈希值
        
        // 系统字段
        valid_start_date: new Date("2024-01-01"),
        valid_end_date: new Date("2025-12-31"),
        user_quota: 100,
        current_user_count: 0,
        status: "active",
        remark: "中国平安保险经纪公司 - 示例数据",
        submitted_by: "system",
        created_at: currentTime,
        updated_at: currentTime
    },
    {
        company_id: "COMP002", 
        company_name: "太平洋保险经纪有限公司",
        company_code: "CPIC001",
        
        // 负责人信息
        contact_person: "李四",
        
        // 联系方式
        tel_no: "021-87654321",
        mobile: "13900139000",
        contact_phone: "13900139000",
        email: "info@cpic-broker.com",
        
        // 中文地址信息
        address_cn_province: "上海市",
        address_cn_city: "上海市",
        address_cn_district: "浦东新区",
        address_cn_detail: "陆家嘴环路1000号",
        
        // 英文地址信息
        address_en_province: "Shanghai",
        address_en_city: "Shanghai", 
        address_en_district: "Pudong New Area",
        address_en_detail: "1000 Lujiazui Ring Road",
        
        address: "上海市浦东新区陆家嘴环路1000号",
        
        // 业务信息
        broker_code: "CPIC-BROKER-001",
        link: "https://www.cpic-broker.com",
        
        // 登录信息
        username: "cpic_admin",
        password_hash: "$2a$10$another_example_hash_value_for_password",
        
        // 系统字段
        valid_start_date: new Date("2024-01-01"),
        valid_end_date: new Date("2025-12-31"),
        user_quota: 50,
        current_user_count: 0,
        status: "active",
        remark: "太平洋保险经纪公司 - 示例数据",
        submitted_by: "system",
        created_at: currentTime,
        updated_at: currentTime
    },
    {
        company_id: "COMP003",
        company_name: "阳光保险经纪有限公司",
        company_code: "SUN001",
        
        // 负责人信息
        contact_person: "王五",
        
        // 联系方式
        tel_no: "0755-88888888",
        mobile: "13700137000",
        contact_phone: "13700137000",
        email: "service@sunshine-broker.com",
        
        // 中文地址信息
        address_cn_province: "广东省",
        address_cn_city: "深圳市",
        address_cn_district: "福田区",
        address_cn_detail: "深南中路2018号",
        
        // 英文地址信息
        address_en_province: "Guangdong Province",
        address_en_city: "Shenzhen",
        address_en_district: "Futian District", 
        address_en_detail: "2018 Shennan Middle Road",
        
        address: "广东省深圳市福田区深南中路2018号",
        
        // 业务信息
        broker_code: "SUN-BROKER-001",
        link: "https://www.sunshine-broker.com",
        
        // 登录信息
        username: "sunshine_admin",
        password_hash: "$2a$10$third_example_hash_value_for_password",
        
        // 系统字段
        valid_start_date: new Date("2024-01-01"),
        valid_end_date: new Date("2025-12-31"),
        user_quota: 30,
        current_user_count: 0,
        status: "inactive", // 示例：停用状态
        remark: "阳光保险经纪公司 - 示例数据（停用状态）",
        submitted_by: "system",
        created_at: currentTime,
        updated_at: currentTime
    }
];

try {
    // 插入示例数据
    const result = db.companies.insertMany(sampleCompanies);
    print(`✅ 成功插入 ${result.insertedIds.length} 条示例数据`);
    
    // 显示插入的数据ID
    print("📋 插入的文档ID:");
    result.insertedIds.forEach((id, index) => {
        print(`  ${index + 1}. ${id} - ${sampleCompanies[index].company_name}`);
    });
    
} catch (error) {
    if (error.code === 11000) {
        print("⚠️  示例数据已存在（重复键错误），跳过插入");
    } else {
        print("❌ 插入示例数据时出错:", error.message);
    }
}

// 4. 验证创建结果
print("\n🔍 验证表创建结果:");

// 显示集合状态
const stats = db.companies.stats();
print(`📊 集合统计:`);
print(`  - 文档数量: ${stats.count}`);
print(`  - 存储大小: ${Math.round(stats.size / 1024)} KB`);
print(`  - 索引数量: ${stats.indexSizes ? Object.keys(stats.indexSizes).length : 'N/A'}`);

// 显示索引信息
print(`📇 已创建的索引:`);
const indexes = db.companies.getIndexes();
indexes.forEach((index, i) => {
    print(`  ${i + 1}. ${index.name} - ${JSON.stringify(index.key)}`);
});

// 查询测试
print(`\n🧪 数据查询测试:`);
const activeCount = db.companies.countDocuments({ status: "active" });
const totalCount = db.companies.countDocuments({});
print(`  - 有效公司数量: ${activeCount}`);
print(`  - 总公司数量: ${totalCount}`);

// 显示一个示例文档结构
print(`\n📋 示例文档结构:`);
const sampleDoc = db.companies.findOne({}, { password_hash: 0 }); // 不显示密码哈希
if (sampleDoc) {
    print(JSON.stringify(sampleDoc, null, 2));
}

print("\n✅ 保险经纪公司表创建完成!");
print("\n💡 使用说明:");
print("  - 集合名称: companies");
print("  - 主要索引: company_id (唯一), company_name (唯一), email (唯一)");
print("  - 支持功能: 全文搜索、地址查询、状态筛选、有效期查询");
print("  - 示例查询:");
print("    db.companies.find({ status: 'active' });");
print("    db.companies.find({ $text: { $search: '平安' } });");
print("    db.companies.find({ 'address_cn_province': '北京市' });");

print("\n🔒 安全提醒:");
print("  - 密码哈希字段 (password_hash) 在查询时应排除");
print("  - 建议在生产环境中设置适当的访问权限");
print("  - 定期备份重要数据"); 