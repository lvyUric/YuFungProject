// MongoDB 完整测试数据初始化脚本
print("========================================");
print("开始初始化完整测试数据...");
print("========================================");

// 连接到数据库
db = db.getSiblingDB('insurance_db');

// 获取当前时间
const now = new Date();

// ========================================
// 1. 清空现有数据
// ========================================
print("\n1. 清空现有数据...");
db.companies.deleteMany({});
db.users.deleteMany({});
db.roles.deleteMany({});
db.menus.deleteMany({});
db.user_roles.deleteMany({});
db.role_permissions.deleteMany({});
print("数据清空完成");

// ========================================
// 2. 初始化公司数据
// ========================================
print("\n2. 初始化公司数据...");

const companies = [
    {
        company_id: "COMP001",
        company_name: "平安保险经纪有限公司",
        company_code: "PINGAN",
        contact_person: "张经理",
        tel_no: "021-12345678",
        mobile: "13800138001",
        email: "pingan@example.com",
        address_cn_province: "上海市",
        address_cn_city: "上海市",
        address_cn_district: "浦东新区",
        address_cn_detail: "陆家嘴金融贸易区",
        broker_code: "PA001",
        username: "pingan_admin",
        password_hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // secret
        valid_start_date: new Date("2024-01-01"),
        valid_end_date: new Date("2025-12-31"),
        user_quota: 50,
        current_user_count: 0,
        status: "active",
        remark: "测试公司1",
        submitted_by: "system",
        created_at: now,
        updated_at: now
    },
    {
        company_id: "COMP002", 
        company_name: "中国人寿保险经纪公司",
        company_code: "CHINALIFE",
        contact_person: "李总监",
        tel_no: "010-87654321",
        mobile: "13900139002",
        email: "chinalife@example.com",
        address_cn_province: "北京市",
        address_cn_city: "北京市", 
        address_cn_district: "朝阳区",
        address_cn_detail: "建国门外大街",
        broker_code: "CL002",
        username: "chinalife_admin",
        password_hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // secret
        valid_start_date: new Date("2024-01-01"),
        valid_end_date: new Date("2025-12-31"),
        user_quota: 30,
        current_user_count: 0,
        status: "active",
        remark: "测试公司2",
        submitted_by: "system",
        created_at: now,
        updated_at: now
    },
    {
        company_id: "COMP003",
        company_name: "太平洋保险经纪",
        company_code: "CPIC",
        contact_person: "王主任",
        tel_no: "020-11223344",
        mobile: "13700137003",
        email: "cpic@example.com",
        address_cn_province: "广东省",
        address_cn_city: "广州市",
        address_cn_district: "天河区",
        address_cn_detail: "珠江新城",
        broker_code: "CP003",
        username: "cpic_admin",
        password_hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // secret
        valid_start_date: new Date("2024-01-01"),
        valid_end_date: new Date("2024-06-30"),
        user_quota: 20,
        current_user_count: 0,
        status: "expired",
        remark: "测试公司3-已过期",
        submitted_by: "system",
        created_at: now,
        updated_at: now
    }
];

db.companies.insertMany(companies);
print(`公司数据初始化完成，共插入 ${companies.length} 条记录`);

// ========================================
// 3. 初始化角色数据
// ========================================
print("\n3. 初始化角色数据...");

