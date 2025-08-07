# 保险经纪公司管理系统

基于 Gin + MongoDB + React + TypeScript + TailwindCSS 构建的保险经纪公司管理平台。

## 🎯 功能特性

✅ **已实现功能**：
- 🔐 用户登录/注册
- 🔑 修改密码
- 🚪 用户注销
- 📱 响应式设计（支持PC端和移动端）
- 🎨 现代化UI界面
- 🔒 JWT认证机制
- 📊 用户仪表盘
- 🛡️ 密码强度验证
- ⚡ 令牌自动刷新
- 🗄️ 完整数据库设计（9个核心表）
- 📋 数据库自动初始化
- 👤 默认管理员账户

## 🚀 快速启动

### 📋 前置条件
- Go 1.21+
- Node.js 16+
- MongoDB 4.4+

### 🔥 一键启动（推荐）

#### Windows用户
```bash
# 启动后端（会自动初始化数据库）
start-backend.bat

# 启动前端（新开一个命令窗口）
start-frontend.bat
```

#### Linux/Mac用户
```bash
# 一键启动前后端（会自动初始化数据库）
./start.sh
```

### 📖 手动启动

#### 1. 启动MongoDB
```bash
# 方法1: 本地安装的MongoDB
# Windows: net start MongoDB
# Linux/Mac: sudo systemctl start mongod

# 方法2: Docker启动MongoDB（推荐）
docker run --restart=always --name mongo-dev \
  -p 27017:27017 \
  -e TZ=Asia/Shanghai \
  -v /data/mongo/data:/data/db \
  -v /data/mongo/log:/data/log \
  --privileged=true \
  -e MONGO_INITDB_ROOT_USERNAME=admin \
  -e MONGO_INITDB_ROOT_PASSWORD=yf2025 \
  -d mongo
```

#### 2. 初始化数据库（首次启动）
```bash
# Windows
init-database.bat

# Linux/Mac
mongo insurance_db scripts/init-mongo.js
```

#### 3. 启动后端服务
```bash
# 安装Go依赖
go mod tidy

# 启动后端
go run cmd/main.go
```

后端服务将在 `http://localhost:8080` 启动

#### 4. 启动前端服务
```bash
# 进入前端目录
cd stelory-admin

# 安装依赖（首次运行）
npm install

# 启动开发服务器
npm start
```

前端服务将在 `http://localhost:3000` 启动

### 🔑 默认账户信息

系统初始化后会自动创建默认管理员账户：

| 字段 | 值 |
|------|-----|
| **用户名** | `admin` |
| **密码** | `admin123` |
| **角色** | 超级管理员 |
| **权限** | 所有功能权限 |

⚠️ **重要提醒**：首次登录后请立即修改默认密码！

## 🏗️ 技术栈

### 后端
- **框架**: Gin (Go语言Web框架)
- **数据库**: MongoDB
- **认证**: JWT + BCrypt密码加密
- **配置管理**: Viper
- **API文档**: RESTful API
- **中间件**: CORS, 认证, 日志

### 前端
- **框架**: React 18 + TypeScript
- **状态管理**: Zustand
- **UI样式**: TailwindCSS + 自定义组件
- **图标**: Heroicons (SVG)
- **构建工具**: Create React App
- **网络请求**: Fetch API

## 🗄️ 数据库设计

### 核心数据表（9个表）

| 表名 | 中文名称 | 描述 | 主要字段 |
|------|----------|------|----------|
| **users** | 用户表 | 存储系统用户信息 | 用户ID、用户名、密码、角色、状态 |
| **companies** | 公司表 | 保险经纪公司信息 | 公司ID、公司名称、有效期、用户配额 |
| **roles** | 角色表 | 系统角色权限 | 角色ID、角色名称、权限范围、菜单权限 |
| **menus** | 菜单表 | 系统菜单权限 | 菜单ID、菜单名称、路由、权限标识 |
| **policies** | 保单表 | 保单业务数据 | 保单ID、公司ID、动态数据、状态 |
| **table_structures** | 表结构表 | 动态表定义 | 表ID、表名、表类型、状态 |
| **field_definitions** | 字段定义表 | 动态字段定义 | 字段ID、字段名称、字段类型、验证规则 |
| **operation_logs** | 操作日志表 | 用户操作记录 | 日志ID、用户ID、操作类型、操作时间 |
| **data_change_logs** | 变更记录表 | 数据变更历史 | 变更ID、表名、变更类型、变更内容 |

