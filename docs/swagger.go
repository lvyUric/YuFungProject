// Package docs 保险经纪管理系统API文档
//
//	@title			保险经纪管理系统API
//	@version		1.0
//	@description	基于Gin + MongoDB构建的保险经纪公司管理平台API文档
//	@description	支持用户认证、权限管理、保单管理等功能
//
//	@contact.name	开发团队
//	@contact.email	dev@insurance.com
//
//	@license.name	MIT
//	@license.url	https://opensource.org/licenses/MIT
//
//	@host						localhost:8080
//	@BasePath					/api
//	@schemes					http https
//
//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				JWT Token，格式: Bearer {token}
//
//	@tag.name			认证管理
//	@tag.description	用户登录、注册、密码管理等认证相关接口
//
//	@tag.name			用户管理
//	@tag.description	用户信息管理、用户列表查询等接口
//
//	@tag.name			公司管理
//	@tag.description	保险经纪公司管理接口
//
//	@tag.name			角色权限
//	@tag.description	角色管理和权限控制接口
//
//	@tag.name			保单管理
//	@tag.description	保单业务管理接口
//
//	@tag.name			系统管理
//	@tag.description	系统配置和管理接口
package docs