const roles = [
    // 平台级角色
    {
        role_id: "ROLE001",
        role_name: "超级管理员",
        role_key: "super_admin", 
        company_id: "", // 平台级角色
        sort_order: 1,
        data_scope: "all",
        menu_ids: [], // 将在后面分配
        status: "enable",
        remark: "系统超级管理员，拥有所有权限",
        created_at: now,
        updated_at: now
    },
    {
        role_id: "ROLE002",
        role_name: "平台管理员",
        role_key: "platform_admin",
        company_id: "", // 平台级角色
        sort_order: 2,
        data_scope: "all",
        menu_ids: [],
        status: "enable", 
        remark: "平台管理员，管理公司和用户",
        created_at: now,
        updated_at: now
    },
    
    // 公司级角色 - 平安保险
    {
        role_id: "ROLE003",
        role_name: "公司管理员",
        role_key: "company_admin",
        company_id: "COMP001",
        sort_order: 1,
        data_scope: "company",
        menu_ids: [],
        status: "enable",
        remark: "公司管理员，管理本公司用户和数据",
        created_at: now,
        updated_at: now
    },
    {
        role_id: "ROLE004", 
        role_name: "业务经理",
        role_key: "business_manager",
        company_id: "COMP001",
        sort_order: 2,
        data_scope: "company",
        menu_ids: [],
        status: "enable",
        remark: "业务经理，管理保单和客户",
        created_at: now,
        updated_at: now
    },
    {
        role_id: "ROLE005",
        role_name: "普通员工",
        role_key: "employee",
        company_id: "COMP001", 
        sort_order: 3,
        data_scope: "self",
        menu_ids: [],
        status: "enable",
        remark: "普通员工，查看自己的数据",
        created_at: now,
        updated_at: now
    },
    
    // 公司级角色 - 中国人寿
    {
        role_id: "ROLE006",
        role_name: "公司管理员",
        role_key: "company_admin",
        company_id: "COMP002",
        sort_order: 1,
        data_scope: "company",
        menu_ids: [],
        status: "enable",
        remark: "中国人寿公司管理员",
        created_at: now,
        updated_at: now
    },
    {
        role_id: "ROLE007",
        role_name: "销售主管",
        role_key: "sales_supervisor", 
        company_id: "COMP002",
        sort_order: 2,
        data_scope: "company",
        menu_ids: [],
        status: "enable",
        remark: "销售主管",
        created_at: now,
        updated_at: now
    }
];

db.roles.insertMany(roles);
print(`角色数据初始化完成，共插入 ${roles.length} 条记录`);

// ========================================
// 4. 初始化菜单数据
// ========================================
print("\n4. 初始化菜单数据...");