### 设计特点
- ✅ **无外键约束**: NoSQL设计，通过业务逻辑保证一致性
- ✅ **业务主键**: 独立的业务标识符，便于分布式扩展
- ✅ **多租户隔离**: 通过company_id实现数据隔离
- ✅ **动态字段**: 支持自定义表结构和字段
- ✅ **完整审计**: 操作日志和数据变更记录
- ✅ **索引优化**: 合理的索引设计提升查询性能

详细的数据库设计文档请查看：[数据库设计文档](docs/database-design.md)

## 📁 项目结构

```
YufungProject/
├── cmd/                     # Go应用程序入口
│   └── main.go             # 主程序
├── configs/                 # 配置文件
│   ├── application.yml     # 应用配置
│   └── config.go          # 配置结构体
├── internal/               # 内部应用代码
│   ├── controller/         # 控制器层
│   │   └── auth_controller.go
│   ├── middleware/         # 中间件
│   │   ├── auth.go        # JWT认证
│   │   └── cors.go        # 跨域处理
│   ├── model/             # 数据模型
│   │   ├── user.go        # 用户模型（含9个完整数据表）
│   │   └── response.go    # 响应模型
│   ├── repository/        # 数据访问层
│   │   └── user_repository.go
│   ├── routes/           # 路由配置
│   │   ├── auth_routes.go
│   │   └── routes.go
│   └── service/          # 业务逻辑层
│       └── auth_service.go
├── pkg/                   # 公共库包
│   ├── database/         # 数据库连接
│   │   └── mongodb.go
│   └── utils/           # 工具函数
│       ├── id.go        # ID生成
│       ├── jwt.go       # JWT工具
│       └── password.go  # 密码工具
├── stelory-admin/        # React前端项目
│   ├── src/
│   │   ├── components/   # React组件
│   │   │   ├── AuthPage.tsx      # 认证页面
│   │   │   ├── LoginForm.tsx     # 登录表单
│   │   │   ├── RegisterForm.tsx  # 注册表单
│   │   │   ├── ChangePasswordForm.tsx # 修改密码
│   │   │   └── Dashboard.tsx     # 仪表盘
│   │   ├── store/       # Zustand状态管理
│   │   │   └── authStore.ts
│   │   ├── App.tsx      # 主应用组件
│   │   └── index.css    # 样式文件
│   ├── tailwind.config.js # TailwindCSS配置
│   └── package.json     # 前端依赖
├── scripts/             # 数据库脚本
│   └── init-mongo.js   # 数据库初始化脚本
├── docs/               # 项目文档
│   └── database-design.md # 数据库设计文档
├── start-backend.bat    # Windows后端启动脚本
├── start-frontend.bat   # Windows前端启动脚本
├── start.sh            # Linux/Mac启动脚本
├── init-database.bat   # Windows数据库初始化脚本
├── go.mod             # Go模块依赖
└── README.md          # 项目说明
```

## 🔌 API接口

### 认证相关
| 方法 | 路径 | 描述 | 权限要求 |
|------|------|------|----------|
| POST | `/api/auth/login` | 用户登录 | 无 |
| POST | `/api/auth/register` | 用户注册 | 无 |
| POST | `/api/auth/logout` | 用户登出 | 需要登录 |
| POST | `/api/auth/change-password` | 修改密码 | 需要登录 |
| POST | `/api/auth/refresh` | 刷新令牌 | 无 |
| GET | `/api/auth/user-info` | 获取用户信息 | 需要登录 |

### 健康检查
| 方法 | 路径 | 描述 |
|------|------|------|
| GET | `/health` | 服务健康状态 |

## 📚 使用说明

### 用户注册
1. 访问 `http://localhost:3000`
2. 点击"立即注册"
3. 填写用户信息：
   - **用户名**（必填，至少3个字符）
   - **显示名称**（必填，至少2个字符）
   - **密码**（必填，至少8个字符，需包含小写字母、数字和大写字母或特殊字符）
   - **邮箱**（可选，需符合邮箱格式）
   - **手机号**（可选，需符合中国手机号格式）
4. 点击"注册账户"

### 用户登录
1. 在登录页面输入用户名和密码
2. 点击"登录"按钮
3. 登录成功后将跳转到仪表盘页面

### 修改密码
1. 登录后，点击右上角用户头像下拉菜单
2. 选择"修改密码"
3. 在弹出的对话框中：
   - 输入当前密码
   - 输入新密码（符合密码强度要求）
   - 确认新密码
4. 点击"确认修改"

### 用户注销
1. 点击右上角用户头像下拉菜单
2. 选择"退出登录"
3. 系统将清除本地认证信息并返回登录页面

## 📱 响应式设计

### 桌面端（≥1024px）
- 完整的布局和功能
- 大卡片式界面
- 完整的用户下拉菜单

