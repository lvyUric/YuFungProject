# 保险经纪公司管理系统 PRD 文档

## 1. 项目概述

### 1.1 项目背景
构建一个面向保险经纪公司的综合管理平台，支持多租户模式，提供灵活的数据管理和权限控制功能。

### 1.2 项目目标
- 实现保险经纪公司的统一管理
- 提供灵活的用户权限控制系统
- 支持动态表结构定义，满足业务字段自定义需求
- 提供完整的操作审计和数据变更追踪
- 构建响应式、易用的管理界面

### 1.3 技术架构
- **后端**：Gin + GORM + MongoDB
- **前端**：React + Zustand + TypeScript + TailwindCSS + Radix UI
- **部署**：Docker Compose
- **安全**：JWT认证 + RBAC权限管理 + API签名加密

### 1.4 项目目录结构
```
YufungProject/
├── cmd/                          # 应用程序入口
├── configs/                      # 配置文件目录
│   ├── application-dev.yml       # 开发环境配置
│   ├── application-prod.yml      # 生产环境配置
│   ├── application.yml           # 基础配置
│   ├── config.go                # 配置结构体定义
├── docs/                        # 项目文档
├── internal/                    # 内部应用代码
│   ├── controller/              # 控制器层
│   ├── middleware/              # 中间件
│   ├── model/                   # 数据模型
│   ├── repository/              # 数据访问层
│   ├── routes/                  # 路由配置
│   └── service/                 # 业务逻辑层
├── pkg/                         # 公共库包
├── Yufung-admin-front/              # 前端项目目录
├── .gitignore                   # Git忽略文件
├── LICENSE                      # 开源协议
├── README.md                    # 项目说明
├── go.mod                       # Go模块依赖
└── go.sum                       # Go模块校验
```

## 2. 功能模块详述

### 2.1 保险经纪公司管理模块
**权限要求**：平台管理员（超级管理员）专有

#### 2.1.1 功能描述
管理所有接入平台的保险经纪公司，控制各公司的用户配额和服务有效期。

#### 2.1.2 数据字段
| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| company_id | String | 是 | 公司唯一标识 |
| company_name | String | 是 | 公司名称 |
| address_cn_province | String | 否 | 中文地址-省/自治区/直辖市 |
| address_cn_city | String | 否 | 中文地址-市 |
| address_cn_district | String | 否 | 中文地址-县/区 |
| address_cn_detail | String | 否 | 中文地址-详细地址 |
| address_en_province | String | 否 | 英文地址-省/自治区/直辖市 |
| address_en_city | String | 否 | 英文地址-市 |
| address_en_district | String | 否 | 英文地址-县/区 |
| address_en_detail | String | 否 | 英文地址-详细地址 |
| address | String | 否 | 公司地址（兼容字段） |
| contact_phone | String | 是 | 联系电话 |
| email | String | 是 | 邮箱地址 |
| valid_start_date | Date | 是 | 有效期开始日期 |
| valid_end_date | Date | 是 | 有效期结束日期 |
| user_quota | Number | 是 | 允许创建的用户数量 |
| current_user_count | Number | 是 | 当前用户数量 |
| status | Enum | 是 | 状态：active/inactive/expired |
| remark | String | 否 | 备注信息 |
| created_at | Date | 是 | 创建时间 |
| updated_at | Date | 是 | 更新时间 |

#### 2.1.3 功能清单
- **基础CRUD**：新增、删除、修改、查询公司信息
- **批量操作**：批量启用/停用公司
- **数据导出**：支持Excel格式导出
- **快捷停用**：一键停用公司及其所属用户
- **自动管理**：有效期到期自动禁用对应用户
- **统计报表**：公司用户使用情况统计

### 2.2 用户管理模块
**权限要求**：平台管理员专有（查看所有用户），公司管理员（查看本公司用户）

#### 2.2.1 功能描述
管理系统内所有用户，支持多层级权限分配。

#### 2.2.2 数据字段
| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| user_id | String | 是 | 用户唯一标识 |
| username | String | 是 | 登录账户 |
| display_name | String | 是 | 显示名称 |
| company_id | String | 是 | 所属公司ID |
| role_ids | Array | 是 | 角色ID数组 |
| status | Enum | 是 | 状态：active/inactive/locked |
| last_login_time | Date | 否 | 最后登录时间 |
| password_hash | String | 是 | 密码哈希值 |
| email | String | 否 | 邮箱 |
| phone | String | 否 | 手机号 |
| remark | String | 否 | 备注 |
| created_at | Date | 是 | 创建时间 |
| updated_at | Date | 是 | 更新时间 |

