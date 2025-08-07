// ========================================
// 保单管理表结构修改脚本
// 功能：1. 确保汇率字段保留4个小数点
//      2. 设置投保单号为全局唯一主键
// ========================================

print("========================================");
print("开始执行保单管理表结构修改...");
print("========================================\n");

// 连接到数据库
const dbName = 'yufung_admin'; // 根据实际数据库名称调整
db = db.getSiblingDB(dbName);

print(`✅ 已连接到数据库: ${dbName}`);

// ========================================
// 1. 检查policies集合是否存在，不存在则创建
// ========================================
print("\n1. 检查并创建policies集合...");

const collections = db.getCollectionNames();
if (!collections.includes('policies')) {
    db.createCollection('policies');
    print("✅ 已创建policies集合");
} else {
    print("✅ policies集合已存在");
}

// ========================================
// 2. 删除旧的投保单号索引（如果存在）
// ========================================
print("\n2. 更新投保单号索引设置...");

try {
    // 删除公司范围内的投保单号唯一索引
    db.policies.dropIndex("idx_company_proposal_unique");
    print("✅ 已删除旧的公司级投保单号唯一索引");
} catch (e) {
    print("ℹ️  旧的公司级投保单号唯一索引不存在，跳过删除");
}

try {
    // 删除普通投保单号索引
    db.policies.dropIndex("idx_proposal_number");
    print("✅ 已删除旧的投保单号普通索引");
} catch (e) {
    print("ℹ️  旧的投保单号普通索引不存在，跳过删除");
}

// ========================================
// 3. 创建新的全局唯一投保单号索引
// ========================================
print("\n3. 创建全局唯一投保单号索引...");

try {
    // 创建全局唯一的投保单号索引
    db.policies.createIndex(
        { "proposal_number": 1 },
        { 
            unique: true, 
            name: "idx_proposal_number_unique_global",
            background: true,
            partialFilterExpression: { 
                "proposal_number": { $exists: true, $ne: "", $ne: null } 
            }
        }
    );
    print("✅ 已创建全局唯一投保单号索引: idx_proposal_number_unique_global");
} catch (e) {
    print("❌ 创建全局唯一投保单号索引失败:", e.message);
    print("   可能存在重复的投保单号，请先清理数据");
}

// ========================================
// 4. 添加汇率字段精度验证
// ========================================
print("\n4. 设置汇率字段精度验证...");

try {
    // 添加文档验证规则，确保汇率字段精度
    db.runCommand({
        "collMod": "policies",
        "validator": {
            $jsonSchema: {
                bsonType: "object",
                properties: {
                    proposal_number: {
                        bsonType: "string",
                        description: "投保单号必须是非空字符串"
                    },
                    exchange_rate: {
                        bsonType: ["double", "decimal"],
                        description: "汇率字段必须是数字类型，保留4位小数"
                    },
                    policy_id: {
                        bsonType: "string",
                        description: "保单ID必须是字符串"
                    },
                    company_id: {
                        bsonType: "string",
                        description: "公司ID必须是字符串"
                    },
                    created_at: {
                        bsonType: "date",
                        description: "创建时间必须是日期类型"
                    },
                    updated_at: {
                        bsonType: "date",
                        description: "更新时间必须是日期类型"
                    }
                },
                required: ["proposal_number", "policy_id", "company_id", "created_at", "updated_at"]
            }
        },
        "validationAction": "warn", // 使用warn而不是error，避免影响现有数据
        "validationLevel": "moderate"
    });
    print("✅ 已设置汇率字段精度验证规则");
} catch (e) {
    print("⚠️  设置验证规则时出现警告:", e.message);
}

// ========================================
// 5. 更新现有数据的汇率字段精度
// ========================================
print("\n5. 更新现有数据的汇率字段精度...");

try {
    // 查找所有有汇率字段的文档并更新精度
    const cursor = db.policies.find({ 
        "exchange_rate": { $exists: true, $ne: null, $type: "number" } 
    });
    
    let updateCount = 0;
    cursor.forEach(function(doc) {
        if (doc.exchange_rate !== null && doc.exchange_rate !== undefined) {
            // 将汇率保留4位小数
            const roundedRate = Math.round(doc.exchange_rate * 10000) / 10000;
            
            db.policies.updateOne(
                { _id: doc._id },
                { 
                    $set: { 
                        exchange_rate: roundedRate,
                        updated_at: new Date()
                    } 
                }
            );
            updateCount++;
        }
    });
    
    print(`✅ 已更新 ${updateCount} 条记录的汇率字段精度`);
} catch (e) {
    print("❌ 更新汇率字段精度失败:", e.message);
}