### 平板端（768px-1023px）
- 自适应布局
- 2列网格显示
- 简化的导航

### 手机端（<768px）
- 单列布局
- 堆叠式卡片
- 优化的触控体验
- 隐藏次要信息

## ⚙️ 配置说明

### 后端配置 (configs/application.yml)
```yaml
server:
  port: 8080          # 服务端口
  mode: debug         # 运行模式: debug/release

database:
  mongodb:
    # 默认本地无认证配置
    uri: mongodb://localhost:27017
    # Docker认证配置（取消注释使用）
    # uri: mongodb://admin:yf2025@localhost:27017/admin
    database: insurance_db

jwt:
  secret: your-secret-key-change-in-production  # JWT密钥
  expires_in: 24h                              # 访问令牌有效期
  refresh_expires_in: 168h                     # 刷新令牌有效期

security:
  password_min_length: 8    # 最小密码长度
  max_login_attempts: 5     # 最大登录尝试次数
  lockout_duration: 30m     # 锁定时长
```

### 前端配置
- API地址在 `src/store/authStore.ts` 中配置
- 默认为 `http://localhost:8080/api`
- TailwindCSS配置在 `tailwind.config.js`

## 🛡️ 安全特性

- 🔐 **JWT令牌认证**：无状态认证机制
- 🔒 **BCrypt密码加密**：安全的密码存储
- 🚫 **防暴力破解**：登录失败自动锁定
- ✅ **密码强度验证**：前后端双重验证
- 🔄 **令牌自动刷新**：无缝用户体验
- 🛡️ **CORS防护**：跨域请求保护
- 📝 **输入验证**：全面的数据验证

## 🌐 浏览器支持

| 浏览器 | 版本要求 |
|--------|----------|
| Chrome | ≥ 90 |
| Firefox | ≥ 90 |
| Safari | ≥ 14 |
| Edge | ≥ 90 |

## 🔧 开发指南

### 添加新组件
1. 在 `stelory-admin/src/components/` 创建新组件
2. 使用 TypeScript 和函数式组件
3. 遵循响应式设计原则
4. 使用 TailwindCSS 进行样式设计

### 状态管理
- 使用 Zustand 进行状态管理
- 在 `src/store/` 目录下创建新的 store
- 保持 store 简洁，单一职责

### API集成
- 在 store 中定义 API 调用
- 使用 fetch API 进行网络请求
- 统一错误处理和加载状态

### 代码规范
- TypeScript 严格模式
- ESLint 代码检查
- Prettier 代码格式化
- 组件和函数命名使用 PascalCase 和 camelCase

## 📈 性能优化

- ⚡ **代码分割**：按需加载组件
- 🎯 **懒加载**：优化首屏加载速度
- 📦 **构建优化**：生产环境代码压缩
- 🚀 **缓存策略**：合理的HTTP缓存设置

## 🔮 下一步开发计划

- [ ] 用户管理模块
- [ ] 公司管理模块  
- [ ] 保单管理模块
- [ ] RBAC权限管理
- [ ] 数据统计仪表盘
- [ ] 文件上传功能
- [ ] 数据导出功能
- [ ] 系统日志模块
- [ ] 通知中心
- [ ] 多语言支持

## 🐛 常见问题

### MongoDB连接失败
```bash
# 确保MongoDB服务正在运行
# Windows
net start MongoDB

# 或使用Docker启动（推荐）
docker run --restart=always --name mongo-dev \
-p 27017:27017 \
-e TZ=Asia/Shanghai \
--privileged=true \
-e MONGO_INITDB_ROOT_USERNAME=admin \
-e MONGO_INITDB_ROOT_PASSWORD=yf2025 \
-d mongo

# Linux
sudo systemctl start mongod

# macOS
brew services start mongodb-community
```

### 前端编译错误
```bash
# 清除node_modules重新安装
rm -rf node_modules package-lock.json
npm install
```

### 端口被占用
```bash
# 查看端口占用
netstat -ano | findstr :8080
netstat -ano | findstr :3000

# 终止进程或更改配置文件中的端口
```

### 数据库初始化失败
```bash
# 手动初始化数据库
# Windows
init-database.bat

# Linux/Mac
mongo insurance_db scripts/init-mongo.js
```

## 🤝 贡献指南

1. Fork 项目
2. 创建功能分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 创建 Pull Request

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情

## 👥 联系方式

如果您有任何问题或建议，请通过以下方式联系我们：

- 提交 [Issue](../../issues)
- 发送邮件至：your-email@example.com

---

**保险经纪管理系统** - 让保险业务管理更简单、更高效！ 🏥✨ 