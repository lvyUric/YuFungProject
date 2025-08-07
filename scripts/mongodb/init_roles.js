// 初始化角色数据脚本
// 用于在MongoDB中创建预设角色

// 连接到insurance_db数据库
use('insurance_db');

// 删除现有角色数据（仅用于初始化）
db.roles.deleteMany({});

// 预设角色数据
const defaultRoles = [
    {
        role_id: "ROLE001",
        role_name: "超级管理员",
        role_key: "super_admin",
        company_id: "", // 空表示平台级角色
        sort_order: 1,
        data_scope: "all", // 全部数据权限
        menu_ids: [], // 所有菜单权限，后续可以添加具体菜单ID
        status: "enable",
        remark: "平台最高权限管理员，拥有所有功能权限",
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        role_id: "ROLE002",
        role_name: "平台管理员",
        role_key: "platform_admin",
        company_id: "", // 空表示平台级角色
        sort_order: 2,
        data_scope: "all", // 全部数据权限
        menu_ids: [],
        status: "enable",
        remark: "平台管理员，负责公司和用户管理",
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        role_id: "ROLE003",
        role_name: "公司管理员",
        role_key: "company_admin",
        company_id: "", // 平台级角色，但限制数据范围为公司
        sort_order: 3,
        data_scope: "company", // 本公司数据权限
        menu_ids: [],
        status: "enable",
        remark: "公司内最高权限管理员，管理本公司用户和业务",
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        role_id: "ROLE004",
        role_name: "业务员",
        role_key: "salesman",
        company_id: "", // 平台级角色
        sort_order: 4,
        data_scope: "self", // 个人数据权限
        menu_ids: [],
        status: "enable",
        remark: "业务人员，负责保单录入和客户管理",
        created_at: new Date(),
        updated_at: new Date()
    },
    {
        role_id: "ROLE005",
        role_name: "只读用户",
        role_key: "readonly_user",
        company_id: "", // 平台级角色
        sort_order: 5,
        data_scope: "self", // 个人数据权限
        menu_ids: [],
        status: "enable",
        remark: "只读权限用户，仅可查看数据",
        created_at: new Date(),
        updated_at: new Date()
    }
];

// 插入预设角色
const result = db.roles.insertMany(defaultRoles);
print(`成功插入 ${result.insertedIds.length} 个预设角色`);

// 创建索引
try {
    // 角色ID唯一索引
    db.roles.createIndex({ "role_id": 1 }, { unique: true, name: "idx_role_id_unique" });
    
    // 角色标识符唯一索引
    db.roles.createIndex({ "role_key": 1 }, { unique: true, name: "idx_role_key_unique" });
    
    // 公司ID索引
    db.roles.createIndex({ "company_id": 1 }, { name: "idx_company_id" });
    
    // 状态索引
    db.roles.createIndex({ "status": 1 }, { name: "idx_status" });
    
    // 创建时间索引
    db.roles.createIndex({ "created_at": -1 }, { name: "idx_created_at" });
    
    print("角色表索引创建成功");
} catch (error) {
    print("索引创建失败（可能已存在）: " + error.message);
}

// 验证数据
const roleCount = db.roles.countDocuments();
print(`角色表当前记录数: ${roleCount}`);

// 显示所有角色
print("当前角色列表:");
db.roles.find({}, { 
    role_id: 1, 
    role_name: 1, 
    role_key: 1, 
    data_scope: 1, 
    status: 1 
}).forEach(role => {
    print(`- ${role.role_name} (${role.role_key}) - ${role.data_scope} - ${role.status}`);
});

print("角色初始化完成！"); 