# Yufung Admin 前端

基于 Ant Design Pro 构建的企业级管理后台前端项目。

## 🚀 功能特性

- ✅ **用户登录** - 支持用户名密码登录
- ✅ **用户注册** - 支持新用户注册
- ✅ **修改密码** - 支持用户修改密码
- ✅ **用户信息管理** - 显示用户基本信息
- ✅ **JWT Token 认证** - 自动处理 token 刷新
- ✅ **用户头像** - 右上角头像下拉菜单
- ✅ **国际化支持** - 支持简体中文、繁体中文、英文
- ✅ **响应式设计** - 支持移动端和桌面端

## 📦 技术栈

- **React 19** - 前端框架
- **Ant Design Pro** - 企业级 UI 组件库
- **UmiJS 4** - 应用框架
- **TypeScript** - 类型安全
- **Ant Design 5** - UI 组件库

## 🛠 开发环境要求

- Node.js >= 20.0.0
- npm 或 yarn

## 📖 快速开始

### 1. 启动后端服务器

确保后端服务器已经启动，默认运行在 `http://localhost:8088`

### 2. 启动前端项目

**方式一：使用启动脚本（推荐）**
```bash
# Windows
start-yufung-admin.bat
```

**方式二：手动启动**
```bash
# 进入项目目录
cd Yufung-admin-front

# 安装依赖
npm install

# 启动开发服务器
npm start
```

### 3. 访问应用

浏览器自动打开 `http://localhost:8000`

## 🌐 国际化支持

系统支持三种语言：

| 语言 | 代码 | 状态 |
|------|------|------|
| 简体中文 | zh-CN | ✅ 完整支持 |
| 繁体中文 | zh-TW | ✅ 完整支持 |
| 英文 | en-US | ✅ 完整支持 |

**切换语言：**
- 点击右上角的语言切换图标
- 支持浏览器语言自动检测

## 👤 用户头像功能

**头像下拉菜单包含：**
- 个人中心
- 个人设置
- 修改密码
- 退出登录

**头像显示逻辑：**
- 优先显示 `display_name`
- 如果没有则显示 `username`
- 支持自定义头像图片

## 🔧 配置说明

### 后端服务器地址配置

如果后端服务器运行在不同的地址或端口，请修改：

**文件：** `config/proxy.ts`
```typescript
dev: {
  '/api/': {
    target: 'http://localhost:8088', // 修改为实际的后端地址
    changeOrigin: true,
    pathRewrite: { '^': '' },
  },
},
```

### 国际化配置

**文件：** `config/config.ts`
```typescript
locale: {
  default: 'zh-CN', // 默认语言
  antd: true,
  baseNavigator: true, // 自动检测浏览器语言
},
```

## 📋 API 接口

项目对接的后端接口包括：

| 接口 | 方法 | 路径 | 描述 |
|------|------|------|------|
| 登录 | POST | `/api/auth/login` | 用户登录 |
| 注册 | POST | `/api/auth/register` | 用户注册 |
| 登出 | POST | `/api/auth/logout` | 用户登出 |
| 修改密码 | POST | `/api/auth/change-password` | 修改密码 |
| 获取用户信息 | GET | `/api/auth/user-info` | 获取当前用户信息 |
| 刷新令牌 | POST | `/api/auth/refresh` | 刷新访问令牌 |

## 🎯 页面路由

| 路径 | 页面 | 描述 |
|------|------|------|
| `/user/login` | 登录页 | 用户登录 |
| `/user/register` | 注册页 | 用户注册 |
| `/account/change-password` | 修改密码 | 修改用户密码 |
| `/welcome` | 欢迎页 | 系统首页 |

## 🔐 认证机制

- 使用 JWT Token 进行身份验证
- Token 自动存储在 localStorage 中
- 请求拦截器自动添加 Authorization 头
- Token 过期时自动跳转到登录页

## 🚧 开发说明

### 添加新页面

1. 在 `src/pages` 目录下创建页面组件
2. 在 `config/routes.ts` 中添加路由配置
3. 如需要权限控制，在 `src/access.ts` 中配置权限

### 添加新接口

1. 在 `src/services/ant-design-pro/api.ts` 中添加接口函数
2. 在 `src/services/ant-design-pro/typings.d.ts` 中添加类型定义

### 添加国际化文案

1. **简体中文：** `src/locales/zh-CN/` 目录下对应文件
2. **繁体中文：** `src/locales/zh-TW/` 目录下对应文件  
3. **英文：** `src/locales/en-US/` 目录下对应文件

### 自定义样式

- 全局样式：`src/global.less` 或 `src/global.style.ts`
- 组件样式：使用 `antd-style` 的 `createStyles`

## 🐛 常见问题

### 1. 启动时端口被占用

```bash
npm start -- --port 3000  # 使用不同端口
```

### 2. 代理配置不生效

确保后端服务器已启动，检查 `config/proxy.ts` 中的地址配置。

### 3. 登录后页面空白

检查后端接口返回的用户信息格式是否正确，特别是 `display_name` 字段。

### 4. 语言切换不生效

1. 检查浏览器是否缓存了语言设置
2. 尝试清除浏览器缓存或使用无痕模式
3. 检查 `config/config.ts` 中的国际化配置

### 5. 头像不显示

1. 检查用户信息中是否包含 `display_name` 或 `username` 字段
2. 检查头像 URL 是否有效
3. 检查网络请求是否成功获取用户信息

## 📝 更新日志

### v1.1.0 (2024-12-19)
- ✅ 添加用户头像下拉菜单功能
- ✅ 完整的国际化支持（简体中文、繁体中文、英文）
- ✅ 优化用户体验和界面设计

### v1.0.0 (2024-12-19)
- ✅ 完成用户登录、注册、修改密码功能
- ✅ 对接后端 JWT 认证
- ✅ 响应式设计优化

## 📞 技术支持

如有问题请联系开发团队或查看项目文档。 