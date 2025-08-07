# 保险经纪管理系统 - 数据库设计文档

## 数据库架构

- **数据库类型**: MongoDB (NoSQL)
- **连接地址**: mongodb://localhost:27017
- **数据库名**: insurance_db
- **字符集**: UTF-8
- **时区**: Asia/Shanghai

## 设计原则

1. **无外键约束**: 采用NoSQL设计理念，通过业务逻辑保证数据一致性
2. **业务主键**: 每个表都有独立的业务主键，便于分布式扩展
3. **软删除**: 重要数据采用状态标记而非物理删除
4. **审计跟踪**: 记录创建时间、更新时间和操作日志
5. **多租户隔离**: 通过company_id实现数据隔离

---

## 核心业务表

### 1. 用户表 (users)

存储系统所有用户的基本信息和认证数据。

| 字段名 | 数据类型 | 是否必填 | 索引 | 中文说明 | 备注 |
|--------|----------|----------|------|----------|------|
| _id | ObjectId | 是 | 主键 | MongoDB主键ID | 系统自动生成 |
| user_id | String | 是 | 唯一索引 | 用户唯一标识 | 业务主键，格式：USR+时间戳+随机数 |
| username | String | 是 | 唯一索引 | 登录用户名 | 3-50字符，全局唯一 |
| display_name | String | 是 | 普通索引 | 用户显示名称 | 2-100字符，可包含中文 |
| company_id | String | 是 | 普通索引 | 所属保险经纪公司ID | 关联companies表 |
| role_ids | Array | 是 | - | 用户角色ID数组 | 支持多角色，关联roles表 |
| status | String | 是 | 普通索引 | 用户状态 | active=激活, inactive=禁用, locked=锁定 |
| last_login_time | Date | 否 | - | 最后登录时间 | 可为空，记录最近一次登录 |
| password_hash | String | 是 | - | 密码哈希值 | BCrypt加密，不返回给前端 |
| email | String | 否 | 稀疏索引 | 邮箱地址 | 可选，需符合邮箱格式 |
| phone | String | 否 | - | 手机号码 | 可选，支持中国手机号格式 |
| remark | String | 否 | - | 备注信息 | 可选，管理员备注 |
| login_attempts | Number | 是 | - | 登录失败次数 | 防暴力破解，超过阈值锁定账户 |
| locked_until | Date | 否 | - | 账户锁定截止时间 | 可为空，锁定期间禁止登录 |
| created_at | Date | 是 | 普通索引 | 创建时间 | 记录用户注册时间 |
| updated_at | Date | 是 | - | 更新时间 | 每次修改时自动更新 |

### 2. 保险经纪公司表 (companies)

存储接入平台的保险经纪公司信息和配额管理。

| 字段名 | 数据类型 | 是否必填 | 索引 | 中文说明 | 备注 |
|--------|----------|----------|------|----------|------|
| _id | ObjectId | 是 | 主键 | MongoDB主键ID | 系统自动生成 |
| company_id | String | 是 | 唯一索引 | 公司唯一标识 | 业务主键，格式：CMP+时间戳+随机数 |
| company_name | String | 是 | 普通索引 | 公司名称 | 保险经纪公司完整名称 |
| address | String | 否 | - | 公司地址 | 公司注册或办公地址 |
| contact_phone | String | 是 | - | 联系电话 | 公司客服或联系电话 |
| email | String | 是 | - | 邮箱地址 | 公司官方邮箱 |
| valid_start_date | Date | 是 | 普通索引 | 有效期开始日期 | 服务开始日期 |
| valid_end_date | Date | 是 | 普通索引 | 有效期结束日期 | 服务到期日期 |
| user_quota | Number | 是 | - | 用户数量配额 | 允许创建的最大用户数 |
| current_user_count | Number | 是 | - | 当前用户数量 | 已创建的用户数量 |
| status | String | 是 | 普通索引 | 公司状态 | active=有效, inactive=停用, expired=过期 |
| remark | String | 否 | - | 备注信息 | 平台管理员备注 |
| created_at | Date | 是 | 普通索引 | 创建时间 | 公司接入时间 |
| updated_at | Date | 是 | - | 更新时间 | 最后修改时间 |

