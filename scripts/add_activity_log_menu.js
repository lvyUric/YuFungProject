// 添加活动记录菜单的MongoDB脚本
// 在MongoDB中执行以下命令

// 1. 添加活动记录菜单
db.menus.insertOne({
  menu_id: "activity_log",
  parent_id: "", // 根菜单
  menu_name: "活动记录",
  menu_type: "menu",
  route_path: "/activity-log",
  component: "ActivityLog",
  permission_code: "activity:log:list",
  icon: "HistoryOutlined",
  sort_order: 100,
  visible: true,
  status: "enable",
  created_at: new Date(),
  updated_at: new Date()
});

// 2. 为超级管理员角色添加活动记录权限
// 假设超级管理员角色ID为 "super_admin"
db.roles.updateOne(
  { role_id: "super_admin" },
  { 
    $addToSet: { 
      menu_ids: "activity_log" 
    },
    $set: { updated_at: new Date() }
  }
);

// 3. 为其他需要访问活动记录的角色添加权限
// 例如：公司管理员角色
db.roles.updateOne(
  { role_id: "company_admin" },
  { 
    $addToSet: { 
      menu_ids: "activity_log" 
    },
    $set: { updated_at: new Date() }
  }
);

print("活动记录菜单添加完成！"); 