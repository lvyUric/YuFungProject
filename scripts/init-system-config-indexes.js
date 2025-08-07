// ========================================
// 系统配置表索引初始化脚本
// ========================================

print("开始初始化系统配置表索引...");

// 连接到数据库
db = db.getSiblingDB('yufung_admin');

// 删除现有索引（除了_id索引）
try {
    db.system_configs.dropIndexes();
    print("已删除系统配置表现有索引");
} catch (e) {
    print("删除索引时出错（可能是第一次运行）:", e.message);
}

// 创建索引
try {
    // 1. 配置ID唯一索引
    db.system_configs.createIndex(
        { "config_id": 1 },
        { 
            unique: true,
            name: "idx_config_id_unique",
            background: true
        }
    );
    print("✅ 已创建配置ID唯一索引");

    // 2. 公司ID + 配置类型 + 配置键复合索引（确保同一公司同一类型下配置键唯一）
    db.system_configs.createIndex(
        { 
            "company_id": 1,
            "config_type": 1,
            "config_key": 1
        },
        { 
            unique: true,
            name: "idx_company_type_key_unique",
            background: true
        }
    );
    print("✅ 已创建公司+类型+键复合唯一索引");

    // 3. 公司ID + 配置类型 + 状态复合索引（用于获取启用的配置选项）
    db.system_configs.createIndex(
        { 
            "company_id": 1,
            "config_type": 1,
            "status": 1
        },
        { 
            name: "idx_company_type_status",
            background: true
        }
    );
    print("✅ 已创建公司+类型+状态复合索引");

    // 4. 排序索引
    db.system_configs.createIndex(
        { "sort_order": 1 },
        { 
            name: "idx_sort_order",
            background: true
        }
    );
    print("✅ 已创建排序索引");

    // 5. 创建时间索引
    db.system_configs.createIndex(
        { "created_at": -1 },
        { 
            name: "idx_created_at",
            background: true
        }
    );
    print("✅ 已创建创建时间索引");

    // 6. 文本搜索索引（用于关键词搜索）
    db.system_configs.createIndex(
        {
            "config_key": "text",
            "config_value": "text", 
            "display_name": "text"
        },
        {
            name: "idx_text_search",
            background: true,
            weights: {
                "display_name": 10,
                "config_value": 5,
                "config_key": 1
            }
        }
    );
    print("✅ 已创建文本搜索索引");

    print("✅ 系统配置表索引初始化完成！");

    // 显示所有索引
    print("\n当前系统配置表索引列表：");
    db.system_configs.getIndexes().forEach(function(index) {
        print("- " + index.name + ": " + JSON.stringify(index.key));
    });

} catch (e) {
    print("❌ 创建索引时出错:", e.message);
    throw e;
}

print("\n系统配置表索引初始化脚本执行完成！"); 