### 3. 角色表 (roles)

存储系统角色信息，支持平台级和公司级角色。

| 字段名 | 数据类型 | 是否必填 | 索引 | 中文说明 | 备注 |
|--------|----------|----------|------|----------|------|
| _id | ObjectId | 是 | 主键 | MongoDB主键ID | 系统自动生成 |
| role_id | String | 是 | 唯一索引 | 角色唯一标识 | 业务主键，格式：ROL+时间戳+随机数 |
| role_name | String | 是 | 普通索引 | 角色名称 | 角色的显示名称 |
| role_key | String | 是 | 唯一索引 | 角色标识符 | 英文标识，用于权限判断 |
| company_id | String | 否 | 普通索引 | 所属公司ID | 空表示平台级角色 |
| sort_order | Number | 是 | - | 排序号 | 角色显示顺序 |
| data_scope | String | 是 | - | 数据权限范围 | all=全部, company=本公司, self=个人 |
| menu_ids | Array | 是 | - | 菜单权限ID数组 | 关联menus表，控制页面访问 |
| status | String | 是 | 普通索引 | 角色状态 | enable=启用, disable=禁用 |
| remark | String | 否 | - | 备注信息 | 角色说明 |
| created_at | Date | 是 | 普通索引 | 创建时间 | 角色创建时间 |
| updated_at | Date | 是 | - | 更新时间 | 最后修改时间 |

### 4. 菜单表 (menus)

存储系统菜单和权限控制信息。

| 字段名 | 数据类型 | 是否必填 | 索引 | 中文说明 | 备注 |
|--------|----------|----------|------|----------|------|
| _id | ObjectId | 是 | 主键 | MongoDB主键ID | 系统自动生成 |
| menu_id | String | 是 | 唯一索引 | 菜单唯一标识 | 业务主键 |
| parent_id | String | 否 | 普通索引 | 父菜单ID | 根菜单为空，构建树形结构 |
| menu_name | String | 是 | - | 菜单名称 | 显示在界面上的菜单名 |
| menu_type | String | 是 | 普通索引 | 菜单类型 | directory=目录, menu=菜单, button=按钮 |
| route_path | String | 否 | - | 路由路径 | 前端路由地址 |
| component | String | 否 | - | 组件路径 | 前端组件文件路径 |
| permission_code | String | 否 | 普通索引 | 权限标识符 | 用于后端权限校验 |
| icon | String | 否 | - | 菜单图标 | 图标class或图片地址 |
| sort_order | Number | 是 | - | 排序号 | 菜单显示顺序 |
| visible | Boolean | 是 | - | 是否显示 | 控制菜单是否在界面显示 |
| status | String | 是 | 普通索引 | 菜单状态 | enable=启用, disable=禁用 |
| created_at | Date | 是 | - | 创建时间 | 菜单创建时间 |
| updated_at | Date | 是 | - | 更新时间 | 最后修改时间 |

---

## 业务数据表

### 5. 保单表 (policies)

存储保单业务数据，采用动态字段设计。

| 字段名 | 数据类型 | 是否必填 | 索引 | 中文说明 | 备注 |
|--------|----------|----------|------|----------|------|
| _id | ObjectId | 是 | 主键 | MongoDB主键ID | 系统自动生成 |
| policy_id | String | 是 | 唯一索引 | 保单唯一标识 | 业务主键 |
| company_id | String | 是 | 普通索引 | 所属公司ID | 数据隔离 |
| table_id | String | 是 | 普通索引 | 动态表结构ID | 关联table_structures表 |
| user_id | String | 是 | 普通索引 | 创建用户ID | 记录操作人 |
| data | Object | 是 | - | 动态字段数据 | JSON格式存储业务数据 |
| status | String | 是 | 普通索引 | 保单状态 | draft=草稿, active=生效, expired=过期, cancelled=取消 |
| created_at | Date | 是 | 普通索引 | 创建时间 | 保单录入时间 |
| updated_at | Date | 是 | - | 更新时间 | 最后修改时间 |

### 6. 动态表结构定义表 (table_structures)

定义动态表的结构信息。