const menus = [
    // 一级菜单 - 系统管理
    {
        menu_id: "MENU001",
        parent_id: "",
        menu_name: "系统管理", 
        menu_type: "directory",
        route_path: "/system",
        component: "",
        permission_code: "system",
        icon: "setting",
        sort_order: 1,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    
    // 二级菜单 - 系统管理子菜单
    {
        menu_id: "MENU002",
        parent_id: "MENU001",
        menu_name: "用户管理",
        menu_type: "menu",
        route_path: "/system/user",
        component: "./user-management",
        permission_code: "system:user:view",
        icon: "user",
        sort_order: 1,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU003",
        parent_id: "MENU001",
        menu_name: "角色管理",
        menu_type: "menu",
        route_path: "/system/role", 
        component: "./role-management",
        permission_code: "system:role:view",
        icon: "team",
        sort_order: 2,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU004",
        parent_id: "MENU001",
        menu_name: "菜单管理",
        menu_type: "menu",
        route_path: "/system/menu",
        component: "./menu-management",
        permission_code: "system:menu:view",
        icon: "menu",
        sort_order: 3,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU005",
        parent_id: "MENU001",
        menu_name: "公司管理",
        menu_type: "menu",
        route_path: "/system/company",
        component: "./company-management", 
        permission_code: "system:company:view",
        icon: "bank",
        sort_order: 4,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    
    // 三级菜单 - 用户管理按钮权限
    {
        menu_id: "MENU006",
        parent_id: "MENU002",
        menu_name: "用户查询",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:user:list",
        icon: "",
        sort_order: 1,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU007",
        parent_id: "MENU002",
        menu_name: "用户新增",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:user:add",
        icon: "",
        sort_order: 2,
        visible: false,
        status: "enable", 
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU008",
        parent_id: "MENU002",
        menu_name: "用户修改",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:user:edit",
        icon: "",
        sort_order: 3,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU009",
        parent_id: "MENU002",
        menu_name: "用户删除",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:user:remove",
        icon: "",
        sort_order: 4,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU010",
        parent_id: "MENU002",
        menu_name: "重置密码",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:user:resetPwd",
        icon: "",
        sort_order: 5,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    
    // 三级菜单 - 角色管理按钮权限
    {
        menu_id: "MENU011",
        parent_id: "MENU003",
        menu_name: "角色查询",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:role:list",
        icon: "",
        sort_order: 1,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU012",
        parent_id: "MENU003",
        menu_name: "角色新增",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:role:add",
        icon: "",
        sort_order: 2,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU013",
        parent_id: "MENU003",
        menu_name: "角色修改",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:role:edit",
        icon: "",
        sort_order: 3,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU014",
        parent_id: "MENU003",
        menu_name: "角色删除",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:role:remove",
        icon: "",
        sort_order: 4,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    
    // 三级菜单 - 菜单管理按钮权限
    {
        menu_id: "MENU015",
        parent_id: "MENU004",
        menu_name: "菜单查询",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:menu:list",
        icon: "",
        sort_order: 1,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU016",
        parent_id: "MENU004",
        menu_name: "菜单新增",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:menu:add",
        icon: "",
        sort_order: 2,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU017",
        parent_id: "MENU004",
        menu_name: "菜单修改",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:menu:edit",
        icon: "",
        sort_order: 3,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU018",
        parent_id: "MENU004",
        menu_name: "菜单删除",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:menu:remove",
        icon: "",
        sort_order: 4,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    
    // 三级菜单 - 公司管理按钮权限
    {
        menu_id: "MENU019",
        parent_id: "MENU005",
        menu_name: "公司查询",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:company:list",
        icon: "",
        sort_order: 1,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU020",
        parent_id: "MENU005",
        menu_name: "公司新增",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:company:add",
        icon: "",
        sort_order: 2,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU021",
        parent_id: "MENU005",
        menu_name: "公司修改",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:company:edit",
        icon: "",
        sort_order: 3,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU022",
        parent_id: "MENU005",
        menu_name: "公司删除",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "system:company:remove",
        icon: "",
        sort_order: 4,
        visible: false,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    
    // 一级菜单 - 业务管理
    {
        menu_id: "MENU023",
        parent_id: "",
        menu_name: "业务管理",
        menu_type: "directory",
        route_path: "/business",
        component: "",
        permission_code: "business",
        icon: "solution",
        sort_order: 2,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU024",
        parent_id: "MENU023",
        menu_name: "保单管理",
        menu_type: "menu",
        route_path: "/business/policy",
        component: "./policy-management",
        permission_code: "business:policy:view",
        icon: "file-text",
        sort_order: 1,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU025",
        parent_id: "MENU023", 
        menu_name: "客户管理",
        menu_type: "menu",
        route_path: "/business/customer",
        component: "./customer-management",
        permission_code: "business:customer:view",
        icon: "contacts",
        sort_order: 2,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    
    // 一级菜单 - 首页
    {
        menu_id: "MENU026",
        parent_id: "",
        menu_name: "首页",
        menu_type: "menu",
        route_path: "/dashboard",
        component: "./dashboard",
        permission_code: "dashboard",
        icon: "dashboard",
        sort_order: 0,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    }
];

db.menus.insertMany(menus);
print(`菜单数据初始化完成，共插入 ${menus.length} 条记录`);

// ========================================
// 5. 初始化用户数据
// ========================================
print("\n5. 初始化用户数据...");

const users = [
    // 平台管理员
    {
        user_id: "USER001",
        username: "admin",
        display_name: "系统管理员",
        company_id: "",
        role_ids: [], // 将在用户角色关联表中设置
        status: "active",
        last_login_time: null,
        password_hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // secret
        email: "admin@system.com",
        phone: "13800000000",
        remark: "系统超级管理员",
        login_attempts: 0,
        locked_until: null,
        created_at: now,
        updated_at: now
    },
    {
        user_id: "USER002",
        username: "platform_admin",
        display_name: "平台管理员",
        company_id: "",
        role_ids: [],
        status: "active",
        last_login_time: null,
        password_hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // secret
        email: "platform@system.com",
        phone: "13800000001",
        remark: "平台管理员",
        login_attempts: 0,
        locked_until: null,
        created_at: now,
        updated_at: now
    },
    
    // 平安保险用户
    {
        user_id: "USER003",
        username: "pingan_admin",
        display_name: "平安管理员",
        company_id: "COMP001",
        role_ids: [],
        status: "active",
        last_login_time: null,
        password_hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // secret
        email: "admin@pingan.com",
        phone: "13800138001",
        remark: "平安保险公司管理员",
        login_attempts: 0,
        locked_until: null,
        created_at: now,
        updated_at: now
    },
    {
        user_id: "USER004",
        username: "zhang_manager",
        display_name: "张经理",
        company_id: "COMP001", 
        role_ids: [],
        status: "active",
        last_login_time: null,
        password_hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // secret
        email: "zhang@pingan.com",
        phone: "13800138002",
        remark: "平安保险业务经理",
        login_attempts: 0,
        locked_until: null,
        created_at: now,
        updated_at: now
    },
    {
        user_id: "USER005",
        username: "li_employee",
        display_name: "李员工",
        company_id: "COMP001",
        role_ids: [],
        status: "active",
        last_login_time: null,
        password_hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // secret
        email: "li@pingan.com",
        phone: "13800138003",
        remark: "平安保险普通员工",
        login_attempts: 0,
        locked_until: null,
        created_at: now,
        updated_at: now
    },
    
    // 中国人寿用户
    {
        user_id: "USER006",
        username: "chinalife_admin",
        display_name: "人寿管理员",
        company_id: "COMP002",
        role_ids: [],
        status: "active",
        last_login_time: null,
        password_hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // secret
        email: "admin@chinalife.com",
        phone: "13900139001",
        remark: "中国人寿公司管理员",
        login_attempts: 0,
        locked_until: null,
        created_at: now,
        updated_at: now
    },
    {
        user_id: "USER007",
        username: "wang_supervisor",
        display_name: "王主管",
        company_id: "COMP002",
        role_ids: [],
        status: "active",
        last_login_time: null,
        password_hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // secret
        email: "wang@chinalife.com",
        phone: "13900139002",
        remark: "中国人寿销售主管",
        login_attempts: 0,
        locked_until: null,
        created_at: now,
        updated_at: now
    },
    
    // 禁用用户示例
    {
        user_id: "USER008",
        username: "disabled_user",
        display_name: "禁用用户",
        company_id: "COMP001",
        role_ids: [],
        status: "inactive",
        last_login_time: null,
        password_hash: "$2a$10$N9qo8uLOickgx2ZMRZoMyeIjZAgcfl7p92ldGxad68LJZdL17lhWy", // secret
        email: "disabled@pingan.com",
        phone: "13800138004",
        remark: "已禁用的测试用户",
        login_attempts: 0,
        locked_until: null,
        created_at: now,
        updated_at: now
    }
];

db.users.insertMany(users);
print(`用户数据初始化完成，共插入 ${users.length} 条记录`);

// ========================================
// 6. 初始化用户角色关联
// ========================================
print("\n6. 初始化用户角色关联...");

const userRoles = [
    // 超级管理员
    {
        user_id: "USER001",
        role_id: "ROLE001",
        created_at: now,
        updated_at: now
    },
    
    // 平台管理员
    {
        user_id: "USER002", 
        role_id: "ROLE002",
        created_at: now,
        updated_at: now
    },
    
    // 平安保险用户角色
    {
        user_id: "USER003",
        role_id: "ROLE003",
        created_at: now,
        updated_at: now
    },
    {
        user_id: "USER004",
        role_id: "ROLE004",
        created_at: now,
        updated_at: now
    },
    {
        user_id: "USER005",
        role_id: "ROLE005",
        created_at: now,
        updated_at: now
    },
    
    // 中国人寿用户角色
    {
        user_id: "USER006",
        role_id: "ROLE006",
        created_at: now,
        updated_at: now
    },
    {
        user_id: "USER007",
        role_id: "ROLE007",
        created_at: now,
        updated_at: now
    },
    
    // 禁用用户也分配一个角色
    {
        user_id: "USER008",
        role_id: "ROLE005",
        created_at: now,
        updated_at: now
    }
];

db.user_roles.insertMany(userRoles);
print(`用户角色关联初始化完成，共插入 ${userRoles.length} 条记录`);

// ========================================
// 7. 初始化角色权限关联
// ========================================
print("\n7. 初始化角色权限关联...");

// 获取所有菜单ID
const allMenuIds = db.menus.find({}).toArray().map(menu => menu.menu_id);

// 超级管理员 - 拥有所有权限
const superAdminPermissions = allMenuIds.map(menuId => ({
    role_id: "ROLE001",
    menu_id: menuId,
    permission_type: "menu",
    created_at: now,
    updated_at: now
}));

// 平台管理员 - 拥有系统管理权限
const systemMenuIds = ["MENU001", "MENU002", "MENU003", "MENU004", "MENU005", "MENU026", 
                      "MENU006", "MENU007", "MENU008", "MENU009", "MENU010", // 用户管理按钮
                      "MENU011", "MENU012", "MENU013", "MENU014", // 角色管理按钮
                      "MENU015", "MENU016", "MENU017", "MENU018", // 菜单管理按钮
                      "MENU019", "MENU020", "MENU021", "MENU022"]; // 公司管理按钮

const platformAdminPermissions = systemMenuIds.map(menuId => ({
    role_id: "ROLE002", 
    menu_id: menuId,
    permission_type: "menu",
    created_at: now,
    updated_at: now
}));

// 公司管理员 - 平安保险（除了公司管理）
const companyAdminMenuIds = ["MENU001", "MENU002", "MENU003", "MENU023", "MENU024", "MENU025", "MENU026",
                           "MENU006", "MENU007", "MENU008", "MENU009", "MENU010", // 用户管理按钮
                           "MENU011", "MENU012", "MENU013", "MENU014"]; // 角色管理按钮

const companyAdminPermissions = companyAdminMenuIds.map(menuId => ({
    role_id: "ROLE003",
    menu_id: menuId,
    permission_type: "menu",
    created_at: now,
    updated_at: now
}));

// 业务经理 - 平安保险
const businessManagerMenuIds = ["MENU023", "MENU024", "MENU025", "MENU026"];

const businessManagerPermissions = businessManagerMenuIds.map(menuId => ({
    role_id: "ROLE004",
    menu_id: menuId,
    permission_type: "menu", 
    created_at: now,
    updated_at: now
}));

// 普通员工 - 平安保险（只能查看）
const employeeMenuIds = ["MENU024", "MENU025", "MENU026"];

const employeePermissions = employeeMenuIds.map(menuId => ({
    role_id: "ROLE005",
    menu_id: menuId,
    permission_type: "menu",
    created_at: now,
    updated_at: now
}));

// 公司管理员 - 中国人寿
const chinalife_companyAdminPermissions = companyAdminMenuIds.map(menuId => ({
    role_id: "ROLE006",
    menu_id: menuId,
    permission_type: "menu",
    created_at: now,
    updated_at: now
}));

// 销售主管 - 中国人寿
const salesSupervisorPermissions = businessManagerMenuIds.map(menuId => ({
    role_id: "ROLE007",
    menu_id: menuId,
    permission_type: "menu",
    created_at: now,
    updated_at: now
}));

// 合并所有权限
const allRolePermissions = [
    ...superAdminPermissions,
    ...platformAdminPermissions,
    ...companyAdminPermissions,
    ...businessManagerPermissions,
    ...employeePermissions,
    ...chinalife_companyAdminPermissions,
    ...salesSupervisorPermissions
];

db.role_permissions.insertMany(allRolePermissions);
print(`角色权限关联初始化完成，共插入 ${allRolePermissions.length} 条记录`);

// ========================================
// 8. 创建索引
// ========================================
print("\n8. 创建索引...");

try {
    // 公司表索引
    db.companies.createIndex({ "company_id": 1 }, { unique: true, name: "idx_company_id" });
    db.companies.createIndex({ "company_name": 1 }, { name: "idx_company_name" });
    db.companies.createIndex({ "status": 1 }, { name: "idx_company_status" });
    
    // 用户表索引
    db.users.createIndex({ "user_id": 1 }, { unique: true, name: "idx_user_id" });
    db.users.createIndex({ "username": 1 }, { unique: true, name: "idx_username" });
    db.users.createIndex({ "company_id": 1 }, { name: "idx_user_company" });
    db.users.createIndex({ "status": 1 }, { name: "idx_user_status" });
    
    // 角色表索引
    db.roles.createIndex({ "role_id": 1 }, { unique: true, name: "idx_role_id" });
    db.roles.createIndex({ "role_key": 1, "company_id": 1 }, { unique: true, name: "idx_role_key_company" });
    db.roles.createIndex({ "company_id": 1 }, { name: "idx_role_company" });
    
    // 菜单表索引
    db.menus.createIndex({ "menu_id": 1 }, { unique: true, name: "idx_menu_id" });
    db.menus.createIndex({ "parent_id": 1 }, { name: "idx_menu_parent" });
    db.menus.createIndex({ "menu_type": 1 }, { name: "idx_menu_type" });
    db.menus.createIndex({ "status": 1 }, { name: "idx_menu_status" });
    
    // 用户角色关联表索引
    db.user_roles.createIndex({ "user_id": 1, "role_id": 1 }, { unique: true, name: "idx_user_role_unique" });
    db.user_roles.createIndex({ "user_id": 1 }, { name: "idx_user_id" });
    db.user_roles.createIndex({ "role_id": 1 }, { name: "idx_role_id" });
    
    // 角色权限关联表索引
    db.role_permissions.createIndex({ "role_id": 1, "menu_id": 1 }, { unique: true, name: "idx_role_permission_unique" });
    db.role_permissions.createIndex({ "role_id": 1 }, { name: "idx_role_id_perm" });
    db.role_permissions.createIndex({ "menu_id": 1 }, { name: "idx_menu_id_perm" });
    
    print("所有索引创建成功");
} catch (error) {
    print("创建索引失败: " + error);
}

// ========================================
// 9. 数据统计验证
// ========================================
print("\n9. 数据统计验证...");

const companyCount = db.companies.countDocuments({});
const userCount = db.users.countDocuments({});
const roleCount = db.roles.countDocuments({});
const menuCount = db.menus.countDocuments({});
const userRoleCount = db.user_roles.countDocuments({});
const rolePermissionCount = db.role_permissions.countDocuments({});

print(`公司数量: ${companyCount}`);
print(`用户数量: ${userCount}`);
print(`角色数量: ${roleCount}`);
print(`菜单数量: ${menuCount}`);
print(`用户角色关联数量: ${userRoleCount}`);
print(`角色权限关联数量: ${rolePermissionCount}`);

// ========================================
// 10. 输出测试账号信息
// ========================================
print("\n========================================");
print("测试账号信息:");
print("========================================");
print("超级管理员: admin / secret");
print("平台管理员: platform_admin / secret");
print("平安公司管理员: pingan_admin / secret");
print("平安业务经理: zhang_manager / secret");
print("平安普通员工: li_employee / secret");
print("人寿公司管理员: chinalife_admin / secret");
print("人寿销售主管: wang_supervisor / secret");
print("========================================");

print("\n完整测试数据初始化完成！");
print("========================================"); 