// MongoDB 初始化脚本
// 保险经纪管理系统数据库初始化

// 切换到业务数据库
db = db.getSiblingDB('insurance_db');

print('开始初始化保险经纪管理系统数据库...');

// ===========================
// 创建集合和索引
// ===========================

// 1. 用户表索引
print('创建用户表索引...');
db.users.createIndex({ "user_id": 1 }, { unique: true, name: "idx_user_id" });
db.users.createIndex({ "username": 1 }, { unique: true, name: "idx_username" });
db.users.createIndex({ "company_id": 1 }, { name: "idx_company_id" });
db.users.createIndex({ "email": 1 }, { sparse: true, name: "idx_email" });
db.users.createIndex({ "status": 1 }, { name: "idx_status" });
db.users.createIndex({ "created_at": 1 }, { name: "idx_created_at" });
db.users.createIndex({ "company_id": 1, "status": 1 }, { name: "idx_company_status" });

// 2. 公司表索引
print('创建公司表索引...');
db.companies.createIndex({ "company_id": 1 }, { unique: true, name: "idx_company_id" });
db.companies.createIndex({ "company_name": 1 }, { unique: true, name: "idx_company_name" });
db.companies.createIndex({ "status": 1 }, { name: "idx_status" });
db.companies.createIndex({ "valid_start_date": 1 }, { name: "idx_valid_start" });
db.companies.createIndex({ "valid_end_date": 1 }, { name: "idx_valid_end" });
db.companies.createIndex({ "created_at": 1 }, { name: "idx_created_at" });

// 3. 角色表索引
print('创建角色表索引...');
db.roles.createIndex({ "role_id": 1 }, { unique: true, name: "idx_role_id" });
db.roles.createIndex({ "role_key": 1 }, { unique: true, name: "idx_role_key" });
db.roles.createIndex({ "company_id": 1 }, { name: "idx_company_id" });
db.roles.createIndex({ "status": 1 }, { name: "idx_status" });
db.roles.createIndex({ "created_at": 1 }, { name: "idx_created_at" });

// 4. 菜单表索引
print('创建菜单表索引...');
db.menus.createIndex({ "menu_id": 1 }, { unique: true, name: "idx_menu_id" });
db.menus.createIndex({ "parent_id": 1 }, { name: "idx_parent_id" });
db.menus.createIndex({ "menu_type": 1 }, { name: "idx_menu_type" });
db.menus.createIndex({ "permission_code": 1 }, { name: "idx_permission_code" });
db.menus.createIndex({ "status": 1 }, { name: "idx_status" });

// 5. 保单表索引
print('创建保单表索引...');
db.policies.createIndex({ "policy_id": 1 }, { unique: true, name: "idx_policy_id" });
db.policies.createIndex({ "company_id": 1 }, { name: "idx_company_id" });
db.policies.createIndex({ "table_id": 1 }, { name: "idx_table_id" });
db.policies.createIndex({ "user_id": 1 }, { name: "idx_user_id" });
db.policies.createIndex({ "status": 1 }, { name: "idx_status" });
db.policies.createIndex({ "created_at": 1 }, { name: "idx_created_at" });
db.policies.createIndex({ "company_id": 1, "status": 1 }, { name: "idx_company_status" });

// 6. 表结构定义表索引
print('创建表结构表索引...');
db.table_structures.createIndex({ "table_id": 1 }, { unique: true, name: "idx_table_id" });
db.table_structures.createIndex({ "table_name": 1 }, { name: "idx_table_name" });
db.table_structures.createIndex({ "company_id": 1 }, { name: "idx_company_id" });
db.table_structures.createIndex({ "table_type": 1 }, { name: "idx_table_type" });
db.table_structures.createIndex({ "status": 1 }, { name: "idx_status" });

