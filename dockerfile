# 构建 amd64 版本
FROM golang:alpine AS builder-amd64
WORKDIR /app
COPY . .
RUN apk add --no-cache git
RUN GOOS=linux GOARCH=amd64 go build -o GoFileserver

# 构建 arm64 版本
FROM golang:alpine AS builder-arm64
WORKDIR /app
COPY . .
RUN apk add --no-cache git
RUN GOOS=linux GOARCH=arm64 go build -o GoFileserver

# 最终镜像
FROM alpine:latest

ENV WEB_PORT=8080
ENV WEB_PATH=/web/www
ENV AUTH_USER=
ENV AUTH_PASS=
WORKDIR /web/

# 根据架构复制对应二进制文件
COPY --from=builder-amd64 /app/GoFileserver /web/GoFileserver-amd64
COPY --from=builder-arm64 /app/GoFileserver /web/GoFileserver-arm64

# 创建启动脚本
RUN mkdir -p -m 755 /web/www && \
    chmod 755 /web/GoFileserver-* && \
    echo '#!/bin/sh' > /web/start.sh && \
    echo 'ARCH=$(uname -m)' >> /web/start.sh && \
    echo 'if [ "$ARCH" = "x86_64" ]; then' >> /web/start.sh && \
    echo '    exec /web/GoFileserver-amd64 -c env' >> /web/start.sh && \
    echo 'else' >> /web/start.sh && \
    echo '    exec /web/GoFileserver-arm64 -c env' >> /web/start.sh && \
    echo 'fi' >> /web/start.sh && \
    chmod +x /web/start.sh && \
    rm -rf /var/cache/apk/*

VOLUME $WEB_PATH

EXPOSE $WEB_PORT

# 使用启动脚本
CMD ["/web/start.sh"]