| 字段名 | 数据类型 | 是否必填 | 索引 | 中文说明 | 备注 |
|--------|----------|----------|------|----------|------|
| _id | ObjectId | 是 | 主键 | MongoDB主键ID | 系统自动生成 |
| table_id | String | 是 | 唯一索引 | 表结构唯一标识 | 业务主键 |
| table_name | String | 是 | 普通索引 | 表名 | 英文标识 |
| display_name | String | 是 | - | 表显示名称 | 中文名称 |
| table_type | String | 是 | 普通索引 | 表类型 | system=系统表, custom=自定义表 |
| company_id | String | 否 | 普通索引 | 所属公司ID | 空表示平台级表 |
| description | String | 否 | - | 表描述 | 表的用途说明 |
| status | String | 是 | 普通索引 | 表状态 | active=启用, inactive=禁用 |
| created_at | Date | 是 | - | 创建时间 | 表创建时间 |
| updated_at | Date | 是 | - | 更新时间 | 最后修改时间 |

### 7. 字段定义表 (field_definitions)

定义动态表的字段结构。

| 字段名 | 数据类型 | 是否必填 | 索引 | 中文说明 | 备注 |
|--------|----------|----------|------|----------|------|
| _id | ObjectId | 是 | 主键 | MongoDB主键ID | 系统自动生成 |
| field_id | String | 是 | 唯一索引 | 字段唯一标识 | 业务主键 |
| table_id | String | 是 | 普通索引 | 所属表ID | 关联table_structures表 |
| field_name | String | 是 | - | 字段名 | 英文标识 |
| display_name | String | 是 | - | 字段显示名称 | 中文名称 |
| field_type | String | 是 | - | 字段类型 | string=文本, number=数字, date=日期, boolean=布尔, enum=枚举, file=文件 |
| field_length | Number | 否 | - | 字段长度限制 | 文本类型的最大长度 |
| required | Boolean | 是 | - | 是否必填 | 表单验证规则 |
| default_value | String | 否 | - | 默认值 | 字段默认值 |
| enum_options | Array | 否 | - | 枚举选项 | 字段类型为enum时的选项 |
| validation_rules | Object | 否 | - | 验证规则 | JSON格式的自定义验证规则 |
| sort_order | Number | 是 | - | 排序号 | 字段显示顺序 |
| visible | Boolean | 是 | - | 是否显示 | 控制字段是否在表单显示 |
| created_at | Date | 是 | - | 创建时间 | 字段创建时间 |
| updated_at | Date | 是 | - | 更新时间 | 最后修改时间 |

---

## 系统日志表

### 8. 系统操作日志表 (operation_logs)

记录用户在系统中的所有操作行为。

| 字段名 | 数据类型 | 是否必填 | 索引 | 中文说明 | 备注 |
|--------|----------|----------|------|----------|------|
| _id | ObjectId | 是 | 主键 | MongoDB主键ID | 系统自动生成 |
| log_id | String | 是 | 唯一索引 | 日志唯一标识 | 业务主键 |
| user_id | String | 是 | 普通索引 | 操作用户ID | 关联users表 |
| username | String | 是 | 普通索引 | 操作用户名 | 冗余存储，便于查询 |
| company_id | String | 是 | 普通索引 | 所属公司ID | 数据隔离 |
| operation_type | String | 是 | 普通索引 | 操作类型 | create=创建, update=更新, delete=删除, view=查看, export=导出 |
| module_name | String | 是 | 普通索引 | 模块名称 | 如：用户管理、保单管理等 |
| operation_desc | String | 是 | - | 操作描述 | 具体操作说明 |
| request_url | String | 是 | - | 请求URL | API接口地址 |
| request_method | String | 是 | - | 请求方法 | GET, POST, PUT, DELETE等 |
| request_params | Object | 否 | - | 请求参数 | JSON格式的请求参数 |
| ip_address | String | 是 | 普通索引 | 操作IP地址 | 客户端IP |
| user_agent | String | 是 | - | 浏览器标识 | 客户端浏览器信息 |
| operation_time | Date | 是 | 普通索引 | 操作时间 | 操作发生时间 |
| execution_time | Number | 是 | - | 执行耗时 | 毫秒为单位 |
| result_status | String | 是 | 普通索引 | 执行结果 | success=成功, failure=失败 |
| error_message | String | 否 | - | 错误信息 | 操作失败时的错误详情 |

