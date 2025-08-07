#!/bin/bash

# 昱丰保险经纪管理系统 Docker 部署脚本
# 作者: AI Assistant
# 日期: $(date +%Y-%m-%d)

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 日志函数
log_info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

log_warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

log_step() {
    echo -e "${BLUE}[STEP]${NC} $1"
}

# 检查Docker是否安装
check_docker() {
    log_step "检查Docker环境..."
    if ! command -v docker &> /dev/null; then
        log_error "Docker未安装，请先安装Docker"
        exit 1
    fi
    
    if ! command -v docker-compose &> /dev/null; then
        log_error "Docker Compose未安装，请先安装Docker Compose"
        exit 1
    fi
    
    log_info "Docker环境检查通过"
}

# 创建必要的目录
create_directories() {
    log_step "创建必要的目录..."
    mkdir -p logs uploads
    log_info "目录创建完成"
}

# 停止现有容器
stop_containers() {
    log_step "停止现有容器..."
    docker-compose -f docker-compose.prod.yml down --remove-orphans || true
    log_info "现有容器已停止"
}

# 构建镜像
build_images() {
    log_step "构建Docker镜像..."
    docker-compose -f docker-compose.prod.yml build --no-cache
    log_info "镜像构建完成"
}

# 启动服务
start_services() {
    log_step "启动服务..."
    docker-compose -f docker-compose.prod.yml up -d
    log_info "服务启动完成"
}

# 检查服务状态
check_services() {
    log_step "检查服务状态..."
    sleep 10
    
    # 检查后端服务
    if curl -f http://localhost:8088/health > /dev/null 2>&1; then
        log_info "✅ 后端服务运行正常"
    else
        log_warn "⚠️  后端服务可能未完全启动，请稍后检查"
    fi
    
    # 检查前端服务
    if curl -f http://localhost:80/health > /dev/null 2>&1; then
        log_info "✅ 前端服务运行正常"
    else
        log_warn "⚠️  前端服务可能未完全启动，请稍后检查"
    fi
}

# 显示服务信息
show_info() {
    log_step "部署完成！服务信息如下："
    echo ""
    echo "🌐 前端访问地址: http://localhost"
    echo "🔧 后端API地址: http://localhost:8088"
    echo "📖 API文档地址: http://localhost:8088/swagger/index.html"
    echo "💾 数据库地址: mongodb://106.52.172.124:27017"
    echo ""
    echo "📋 常用命令："
    echo "  查看服务状态: docker-compose -f docker-compose.prod.yml ps"
    echo "  查看服务日志: docker-compose -f docker-compose.prod.yml logs -f"
    echo "  停止服务: docker-compose -f docker-compose.prod.yml down"
    echo "  重启服务: docker-compose -f docker-compose.prod.yml restart"
    echo ""
}

# 主函数
main() {
    echo "🚀 开始部署昱丰保险经纪管理系统..."
    echo ""
    
    check_docker
    create_directories
    stop_containers
    build_images
    start_services
    check_services
    show_info
    
    log_info "部署完成！"
}

# 错误处理
trap 'log_error "部署过程中发生错误，请检查日志"; exit 1' ERR

# 执行主函数
main "$@" 