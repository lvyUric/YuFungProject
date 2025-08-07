# 使用官方Go镜像作为构建环境
FROM golang:1.24-alpine AS builder

# 设置Go环境变量，使用国内代理
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
ENV GO111MODULE=on
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64

# 设置工作目录
WORKDIR /app

# 安装必要的系统依赖，使用国内镜像源
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk add --no-cache git ca-certificates tzdata

# 复制go mod文件
COPY go.mod go.sum ./

# 下载依赖，增加超时设置和详细输出
RUN go mod download -x || (echo "下载失败，重试..." && sleep 5 && go mod download -x)

# 复制源代码
COPY . .

# 构建应用，增加构建参数和优化
RUN go build \
    -ldflags="-w -s -extldflags=-static" \
    -a -installsuffix cgo \
    -o main ./cmd/main.go

# 使用轻量级的alpine镜像作为运行环境
FROM alpine:latest

# 使用国内镜像源安装依赖
RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
    apk --no-cache add ca-certificates tzdata curl wget netcat-openbsd

# 设置时区
ENV TZ=Asia/Shanghai

# 创建非root用户
RUN addgroup -g 1001 -S appgroup && \
    adduser -u 1001 -S appuser -G appgroup

# 创建必要的目录
RUN mkdir -p /app/logs /app/uploads && \
    chown -R appuser:appgroup /app

# 设置工作目录
WORKDIR /app

# 从builder阶段复制二进制文件
COPY --from=builder /app/main .

# 复制配置文件
COPY --from=builder /app/configs ./configs

# 复制网络测试脚本
COPY test-mongodb-from-container.sh /app/
RUN chmod +x /app/test-mongodb-from-container.sh

# 切换到非root用户
USER appuser

# 暴露端口
EXPOSE 8088

# 健康检查
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8088/health || exit 1

# 启动应用
CMD ["./main"] 