### 9. 数据变更记录表 (data_change_logs)

记录重要数据的变更历史。

| 字段名 | 数据类型 | 是否必填 | 索引 | 中文说明 | 备注 |
|--------|----------|----------|------|----------|------|
| _id | ObjectId | 是 | 主键 | MongoDB主键ID | 系统自动生成 |
| change_id | String | 是 | 唯一索引 | 变更记录唯一标识 | 业务主键 |
| table_name | String | 是 | 普通索引 | 操作表名 | 被变更的表名 |
| record_id | String | 是 | 普通索引 | 记录ID | 被变更记录的唯一标识 |
| user_id | String | 是 | 普通索引 | 操作用户ID | 关联users表 |
| company_id | String | 是 | 普通索引 | 所属公司ID | 数据隔离 |
| change_type | String | 是 | 普通索引 | 变更类型 | insert=新增, update=更新, delete=删除 |
| old_values | Object | 否 | - | 变更前数据 | JSON格式，更新和删除时记录 |
| new_values | Object | 否 | - | 变更后数据 | JSON格式，新增和更新时记录 |
| changed_fields | Array | 否 | - | 变更字段列表 | 具体变更的字段名数组 |
| change_time | Date | 是 | 普通索引 | 变更时间 | 数据变更时间 |
| change_reason | String | 否 | - | 变更原因 | 变更说明或原因 |

---

## 索引设计

### 单字段索引
- `users.user_id` (唯一)
- `users.username` (唯一)
- `users.company_id` (普通)
- `companies.company_id` (唯一)
- `roles.role_id` (唯一)
- `roles.role_key` (唯一)
- `menus.menu_id` (唯一)

### 复合索引
- `users.company_id + users.status` (用户查询优化)
- `policies.company_id + policies.status` (保单查询优化)
- `operation_logs.user_id + operation_logs.operation_time` (日志查询优化)
- `data_change_logs.table_name + data_change_logs.record_id` (变更记录查询优化)

### 时间索引
- `users.created_at` (创建时间查询)
- `operation_logs.operation_time` (日志时间查询)
- `data_change_logs.change_time` (变更时间查询)

---

## 数据完整性

### 1. 业务规则
- 用户名全局唯一
- 公司名称全局唯一
- 角色标识符全局唯一
- 用户数量不能超过公司配额

### 2. 状态控制
- 用户状态：active, inactive, locked
- 公司状态：active, inactive, expired
- 角色状态：enable, disable
- 保单状态：draft, active, expired, cancelled

### 3. 级联操作
- 公司禁用时，该公司所有用户自动禁用
- 角色删除时，用户的角色关联自动移除
- 菜单删除时，角色的菜单权限自动移除

---

## 性能优化

### 1. 查询优化
- 合理设置索引，避免全表扫描
- 使用复合索引优化常用查询
- 分页查询使用skip + limit

### 2. 存储优化
- 使用MongoDB的GridFS存储大文件
- 定期归档历史日志数据
- 使用TTL索引自动清理过期数据

### 3. 并发控制
- 使用乐观锁控制并发更新
- 关键操作使用分布式锁
- 读写分离提高查询性能

---

## 安全设计

### 1. 数据加密
- 密码使用BCrypt加密存储
- 敏感字段可考虑AES加密
- 传输过程使用HTTPS

### 2. 访问控制
- 基于角色的权限控制(RBAC)
- 数据权限按公司隔离
- API接口权限验证

### 3. 审计跟踪
- 记录所有重要操作日志
- 数据变更历史可追溯
- 异常操作告警机制

---

## 备份策略

### 1. 数据备份
- 每日全量备份
- 实时增量备份
- 异地容灾备份

### 2. 恢复测试
- 定期备份恢复测试
- 制定应急恢复预案
- 文档化恢复流程

---

*本文档版本：v1.0*  
*更新日期：2024年12月*  
*维护人员：开发团队* 