#### 2.2.3 功能清单
- **基础CRUD**：用户信息管理
- **角色分配**：为用户分配多个角色
- **密码管理**：重置密码、强制修改密码
- **批量操作**：批量启用/停用用户
- **数据导出**：用户信息导出
- **登录监控**：用户登录状态追踪

### 2.3 RBAC权限管理模块（仿若依）
**权限要求**：平台管理员专有

#### 2.3.1 权限模型设计
```
用户(User) ←→ 角色(Role) ←→ 权限(Permission)
           N:N            N:N
```

#### 2.3.2 菜单管理
| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| menu_id | String | 是 | 菜单ID |
| parent_id | String | 否 | 父菜单ID |
| menu_name | String | 是 | 菜单名称 |
| menu_type | Enum | 是 | 类型：directory/menu/button |
| route_path | String | 否 | 路由路径 |
| component | String | 否 | 组件路径 |
| permission_code | String | 否 | 权限标识 |
| icon | String | 否 | 菜单图标 |
| sort_order | Number | 是 | 排序 |
| visible | Boolean | 是 | 是否显示 |
| status | Enum | 是 | 状态：enable/disable |

#### 2.3.3 角色管理
| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| role_id | String | 是 | 角色ID |
| role_name | String | 是 | 角色名称 |
| role_key | String | 是 | 角色标识 |
| company_id | String | 否 | 所属公司（空表示平台角色） |
| sort_order | Number | 是 | 排序 |
| data_scope | Enum | 是 | 数据权限：all/company/self |
| menu_ids | Array | 是 | 菜单权限数组 |
| status | Enum | 是 | 状态：enable/disable |
| remark | String | 否 | 备注 |

#### 2.3.4 预设角色
- **超级管理员**：平台最高权限
- **公司管理员**：公司内最高权限
- **普通用户**：基础操作权限
- **只读用户**：仅查看权限

### 2.4 保单管理模块
**权限要求**：根据数据权限控制访问范围

#### 2.4.1 功能描述
- 公司用户：管理本公司录入的保单
- 平台管理员：查看所有公司保单（只读）

#### 2.4.2 动态字段设计
基于动态表结构管理，用户可自定义保单字段。

#### 2.4.3 核心功能
- **基础CRUD**：保单信息管理
- **文件上传**：支持保单相关文档上传
- **数据导出**：Excel格式导出
- **高级搜索**：多条件组合搜索
- **批量操作**：批量修改保单状态

### 2.5 动态表结构管理模块
**权限要求**：平台管理员或具备表结构管理权限的用户

#### 2.5.1 功能描述
允许用户自定义表结构和字段，无需修改代码即可扩展业务字段。

#### 2.5.2 表结构定义
| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| table_id | String | 是 | 表唯一标识 |
| table_name | String | 是 | 表名 |
| display_name | String | 是 | 显示名称 |
| table_type | Enum | 是 | 表类型：system/custom |
| company_id | String | 否 | 所属公司（空表示平台表） |
| description | String | 否 | 表描述 |
| status | Enum | 是 | 状态：active/inactive |

#### 2.5.3 字段定义
| 字段名 | 类型 | 必填 | 说明 |
|--------|------|------|------|
| field_id | String | 是 | 字段ID |
| table_id | String | 是 | 所属表ID |
| field_name | String | 是 | 字段名 |
| display_name | String | 是 | 显示名称 |
| field_type | Enum | 是 | 字段类型：string/number/date/boolean/enum/file |
| field_length | Number | 否 | 字段长度 |
| required | Boolean | 是 | 是否必填 |
| default_value | String | 否 | 默认值 |
| enum_options | Array | 否 | 枚举选项 |
| validation_rules | Object | 否 | 验证规则 |
| sort_order | Number | 是 | 排序 |
| visible | Boolean | 是 | 是否显示 |

#### 2.5.4 功能特性
- **可视化设计器**：拖拽式字段设计
- **字段类型支持**：文本、数字、日期、选择、文件等
- **验证规则配置**：长度、格式、范围等验证
- **实时预览**：表单实时预览效果
- **版本管理**：表结构变更版本控制

### 2.6 系统基础功能模块

#### 2.6.1 用户中心
- **个人信息管理**：修改基本信息
- **密码管理**：修改密码、找回密码
- **登录日志**：个人登录记录查看
- **用户注销**：账户注销功能

#### 2.6.2 系统设置
- **系统参数配置**：基础参数设置
- **字典数据管理**：下拉选项数据维护
- **通知公告**：系统消息发布

### 2.7 增强功能模块

#### 2.7.1 系统活动记录
**功能描述**：记录用户在系统中的所有操作行为。