// 7. 字段定义表索引
print('创建字段定义表索引...');
db.field_definitions.createIndex({ "field_id": 1 }, { unique: true, name: "idx_field_id" });
db.field_definitions.createIndex({ "table_id": 1 }, { name: "idx_table_id" });
db.field_definitions.createIndex({ "field_name": 1 }, { name: "idx_field_name" });

// 8. 操作日志表索引
print('创建操作日志表索引...');
db.operation_logs.createIndex({ "log_id": 1 }, { unique: true, name: "idx_log_id" });
db.operation_logs.createIndex({ "user_id": 1 }, { name: "idx_user_id" });
db.operation_logs.createIndex({ "username": 1 }, { name: "idx_username" });
db.operation_logs.createIndex({ "company_id": 1 }, { name: "idx_company_id" });
db.operation_logs.createIndex({ "operation_type": 1 }, { name: "idx_operation_type" });
db.operation_logs.createIndex({ "module_name": 1 }, { name: "idx_module_name" });
db.operation_logs.createIndex({ "operation_time": 1 }, { name: "idx_operation_time" });
db.operation_logs.createIndex({ "result_status": 1 }, { name: "idx_result_status" });
db.operation_logs.createIndex({ "user_id": 1, "operation_time": 1 }, { name: "idx_user_time" });

// 9. 数据变更记录表索引
print('创建数据变更记录表索引...');
db.data_change_logs.createIndex({ "change_id": 1 }, { unique: true, name: "idx_change_id" });
db.data_change_logs.createIndex({ "table_name": 1 }, { name: "idx_table_name" });
db.data_change_logs.createIndex({ "record_id": 1 }, { name: "idx_record_id" });
db.data_change_logs.createIndex({ "user_id": 1 }, { name: "idx_user_id" });
db.data_change_logs.createIndex({ "company_id": 1 }, { name: "idx_company_id" });
db.data_change_logs.createIndex({ "change_type": 1 }, { name: "idx_change_type" });
db.data_change_logs.createIndex({ "change_time": 1 }, { name: "idx_change_time" });
db.data_change_logs.createIndex({ "table_name": 1, "record_id": 1 }, { name: "idx_table_record" });

// ===========================
// 插入默认数据
// ===========================

// 获取当前时间
var now = new Date();

// 1. 创建默认平台公司
print('创建默认平台公司...');
var defaultCompany = {
    company_id: "CMP_PLATFORM_001",
    company_name: "平台管理公司",
    address: "系统内置",
    contact_phone: "400-000-0000",
    email: "admin@platform.com",
    valid_start_date: new Date("2024-01-01"),
    valid_end_date: new Date("2099-12-31"),
    user_quota: 9999,
    current_user_count: 1,
    status: "active",
    remark: "系统默认平台管理公司",
    created_at: now,
    updated_at: now
};

// 检查是否已存在默认公司
if (db.companies.countDocuments({ company_id: "CMP_PLATFORM_001" }) === 0) {
    db.companies.insertOne(defaultCompany);
    print('默认平台公司创建成功');
} else {
    print('默认平台公司已存在，跳过创建');
}

// 2. 创建默认角色
print('创建默认角色...');
var defaultRoles = [
    {
        role_id: "ROL_SUPER_ADMIN",
        role_name: "超级管理员",
        role_key: "super_admin",
        company_id: "",
        sort_order: 1,
        data_scope: "all",
        menu_ids: ["*"],
        status: "enable",
        remark: "平台超级管理员，拥有所有权限",
        created_at: now,
        updated_at: now
    },
    {
        role_id: "ROL_COMPANY_ADMIN",
        role_name: "公司管理员",
        role_key: "company_admin",
        company_id: "",
        sort_order: 2,
        data_scope: "company",
        menu_ids: [],
        status: "enable",
        remark: "公司管理员，管理本公司用户和数据",
        created_at: now,
        updated_at: now
    },
    {
        role_id: "ROL_NORMAL_USER",
        role_name: "普通用户",
        role_key: "normal_user",
        company_id: "",
        sort_order: 3,
        data_scope: "self",
        menu_ids: [],
        status: "enable",
        remark: "普通用户，基础操作权限",
        created_at: now,
        updated_at: now
    },
    {
        role_id: "ROL_READONLY_USER",
        role_name: "只读用户",
        role_key: "readonly_user",
        company_id: "",
        sort_order: 4,
        data_scope: "self",
        menu_ids: [],
        status: "enable",
        remark: "只读用户，仅查看权限",
        created_at: now,
        updated_at: now
    }
];

