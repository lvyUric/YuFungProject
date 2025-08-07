// MongoDB JavaScript脚本 - 创建保险经纪公司表（简化版）
// 使用方法: mongo your_database_name create_company_table_simple.js

print("🚀 创建保险经纪公司表结构...");

// 创建集合
db.createCollection("companies", {
    validator: {
        $jsonSchema: {
            bsonType: "object",
            required: ["company_id", "company_name", "email", "status", "created_at", "updated_at"],
            properties: {
                // 基本信息
                company_id: { bsonType: "string", description: "公司唯一标识，必填" },
                company_name: { bsonType: "string", description: "公司名称，必填" },
                company_code: { bsonType: "string", description: "内部公司代码" },
                
                // 负责人信息
                contact_person: { bsonType: "string", description: "联络人" },
                
                // 联系方式
                tel_no: { bsonType: "string", description: "固定电话" },
                mobile: { bsonType: "string", description: "移动电话" },
                contact_phone: { bsonType: "string", description: "联系电话（兼容字段）" },
                email: { bsonType: "string", description: "邮箱地址，必填" },
                
                // 中文地址信息
                address_cn_province: { bsonType: "string", description: "中文地址-省/自治区/直辖市" },
                address_cn_city: { bsonType: "string", description: "中文地址-市" },
                address_cn_district: { bsonType: "string", description: "中文地址-县/区" },
                address_cn_detail: { bsonType: "string", description: "中文地址-详细地址" },
                
                // 英文地址信息
                address_en_province: { bsonType: "string", description: "英文地址-省/自治区/直辖市" },
                address_en_city: { bsonType: "string", description: "英文地址-市" },
                address_en_district: { bsonType: "string", description: "英文地址-县/区" },
                address_en_detail: { bsonType: "string", description: "英文地址-详细地址" },
                
                address: { bsonType: "string", description: "原有地址字段（兼容）" },
                
                // 业务信息
                broker_code: { bsonType: "string", description: "经纪人代码" },
                link: { bsonType: "string", description: "相关链接" },
                
                // 登录信息
                username: { bsonType: "string", description: "用户名" },
                password_hash: { bsonType: "string", description: "密码哈希值" },
                
                // 系统字段
                valid_start_date: { bsonType: "date", description: "有效期开始日期" },
                valid_end_date: { bsonType: "date", description: "有效期结束日期" },
                user_quota: { bsonType: "int", description: "用户配额" },
                current_user_count: { bsonType: "int", description: "当前用户数量" },
                status: { 
                    bsonType: "string", 
                    enum: ["active", "inactive", "expired"],
                    description: "状态：active=有效, inactive=停用, expired=过期"
                },
                remark: { bsonType: "string", description: "备注信息" },
                submitted_by: { bsonType: "string", description: "提交人" },
                created_at: { bsonType: "date", description: "创建时间" },
                updated_at: { bsonType: "date", description: "更新时间" }
            }
        }
    }
});

print("✅ 集合创建完成");

// 创建核心索引
print("📋 创建索引...");

// 业务主键唯一索引
db.companies.createIndex({ "company_id": 1 }, { unique: true, name: "idx_company_id" });

// 公司名称唯一索引
db.companies.createIndex({ "company_name": 1 }, { unique: true, name: "idx_company_name" });

// 邮箱唯一索引
db.companies.createIndex({ "email": 1 }, { unique: true, sparse: true, name: "idx_email" });

// 用户名唯一索引
db.companies.createIndex({ "username": 1 }, { unique: true, sparse: true, name: "idx_username" });

// 状态查询索引
db.companies.createIndex({ "status": 1 }, { name: "idx_status" });

// 有效期查询索引
db.companies.createIndex({ "valid_start_date": 1, "valid_end_date": 1 }, { name: "idx_valid_period" });

// 创建时间索引
db.companies.createIndex({ "created_at": -1 }, { name: "idx_created_at" });

// 地址查询索引
db.companies.createIndex({ "address_cn_province": 1, "address_cn_city": 1 }, { name: "idx_address" });

// 文本搜索索引
db.companies.createIndex({
    "company_name": "text",
    "company_code": "text", 
    "contact_person": "text",
    "email": "text"
}, { name: "idx_text_search", default_language: "none" });

print("✅ 索引创建完成");

// 显示创建结果
print("\n📊 表结构创建完成:");
print("  - 集合名称: companies");
print("  - 字段验证: 已启用");
print("  - 索引数量: " + db.companies.getIndexes().length);

print("\n📋 主要索引:");
db.companies.getIndexes().forEach((index, i) => {
    print(`  ${i + 1}. ${index.name}`);
});

print("\n✅ 公司表创建完成，可以开始使用！"); 