**数据结构**：
| 字段名 | 类型 | 说明 |
|--------|------|------|
| log_id | String | 日志ID |
| user_id | String | 操作用户ID |
| username | String | 用户名 |
| company_id | String | 所属公司ID |
| operation_type | Enum | 操作类型：create/update/delete/view/export |
| module_name | String | 模块名称 |
| operation_desc | String | 操作描述 |
| request_url | String | 请求URL |
| request_method | String | 请求方法 |
| request_params | Object | 请求参数 |
| ip_address | String | IP地址 |
| user_agent | String | 浏览器信息 |
| operation_time | Date | 操作时间 |
| execution_time | Number | 执行耗时(ms) |
| result_status | Enum | 执行结果：success/failure |

#### 2.7.2 数据变更记录
**功能描述**：记录数据表每条记录的变更历史。

**数据结构**：
| 字段名 | 类型 | 说明 |
|--------|------|------|
| change_id | String | 变更记录ID |
| table_name | String | 表名 |
| record_id | String | 记录ID |
| user_id | String | 操作用户ID |
| company_id | String | 所属公司ID |
| change_type | Enum | 变更类型：insert/update/delete |
| old_values | Object | 变更前数据 |
| new_values | Object | 变更后数据 |
| changed_fields | Array | 变更字段列表 |
| change_time | Date | 变更时间 |
| change_reason | String | 变更原因 |

#### 2.7.3 首页仪表盘
**功能组件**：
- **统计卡片**：用户数、公司数、保单数等关键指标
- **操作记录**：最近操作活动（可配置显示条数）
- **数据图表**：用户增长趋势、保单统计等
- **快捷入口**：常用功能快速访问
- **系统通知**：重要公告和提醒

**权限控制**：
- 平台管理员：查看所有公司数据
- 公司用户：仅查看本公司数据

## 3. 技术实现规范

### 3.1 后端架构设计

#### 3.1.1 分层架构
```
Controller Layer (控制器层)
    ↓
Service Layer (业务逻辑层)
    ↓
Repository Layer (数据访问层)
    ↓
Model Layer (数据模型层)
```