defaultRoles.forEach(function(role) {
    if (db.roles.countDocuments({ role_id: role.role_id }) === 0) {
        db.roles.insertOne(role);
        print('角色创建成功: ' + role.role_name);
    } else {
        print('角色已存在，跳过创建: ' + role.role_name);
    }
});

// 3. 创建默认菜单
print('创建默认菜单...');
var defaultMenus = [
    {
        menu_id: "MENU_DASHBOARD",
        parent_id: "",
        menu_name: "仪表盘",
        menu_type: "menu",
        route_path: "/dashboard",
        component: "Dashboard",
        permission_code: "dashboard:view",
        icon: "dashboard",
        sort_order: 1,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_USER_MGMT",
        parent_id: "",
        menu_name: "用户管理",
        menu_type: "directory",
        route_path: "/user",
        component: "",
        permission_code: "",
        icon: "user",
        sort_order: 2,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_USER_LIST",
        parent_id: "MENU_USER_MGMT",
        menu_name: "用户列表",
        menu_type: "menu",
        route_path: "/user/list",
        component: "UserList",
        permission_code: "user:list",
        icon: "",
        sort_order: 1,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_COMPANY_MGMT",
        parent_id: "",
        menu_name: "公司管理",
        menu_type: "menu",
        route_path: "/company",
        component: "CompanyManagement",
        permission_code: "company:view",
        icon: "company",
        sort_order: 3,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_POLICY_MGMT",
        parent_id: "",
        menu_name: "保单管理",
        menu_type: "menu",
        route_path: "/policy",
        component: "PolicyManagement",
        permission_code: "policy:view",
        icon: "policy",
        sort_order: 4,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_SYSTEM_MGMT",
        parent_id: "",
        menu_name: "系统管理",
        menu_type: "directory",
        route_path: "/system",
        component: "",
        permission_code: "",
        icon: "system",
        sort_order: 5,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_ROLE_MGMT",
        parent_id: "MENU_SYSTEM_MGMT",
        menu_name: "角色管理",
        menu_type: "menu",
        route_path: "/system/role",
        component: "RoleManagement",
        permission_code: "role:view",
        icon: "",
        sort_order: 1,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_MENU_MGMT",
        parent_id: "MENU_SYSTEM_MGMT",
        menu_name: "菜单管理",
        menu_type: "menu",
        route_path: "/system/menu",
        component: "MenuManagement",
        permission_code: "menu:view",
        icon: "",
        sort_order: 2,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    }
];

defaultMenus.forEach(function(menu) {
    if (db.menus.countDocuments({ menu_id: menu.menu_id }) === 0) {
        db.menus.insertOne(menu);
        print('菜单创建成功: ' + menu.menu_name);
    } else {
        print('菜单已存在，跳过创建: ' + menu.menu_name);
    }
});

// 4. 创建默认超级管理员用户
print('创建默认超级管理员用户...');
// 密码：admin123（BCrypt加密后的值）
var defaultAdmin = {
    user_id: "USR_ADMIN_001",
    username: "admin",
    display_name: "系统管理员",
    company_id: "CMP_PLATFORM_001",
    role_ids: ["ROL_SUPER_ADMIN"],
    status: "active",
    last_login_time: null,
    password_hash: "$2a$10$N.zmdr9k7uOCQb376NoUnuTJ8iKVn5fyeCp53rjz0scdMCQQCKRha", // admin123
    email: "admin@platform.com",
    phone: "",
    remark: "系统默认管理员账户",
    login_attempts: 0,
    locked_until: null,
    created_at: now,
    updated_at: now
};

