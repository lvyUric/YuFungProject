// MongoDB初始化脚本 - 变更记录集合
// 使用方法: mongosh --eval "load('scripts/init-change-records.js')"

print("========================================");
print("开始初始化变更记录集合...");
print("========================================");

// 切换到目标数据库
db = db.getSiblingDB('insurance_db');

// 变更记录集合名称
const CHANGE_RECORDS_COLLECTION = 'change_records';

try {
  // 检查集合是否已存在
  const collections = db.getCollectionNames();
  if (collections.includes(CHANGE_RECORDS_COLLECTION)) {
    print(`⚠️  集合 ${CHANGE_RECORDS_COLLECTION} 已存在，跳过创建`);
  } else {
    // 创建变更记录集合
    db.createCollection(CHANGE_RECORDS_COLLECTION, {
      validator: {
        $jsonSchema: {
          bsonType: "object",
          required: ["change_id", "table_name", "record_id", "user_id", "username", "company_id", "change_type", "changed_fields", "change_time"],
          properties: {
            change_id: {
              bsonType: "string",
              description: "变更记录唯一标识，必填"
            },
            table_name: {
              bsonType: "string",
              description: "表名，必填"
            },
            record_id: {
              bsonType: "string", 
              description: "记录ID，必填"
            },
            user_id: {
              bsonType: "string",
              description: "操作用户ID，必填"
            },
            username: {
              bsonType: "string",
              description: "用户名，必填"
            },
            company_id: {
              bsonType: "string",
              description: "所属公司ID，必填"
            },
            change_type: {
              bsonType: "string",
              enum: ["insert", "update", "delete"],
              description: "变更类型，必填"
            },
            old_values: {
              bsonType: "object",
              description: "变更前数据，可选"
            },
            new_values: {
              bsonType: "object", 
              description: "变更后数据，可选"
            },
            changed_fields: {
              bsonType: "array",
              items: {
                bsonType: "string"
              },
              description: "变更字段列表，必填"
            },
            change_time: {
              bsonType: "date",
              description: "变更时间，必填"
            },
            change_reason: {
              bsonType: "string",
              description: "变更原因，可选"
            },
            ip_address: {
              bsonType: "string",
              description: "IP地址，可选"
            },
            user_agent: {
              bsonType: "string",
              description: "浏览器信息，可选"
            }
          }
        }
      }
    });
    print(`✅ 变更记录集合 ${CHANGE_RECORDS_COLLECTION} 创建成功`);
  }

  // 创建索引
  print("\n开始创建索引...");
  
  const changeRecordsCollection = db.getCollection(CHANGE_RECORDS_COLLECTION);
  
  // 1. 唯一索引 - change_id
  try {
    changeRecordsCollection.createIndex(
      { change_id: 1 },
      { 
        unique: true,
        name: "idx_change_id_unique",
        background: true
      }
    );
    print("✅ 创建唯一索引: change_id");
  } catch (error) {
    print(`⚠️  索引 change_id 可能已存在: ${error.message}`);
  }

  // 2. 复合索引 - 表名和记录ID和时间（最重要的查询模式）
  try {
    changeRecordsCollection.createIndex(
      { table_name: 1, record_id: 1, change_time: -1 },
      {
        name: "idx_table_record_time",
        background: true
      }
    );
    print("✅ 创建复合索引: table_name + record_id + change_time");
  } catch (error) {
    print(`⚠️  索引可能已存在: ${error.message}`);
  }

  // 3. 复合索引 - 公司ID和时间
  try {
    changeRecordsCollection.createIndex(
      { company_id: 1, change_time: -1 },
      {
        name: "idx_company_time", 
        background: true
      }
    );
    print("✅ 创建复合索引: company_id + change_time");
  } catch (error) {
    print(`⚠️  索引可能已存在: ${error.message}`);
  }

  // 4. 复合索引 - 用户ID和时间
  try {
    changeRecordsCollection.createIndex(
      { user_id: 1, change_time: -1 },
      {
        name: "idx_user_time",
        background: true
      }
    );
    print("✅ 创建复合索引: user_id + change_time");
  } catch (error) {
    print(`⚠️  索引可能已存在: ${error.message}`);
  }

  // 5. 单字段索引 - 时间（用于清理旧记录）
  try {
    changeRecordsCollection.createIndex(
      { change_time: -1 },
      {
        name: "idx_change_time",
        background: true
      }
    );
    print("✅ 创建单字段索引: change_time");
  } catch (error) {
    print(`⚠️  索引可能已存在: ${error.message}`);
  }

  // 6. 复合索引 - 变更类型和时间
  try {
    changeRecordsCollection.createIndex(
      { change_type: 1, change_time: -1 },
      {
        name: "idx_type_time",
        background: true
      }
    );
    print("✅ 创建复合索引: change_type + change_time");
  } catch (error) {
    print(`⚠️  索引可能已存在: ${error.message}`);
  }

  // 显示所有索引
  print("\n📋 当前变更记录集合的索引:");
  const indexes = changeRecordsCollection.getIndexes();
  indexes.forEach((index, i) => {
    print(`${i + 1}. ${index.name}: ${JSON.stringify(index.key)}`);
  });

  // 插入一些测试数据（可选）
  print("\n🧪 插入测试数据...");
  const testRecord = {
    change_id: "CHG_TEST_" + new Date().getTime(),
    table_name: "policies",
    record_id: "POL_TEST_001", 
    user_id: "USR_ADMIN_001",
    username: "admin",
    company_id: "COM_TEST_001",
    change_type: "insert",
    old_values: {},
    new_values: {
      customer_name_cn: "测试客户",
      policy_currency: "USD",
      actual_premium: 10000
    },
    changed_fields: ["customer_name_cn", "policy_currency", "actual_premium"],
    change_time: new Date(),
    change_reason: "测试数据",
    ip_address: "127.0.0.1",
    user_agent: "Test Script"
  };

  const result = changeRecordsCollection.insertOne(testRecord);
  print(`✅ 插入测试记录成功，ID: ${result.insertedId}`);

  // 验证查询
  const count = changeRecordsCollection.countDocuments({});
  print(`📊 变更记录集合当前文档数量: ${count}`);

  print("\n========================================");
  print("✅ 变更记录集合初始化完成！");
  print("========================================");
  print("\n📝 使用说明:");
  print("1. 变更记录会在保单更新时自动创建");
  print("2. 默认保留所有历史记录，可通过定时任务清理旧数据");
  print("3. 查询API:");
  print("   - GET /api/policies/{id}/change-records - 获取保单变更记录");
  print("   - GET /api/change-records - 获取所有变更记录");
  print("4. 前端界面:");
  print("   - 保单详情页面 -> 变更记录标签页");

} catch (error) {
  print("❌ 初始化过程中发生错误:");
  print(error);
} 