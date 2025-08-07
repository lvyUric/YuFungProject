// MongoDB RBAC中间表初始化脚本
print("开始初始化RBAC中间表...");

// 连接到数据库
db = db.getSiblingDB('insurance_db');

// 清空中间表（重新初始化）
db.user_roles.deleteMany({});
db.role_permissions.deleteMany({});
print("清空RBAC中间表完成");

// 获取当前时间
const now = new Date();

// 查询现有数据
const adminUser = db.users.findOne({ username: "admin" });
const superAdminRole = db.roles.findOne({ role_key: "super_admin" });

// 如果存在管理员用户和超级管理员角色，建立关联
if (adminUser && superAdminRole) {
    // 为admin用户分配super_admin角色
    const userRole = {
        user_id: adminUser.user_id,
        role_id: superAdminRole.role_id,
        created_at: now,
        updated_at: now
    };
    
    db.user_roles.insertOne(userRole);
    print(`为用户 ${adminUser.username} 分配角色 ${superAdminRole.role_name}`);
    
    // 获取所有菜单为超级管理员角色分配权限
    const allMenus = db.menus.find({}).toArray();
    
    if (allMenus.length > 0) {
        const rolePermissions = allMenus.map(menu => ({
            role_id: superAdminRole.role_id,
            menu_id: menu.menu_id,
            permission_type: "menu",
            created_at: now,
            updated_at: now
        }));
        
        db.role_permissions.insertMany(rolePermissions);
        print(`为超级管理员角色分配 ${allMenus.length} 个菜单权限`);
    }
}

// 创建索引
print("开始创建RBAC中间表索引...");

try {
    // 用户角色表索引
    db.user_roles.createIndex({ "user_id": 1, "role_id": 1 }, { unique: true, name: "idx_user_role_unique" });
    db.user_roles.createIndex({ "user_id": 1 }, { name: "idx_user_id" });
    db.user_roles.createIndex({ "role_id": 1 }, { name: "idx_role_id" });
    
    // 角色权限表索引
    db.role_permissions.createIndex({ "role_id": 1, "menu_id": 1 }, { unique: true, name: "idx_role_permission_unique" });
    db.role_permissions.createIndex({ "role_id": 1 }, { name: "idx_role_id_perm" });
    db.role_permissions.createIndex({ "menu_id": 1 }, { name: "idx_menu_id_perm" });
    
    print("RBAC中间表索引创建成功");
} catch (error) {
    print("创建RBAC中间表索引失败: " + error);
}

// 验证数据
const userRoleCount = db.user_roles.countDocuments({});
const rolePermissionCount = db.role_permissions.countDocuments({});
print(`用户角色关联表当前共有 ${userRoleCount} 条记录`);
print(`角色权限关联表当前共有 ${rolePermissionCount} 条记录`);

print("\nRBAC中间表初始化完成！"); 