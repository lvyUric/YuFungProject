// MongoDB 菜单初始化脚本
print("开始初始化菜单数据...");

// 连接到数据库
db = db.getSiblingDB('insurance_db');

// 清空菜单表（重新初始化）
db.menus.deleteMany({});
print("清空菜单表完成");

// 获取当前时间
const now = new Date();

// 系统默认菜单数据
const defaultMenus = [
    // 一级菜单 - 系统概览
    {
        menu_id: "MENU_001",
        parent_id: "",
        menu_name: "系统概览",
        menu_type: "directory",
        route_path: "/dashboard",
        component: "",
        permission_code: "dashboard",
        icon: "dashboard",
        sort_order: 1,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_001_001",
        parent_id: "MENU_001",
        menu_name: "首页",
        menu_type: "menu",
        route_path: "/welcome",
        component: "./Welcome",
        permission_code: "dashboard:welcome",
        icon: "home",
        sort_order: 1,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },

    // 一级菜单 - 公司管理
    {
        menu_id: "MENU_002",
        parent_id: "",
        menu_name: "公司管理",
        menu_type: "directory",
        route_path: "/company",
        component: "",
        permission_code: "company",
        icon: "bank",
        sort_order: 2,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_002_001",
        parent_id: "MENU_002",
        menu_name: "公司列表",
        menu_type: "menu",
        route_path: "/company/list",
        component: "./company-list",
        permission_code: "company:list",
        icon: "unordered-list",
        sort_order: 1,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_002_001_001",
        parent_id: "MENU_002_001",
        menu_name: "新增公司",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "company:create",
        icon: "plus",
        sort_order: 1,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_002_001_002",
        parent_id: "MENU_002_001",
        menu_name: "编辑公司",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "company:edit",
        icon: "edit",
        sort_order: 2,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_002_001_003",
        parent_id: "MENU_002_001",
        menu_name: "删除公司",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "company:delete",
        icon: "delete",
        sort_order: 3,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },

    // 一级菜单 - 用户管理
    {
        menu_id: "MENU_003",
        parent_id: "",
        menu_name: "用户管理",
        menu_type: "menu",
        route_path: "/user-management",
        component: "./user-management",
        permission_code: "user",
        icon: "user",
        sort_order: 3,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_003_001",
        parent_id: "MENU_003",
        menu_name: "新增用户",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "user:create",
        icon: "user-add",
        sort_order: 1,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_003_002",
        parent_id: "MENU_003",
        menu_name: "编辑用户",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "user:edit",
        icon: "edit",
        sort_order: 2,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_003_003",
        parent_id: "MENU_003",
        menu_name: "删除用户",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "user:delete",
        icon: "user-delete",
        sort_order: 3,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },

    // 一级菜单 - 权限管理
    {
        menu_id: "MENU_004",
        parent_id: "",
        menu_name: "权限管理",
        menu_type: "directory",
        route_path: "/permission",
        component: "",
        permission_code: "permission",
        icon: "safety-certificate",
        sort_order: 4,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_004_001",
        parent_id: "MENU_004",
        menu_name: "角色管理",
        menu_type: "menu",
        route_path: "/role-management",
        component: "./role-management",
        permission_code: "role",
        icon: "team",
        sort_order: 1,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_004_001_001",
        parent_id: "MENU_004_001",
        menu_name: "新增角色",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "role:create",
        icon: "plus",
        sort_order: 1,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_004_001_002",
        parent_id: "MENU_004_001",
        menu_name: "编辑角色",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "role:edit",
        icon: "edit",
        sort_order: 2,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_004_001_003",
        parent_id: "MENU_004_001",
        menu_name: "删除角色",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "role:delete",
        icon: "delete",
        sort_order: 3,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_004_002",
        parent_id: "MENU_004",
        menu_name: "菜单管理",
        menu_type: "menu",
        route_path: "/menu-management",
        component: "./menu-management",
        permission_code: "menu",
        icon: "menu",
        sort_order: 2,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_004_002_001",
        parent_id: "MENU_004_002",
        menu_name: "新增菜单",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "menu:create",
        icon: "plus",
        sort_order: 1,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_004_002_002",
        parent_id: "MENU_004_002",
        menu_name: "编辑菜单",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "menu:edit",
        icon: "edit",
        sort_order: 2,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_004_002_003",
        parent_id: "MENU_004_002",
        menu_name: "删除菜单",
        menu_type: "button",
        route_path: "",
        component: "",
        permission_code: "menu:delete",
        icon: "delete",
        sort_order: 3,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },

    // 一级菜单 - 账户管理
    {
        menu_id: "MENU_005",
        parent_id: "",
        menu_name: "账户管理",
        menu_type: "directory",
        route_path: "/account",
        component: "",
        permission_code: "account",
        icon: "user",
        sort_order: 5,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    },
    {
        menu_id: "MENU_005_001",
        parent_id: "MENU_005",
        menu_name: "修改密码",
        menu_type: "menu",
        route_path: "/account/change-password",
        component: "./user/changePassword",
        permission_code: "account:password",
        icon: "key",
        sort_order: 1,
        visible: true,
        status: "enable",
        created_at: now,
        updated_at: now
    }
];