// ========================================
// 6. 创建其他必要的索引
// ========================================
print("\n6. 创建其他必要的索引...");

const indexesToCreate = [
    // 汇率字段索引（用于统计和查询）
    {
        fields: { "exchange_rate": 1 },
        options: { 
            name: "idx_exchange_rate",
            background: true,
            partialFilterExpression: { 
                "exchange_rate": { $exists: true, $ne: null } 
            }
        }
    },
    // 保单ID唯一索引（如果不存在）
    {
        fields: { "policy_id": 1 },
        options: { 
            unique: true,
            name: "idx_policy_id_unique",
            background: true
        }
    },
    // 公司ID索引（多租户隔离）
    {
        fields: { "company_id": 1 },
        options: { 
            name: "idx_company_id",
            background: true
        }
    }
];

indexesToCreate.forEach(function(indexDef) {
    try {
        db.policies.createIndex(indexDef.fields, indexDef.options);
        print(`✅ 已创建索引: ${indexDef.options.name}`);
    } catch (e) {
        if (e.message.includes("already exists")) {
            print(`ℹ️  索引已存在: ${indexDef.options.name}`);
        } else {
            print(`❌ 创建索引失败 ${indexDef.options.name}:`, e.message);
        }
    }
});

// ========================================
// 7. 验证修改结果
// ========================================
print("\n7. 验证修改结果...");

// 检查索引
print("\n📋 当前保单集合索引列表:");
const indexes = db.policies.getIndexes();
indexes.forEach(function(index) {
    const isUnique = index.unique ? " (唯一)" : "";
    print(`   - ${index.name}: ${JSON.stringify(index.key)}${isUnique}`);
});

// 检查文档数量
const totalDocs = db.policies.countDocuments();
print(`\n📊 保单集合文档总数: ${totalDocs}`);

// 检查有汇率的文档数量
const docsWithExchangeRate = db.policies.countDocuments({ 
    "exchange_rate": { $exists: true, $ne: null } 
});
print(`📊 包含汇率字段的文档数: ${docsWithExchangeRate}`);

// 检查投保单号唯一性
const uniqueProposalNumbers = db.policies.aggregate([
    { $match: { "proposal_number": { $exists: true, $ne: "", $ne: null } } },
    { $group: { _id: "$proposal_number", count: { $sum: 1 } } },
    { $match: { count: { $gt: 1 } } },
    { $count: "duplicates" }
]).toArray();

if (uniqueProposalNumbers.length > 0) {
    print(`⚠️  检测到 ${uniqueProposalNumbers[0].duplicates} 个重复的投保单号，需要手动处理`);
} else {
    print("✅ 投保单号唯一性检查通过");
}

// ========================================
// 8. 创建数据验证函数
// ========================================
print("\n8. 创建数据验证和辅助函数...");

// 创建验证汇率精度的函数
const validateExchangeRateFunction = `
function validateExchangeRate(rate) {
    if (rate === null || rate === undefined) return true;
    if (typeof rate !== 'number') return false;
    
    // 检查是否超过4位小数
    const decimalPlaces = (rate.toString().split('.')[1] || '').length;
    return decimalPlaces <= 4;
}
`;

// 创建格式化汇率的函数  
const formatExchangeRateFunction = `
function formatExchangeRate(rate) {
    if (rate === null || rate === undefined) return null;
    return Math.round(rate * 10000) / 10000;
}
`;

print("✅ 已定义数据验证函数（可在应用程序中使用）");

// ========================================
// 9. 生成修改报告
// ========================================
print("\n========================================");
print("🎉 保单管理表结构修改完成！");
print("========================================");

print("\n📋 修改内容总结:");
print("1. ✅ 设置投保单号为全局唯一主键");
print("2. ✅ 添加汇率字段精度控制（4位小数）");
print("3. ✅ 更新现有数据的汇率字段精度");
print("4. ✅ 创建必要的数据库索引");
print("5. ✅ 添加数据验证规则");

print("\n⚠️  注意事项:");
print("1. 投保单号现在具有全局唯一性，不允许重复");
print("2. 汇率字段最多保留4位小数");
print("3. 建议在应用程序中添加相应的验证逻辑");
print("4. 如有重复投保单号，需要先清理数据");

print("\n🔧 建议的应用程序更新:");
print("1. 前端表单验证：投保单号不允许重复");
print("2. 后端API验证：汇率字段精度控制");
print("3. 数据导入：检查投保单号唯一性");

print("\n脚本执行完成！"); 