#### 3.1.2 目录职责说明
- **cmd/**：应用程序启动入口，包含main.go
- **configs/**：配置文件管理，支持多环境配置
- **internal/controller/**：HTTP请求处理，参数验证，响应格式化
- **internal/service/**：核心业务逻辑，事务处理
- **internal/repository/**：数据访问抽象，CRUD操作封装
- **internal/model/**：数据模型定义，包括数据库模型和DTO
- **internal/middleware/**：中间件，如认证、权限、日志、CORS等
- **internal/routes/**：路由配置，API版本管理
- **pkg/**：可复用的公共包，工具函数

#### 3.1.3 配置管理策略
- **多环境配置**：development、test、production
- **配置热加载**：支持配置文件动态加载
- **敏感信息**：通过环境变量覆盖配置文件
- **RBAC模型**：使用Casbin进行权限管理

### 3.2 数据验证规范
- **前端验证**：实时输入验证，提供良好用户体验
- **后端验证**：完整的数据格式、权限、业务规则验证
- **字段过滤**：修改操作时过滤不允许修改的敏感字段

### 3.2 安全规范
- **密码安全**：BCrypt加密存储，不允许在普通修改接口中处理
- **API安全**：JWT认证 + 公私钥签名加密
- **权限控制**：基于RBAC的细粒度权限控制
- **操作审计**：所有敏感操作记录审计日志

### 3.3 响应式设计
- **移动端适配**：支持手机、平板设备访问
- **界面自适应**：基于TailwindCSS实现响应式布局
- **组件化设计**：基于Radix UI构建一致的用户界面

### 3.4 数据库设计
- **MongoDB**：主数据存储
- **索引优化**：关键查询字段建立索引
- **数据分区**：按公司维度进行数据隔离

## 4. 部署架构

### 4.1 Docker Compose 配置
```yaml
version: '3.8'
services:
  # 后端服务
  insurance-backend:
    build:
      context: .
      dockerfile: Dockerfile
      target: production
    container_name: insurance-backend
    ports:
      - "8088:8088"
    environment:
      - GO_ENV=production
      - DB_CONNECTION=mongodb://mongodb:27017
      - DB_NAME=insurance_db
      - JWT_SECRET=${JWT_SECRET}
      - REDIS_URL=redis://redis:6379
    depends_on:
      - mongodb
      - redis
    volumes:
      - ./configs:/app/configs
    networks:
      - insurance-network

  # 前端服务
  insurance-frontend:
    build:
      context: ./stelory-admin
      dockerfile: Dockerfile
    container_name: insurance-frontend
    ports:
      - "3000:3000"
    environment:
      - REACT_APP_API_URL=http://localhost:8080/api
    depends_on:
      - insurance-backend
    networks:
      - insurance-network

  # MongoDB数据库
  mongodb:
    image: mongo:6.0
    container_name: insurance-mongodb
    ports:
      - "27017:27017"
    environment:
      - MONGO_INITDB_ROOT_USERNAME=${MONGO_USERNAME}
      - MONGO_INITDB_ROOT_PASSWORD=${MONGO_PASSWORD}
      - MONGO_INITDB_DATABASE=insurance_db
    volumes:
      - mongodb_data:/data/db
      - ./scripts/init-mongo.js:/docker-entrypoint-initdb.d/init-mongo.js:ro
    networks:
      - insurance-network

  # Redis缓存
  redis:
    image: redis:7-alpine
    container_name: insurance-redis
    ports:
      - "6379:6379"
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes
    networks:
      - insurance-network

  # Nginx反向代理
  nginx:
    image: nginx:alpine
    container_name: insurance-nginx
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/ssl:/etc/nginx/ssl:ro
    depends_on:
      - insurance-backend
      - insurance-frontend
    networks:
      - insurance-network

volumes:
  mongodb_data:
    driver: local
  redis_data:
    driver: local

networks:
  insurance-network:
    driver: bridge
```

### 4.2 Dockerfile 配置

#### 4.2.1 后端 Dockerfile
```dockerfile
# 多阶段构建
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

# 生产环境镜像
FROM alpine:latest AS production
RUN apk --no-cache add ca-certificates tzdata
WORKDIR /root/

COPY --from=builder /app/main .
COPY --from=builder /app/configs ./configs

EXPOSE 8080
CMD ["./main"]
```

### 4.3 环境配置示例

#### 4.3.1 application.yml (基础配置)
```yaml
# 服务配置
server:
  port: 8088
  mode: debug # release, debug, test

# 数据库配置
database:
  mongodb:
    uri: mongodb://localhost:27017
    database: insurance_db
    timeout: 10s
    max_pool_size: 100

# Redis配置
redis:
  addr: localhost:6379
  password: ""
  db: 0
  pool_size: 10

# JWT配置
jwt:
  secret: your-secret-key
  expires_in: 24h
  refresh_expires_in: 168h

# RBAC配置
rbac:
  model_path: ./configs/rbac_model.conf
  
# 日志配置
log:
  level: info
  format: json
  output: stdout
  file_path: ./logs/app.log
  max_size: 100
  max_age: 30
  max_backups: 10

# 文件上传配置
upload:
  max_size: 10MB
  allowed_types: [jpg, jpeg, png, gif, pdf, doc, docx, xls, xlsx]
  path: ./uploads
```

#### 4.3.2 application-prod.yml (生产环境)
```yaml
server:
  mode: release

database:
  mongodb:
    uri: ${MONGODB_URI}
    database: ${MONGODB_DATABASE}

redis:
  addr: ${REDIS_ADDR}
  password: ${REDIS_PASSWORD}

jwt:
  secret: ${JWT_SECRET}

log:
  level: warn
  output: file
```

## 5. 项目里程碑

### 5.1 第一阶段（基础框架）
- [ ] 项目初始化和环境搭建
- [ ] 用户认证和基础权限系统
- [ ] 基础UI组件库搭建

### 5.2 第二阶段（核心模块）
- [ ] 保险经纪公司管理模块
- [ ] 用户管理模块
- [ ] RBAC权限管理完整实现

### 5.3 第三阶段（业务模块）
- [ ] 保单管理模块
- [ ] 动态表结构管理
- [ ] 系统活动记录

### 5.4 第四阶段（增强功能）
- [ ] 数据变更记录
- [ ] 首页仪表盘
- [ ] 系统优化和性能调优

## 6. 质量保证

### 6.1 测试策略
- **单元测试**：核心业务逻辑单元测试
- **集成测试**：API接口集成测试
- **端到端测试**：关键业务流程E2E测试

### 6.2 代码规范
- **TypeScript严格模式**：确保类型安全
- **ESLint规则**：代码风格统一
- **Git提交规范**：规范化提交信息

### 6.3 监控告警
- **性能监控**：API响应时间监控
- **错误监控**：异常日志收集
- **业务监控**：关键指标监控

## 7. 风险评估

### 7.1 技术风险
- **数据迁移风险**：动态表结构变更的数据兼容性
- **性能风险**：大量数据查询和复杂权限计算
- **安全风险**：多租户数据隔离和权限控制

### 7.2 业务风险
- **需求变更风险**：保险业务复杂性带来的需求变化
- **用户体验风险**：复杂权限系统的易用性平衡

### 7.3 风险缓解措施
- 充分的测试覆盖和分阶段发布
- 完善的备份和回滚机制
- 详细的用户操作文档和培训

---

**文档版本**：v1.0  
**创建日期**：2025年8月4日  
**最后更新**：2025年8月4日