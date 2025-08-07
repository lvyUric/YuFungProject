// 活动记录表初始化脚本
// 使用方法：在MongoDB shell中执行 db.load("scripts/init_activity_log.js")

// 切换到目标数据库
use insurance_db

// 删除已存在的活动记录集合（如果存在）
db.activity_logs.drop()

// 创建活动记录集合
db.createCollection("activity_logs")

// 创建索引以提高查询性能
db.activity_logs.createIndex({ "user_id": 1 })
db.activity_logs.createIndex({ "company_id": 1 })
db.activity_logs.createIndex({ "operation_time": -1 })
db.activity_logs.createIndex({ "operation_type": 1 })
db.activity_logs.createIndex({ "module_name": 1 })
db.activity_logs.createIndex({ "company_id": 1, "operation_time": -1 })

// 插入示例活动记录数据
db.activity_logs.insertMany([
  {
    "log_id": "AL" + new Date().getTime() + "001",
    "user_id": "admin",
    "username": "admin",
    "company_id": "",
    "company_name": "平台管理",
    "operation_type": "login",
    "module_name": "认证授权",
    "operation_desc": "用户登录系统",
    "request_url": "/api/v1/auth/login",
    "request_method": "POST",
    "request_params": {},
    "ip_address": "192.168.1.100",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    "operation_time": new Date("2024-03-15T09:30:00Z"),
    "execution_time": 150,
    "result_status": "success",
    "target_id": "",
    "target_name": ""
  },
  {
    "log_id": "AL" + new Date().getTime() + "002",
    "user_id": "admin",
    "username": "admin",
    "company_id": "",
    "company_name": "平台管理",
    "operation_type": "create",
    "module_name": "公司管理",
    "operation_desc": "创建新公司：测试公司A",
    "request_url": "/api/v1/companies",
    "request_method": "POST",
    "request_params": {
      "company_name": "测试公司A",
      "contact_phone": "13800138000"
    },
    "ip_address": "192.168.1.100",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    "operation_time": new Date("2024-03-15T10:15:00Z"),
    "execution_time": 320,
    "result_status": "success",
    "target_id": "COMP001",
    "target_name": "测试公司A"
  },
  {
    "log_id": "AL" + new Date().getTime() + "003",
    "user_id": "user001",
    "username": "zhangsan",
    "company_id": "COMP001",
    "company_name": "测试公司A",
    "operation_type": "create",
    "module_name": "保单管理",
    "operation_desc": "创建新保单：A000000001",
    "request_url": "/api/v1/policies",
    "request_method": "POST",
    "request_params": {
      "proposal_number": "A000000001",
      "customer_name_cn": "张三"
    },
    "ip_address": "192.168.1.101",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    "operation_time": new Date("2024-03-15T11:20:00Z"),
    "execution_time": 450,
    "result_status": "success",
    "target_id": "POL001",
    "target_name": "A000000001"
  },
  {
    "log_id": "AL" + new Date().getTime() + "004",
    "user_id": "user001",
    "username": "zhangsan",
    "company_id": "COMP001",
    "company_name": "测试公司A",
    "operation_type": "update",
    "module_name": "保单管理",
    "operation_desc": "更新保单信息：A000000001",
    "request_url": "/api/v1/policies/POL001",
    "request_method": "PUT",
    "request_params": {
      "customer_name_cn": "张三",
      "policy_currency": "USD"
    },
    "ip_address": "192.168.1.101",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    "operation_time": new Date("2024-03-15T14:30:00Z"),
    "execution_time": 280,
    "result_status": "success",
    "target_id": "POL001",
    "target_name": "A000000001"
  },
  {
    "log_id": "AL" + new Date().getTime() + "005",
    "user_id": "admin",
    "username": "admin",
    "company_id": "",
    "company_name": "平台管理",
    "operation_type": "create",
    "module_name": "用户管理",
    "operation_desc": "创建新用户：lisi",
    "request_url": "/api/v1/users",
    "request_method": "POST",
    "request_params": {
      "username": "lisi",
      "display_name": "李四"
    },
    "ip_address": "192.168.1.100",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    "operation_time": new Date("2024-03-15T16:45:00Z"),
    "execution_time": 380,
    "result_status": "success",
    "target_id": "USER002",
    "target_name": "lisi"
  },
  {
    "log_id": "AL" + new Date().getTime() + "006",
    "user_id": "user002",
    "username": "lisi",
    "company_id": "COMP001",
    "company_name": "测试公司A",
    "operation_type": "view",
    "module_name": "保单管理",
    "operation_desc": "查看保单列表",
    "request_url": "/api/v1/policies",
    "request_method": "GET",
    "request_params": {
      "page": 1,
      "page_size": 20
    },
    "ip_address": "192.168.1.102",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    "operation_time": new Date("2024-03-15T17:10:00Z"),
    "execution_time": 120,
    "result_status": "success",
    "target_id": "",
    "target_name": ""
  },
  {
    "log_id": "AL" + new Date().getTime() + "007",
    "user_id": "user001",
    "username": "zhangsan",
    "company_id": "COMP001",
    "company_name": "测试公司A",
    "operation_type": "export",
    "module_name": "保单管理",
    "operation_desc": "导出保单数据",
    "request_url": "/api/v1/policies/export",
    "request_method": "POST",
    "request_params": {
      "export_type": "xlsx"
    },
    "ip_address": "192.168.1.101",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    "operation_time": new Date("2024-03-15T18:20:00Z"),
    "execution_time": 2500,
    "result_status": "success",
    "target_id": "",
    "target_name": ""
  },
  {
    "log_id": "AL" + new Date().getTime() + "008",
    "user_id": "admin",
    "username": "admin",
    "company_id": "",
    "company_name": "平台管理",
    "operation_type": "delete",
    "module_name": "公司管理",
    "operation_desc": "删除公司：测试公司B",
    "request_url": "/api/v1/companies/COMP002",
    "request_method": "DELETE",
    "request_params": {},
    "ip_address": "192.168.1.100",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    "operation_time": new Date("2024-03-15T19:30:00Z"),
    "execution_time": 180,
    "result_status": "success",
    "target_id": "COMP002",
    "target_name": "测试公司B"
  },
  {
    "log_id": "AL" + new Date().getTime() + "009",
    "user_id": "user001",
    "username": "zhangsan",
    "company_id": "COMP001",
    "company_name": "测试公司A",
    "operation_type": "import",
    "module_name": "保单管理",
    "operation_desc": "批量导入保单数据",
    "request_url": "/api/v1/policies/import",
    "request_method": "POST",
    "request_params": {
      "file_name": "policies_import.xlsx",
      "total_count": 50
    },
    "ip_address": "192.168.1.101",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    "operation_time": new Date("2024-03-15T20:15:00Z"),
    "execution_time": 5000,
    "result_status": "success",
    "target_id": "",
    "target_name": ""
  },
  {
    "log_id": "AL" + new Date().getTime() + "010",
    "user_id": "admin",
    "username": "admin",
    "company_id": "",
    "company_name": "平台管理",
    "operation_type": "logout",
    "module_name": "认证授权",
    "operation_desc": "用户退出登录",
    "request_url": "/api/v1/auth/logout",
    "request_method": "POST",
    "request_params": {},
    "ip_address": "192.168.1.100",
    "user_agent": "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36",
    "operation_time": new Date("2024-03-15T22:00:00Z"),
    "execution_time": 80,
    "result_status": "success",
    "target_id": "",
    "target_name": ""
  }
])

// 验证数据插入
print("活动记录表初始化完成！")
print("插入记录数：" + db.activity_logs.count())
print("索引创建完成！")

// 显示前几条记录作为验证
print("\n前3条记录示例：")
db.activity_logs.find().limit(3).forEach(function(doc) {
  print("ID: " + doc.log_id + ", 用户: " + doc.username + ", 操作: " + doc.operation_desc + ", 时间: " + doc.operation_time)
})

print("\n初始化脚本执行完成！") 