// 插入菜单数据
try {
    const result = db.menus.insertMany(defaultMenus);
    print(`成功插入 ${result.insertedIds.length} 个菜单记录`);
} catch (error) {
    print("插入菜单数据失败: " + error);
}

// 创建索引
print("开始创建菜单表索引...");

try {
    // 菜单ID唯一索引
    db.menus.createIndex({ "menu_id": 1 }, { unique: true, name: "idx_menu_id_unique" });
    
    // 父菜单ID索引
    db.menus.createIndex({ "parent_id": 1 }, { name: "idx_parent_id" });
    
    // 菜单类型索引
    db.menus.createIndex({ "menu_type": 1 }, { name: "idx_menu_type" });
    
    // 状态索引
    db.menus.createIndex({ "status": 1 }, { name: "idx_status" });
    
    // 排序索引
    db.menus.createIndex({ "sort_order": 1 }, { name: "idx_sort_order" });
    
    // 权限标识符索引（稀疏索引）
    db.menus.createIndex({ "permission_code": 1 }, { sparse: true, name: "idx_permission_code" });
    
    // 复合索引：父菜单ID + 排序号
    db.menus.createIndex({ "parent_id": 1, "sort_order": 1 }, { name: "idx_parent_sort" });
    
    print("菜单表索引创建成功");
} catch (error) {
    print("创建菜单表索引失败: " + error);
}

// 验证数据
const menuCount = db.menus.countDocuments({});
print(`菜单表当前共有 ${menuCount} 条记录`);

// 显示菜单树结构
print("\n=== 菜单树结构 ===");
const rootMenus = db.menus.find({ "parent_id": "" }).sort({ "sort_order": 1 });
rootMenus.forEach(function(menu) {
    print(`├─ ${menu.menu_name} (${menu.menu_type}) - ${menu.route_path}`);
    
    // 查找子菜单
    const childMenus = db.menus.find({ "parent_id": menu.menu_id }).sort({ "sort_order": 1 });
    childMenus.forEach(function(child) {
        print(`│  ├─ ${child.menu_name} (${child.menu_type}) - ${child.route_path}`);
        
        // 查找子菜单的子菜单（按钮权限）
        const grandChildMenus = db.menus.find({ "parent_id": child.menu_id }).sort({ "sort_order": 1 });
        grandChildMenus.forEach(function(grandChild) {
            print(`│  │  ├─ ${grandChild.menu_name} (${grandChild.menu_type}) - ${grandChild.permission_code}`);
        });
    });
});

print("\n菜单数据初始化完成！"); 