if (db.users.countDocuments({ username: "admin" }) === 0) {
    db.users.insertOne(defaultAdmin);
    print('默认管理员用户创建成功 - 用户名: admin, 密码: admin123');
} else {
    print('默认管理员用户已存在，跳过创建');
}

// 5. 创建默认表结构（保单表）
print('创建默认保单表结构...');
var defaultTableStructure = {
    table_id: "TABLE_POLICY_BASIC",
    table_name: "policy_basic",
    display_name: "基础保单表",
    table_type: "system",
    company_id: "",
    description: "系统默认的基础保单表结构",
    status: "active",
    created_at: now,
    updated_at: now
};

if (db.table_structures.countDocuments({ table_id: "TABLE_POLICY_BASIC" }) === 0) {
    db.table_structures.insertOne(defaultTableStructure);
    print('默认保单表结构创建成功');
} else {
    print('默认保单表结构已存在，跳过创建');
}

// 6. 创建默认字段定义
print('创建默认字段定义...');
var defaultFields = [
    {
        field_id: "FIELD_POLICY_NO",
        table_id: "TABLE_POLICY_BASIC",
        field_name: "policy_no",
        display_name: "保单号",
        field_type: "string",
        field_length: 50,
        required: true,
        default_value: "",
        enum_options: [],
        validation_rules: { "pattern": "^[A-Z0-9]{10,20}$" },
        sort_order: 1,
        visible: true,
        created_at: now,
        updated_at: now
    },
    {
        field_id: "FIELD_CUSTOMER_NAME",
        table_id: "TABLE_POLICY_BASIC",
        field_name: "customer_name",
        display_name: "客户姓名",
        field_type: "string",
        field_length: 100,
        required: true,
        default_value: "",
        enum_options: [],
        validation_rules: {},
        sort_order: 2,
        visible: true,
        created_at: now,
        updated_at: now
    },
    {
        field_id: "FIELD_INSURANCE_AMOUNT",
        table_id: "TABLE_POLICY_BASIC",
        field_name: "insurance_amount",
        display_name: "保险金额",
        field_type: "number",
        field_length: 0,
        required: true,
        default_value: "0",
        enum_options: [],
        validation_rules: { "min": 0 },
        sort_order: 3,
        visible: true,
        created_at: now,
        updated_at: now
    },
    {
        field_id: "FIELD_EFFECTIVE_DATE",
        table_id: "TABLE_POLICY_BASIC",
        field_name: "effective_date",
        display_name: "生效日期",
        field_type: "date",
        field_length: 0,
        required: true,
        default_value: "",
        enum_options: [],
        validation_rules: {},
        sort_order: 4,
        visible: true,
        created_at: now,
        updated_at: now
    },
    {
        field_id: "FIELD_EXPIRY_DATE",
        table_id: "TABLE_POLICY_BASIC",
        field_name: "expiry_date",
        display_name: "到期日期",
        field_type: "date",
        field_length: 0,
        required: true,
        default_value: "",
        enum_options: [],
        validation_rules: {},
        sort_order: 5,
        visible: true,
        created_at: now,
        updated_at: now
    }
];

defaultFields.forEach(function(field) {
    if (db.field_definitions.countDocuments({ field_id: field.field_id }) === 0) {
        db.field_definitions.insertOne(field);
        print('字段定义创建成功: ' + field.display_name);
    } else {
        print('字段定义已存在，跳过创建: ' + field.display_name);
    }
});

print('保险经纪管理系统数据库初始化完成！');
print('');
print('=================================');
print('默认管理员账户信息：');
print('用户名: admin');
print('密码: admin123');
print('=================================');
print('');
print('请修改默认密码以确保系统安全！'); 