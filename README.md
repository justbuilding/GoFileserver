# GoFileserver

基于 Go 语言实现的轻量级文件服务器，支持文件上传下载和参数传递。

## 功能特性

- ✅ 静态文件下载
- ✅ 文件上传功能
- ✅ 带参数上传下载
- ✅ 自动创建目录
- ✅ 多架构支持（amd64/arm64）
- ✅ 无依赖运行

## 快速开始

### 编译运行

```bash
# 编译
cd GoFileserver
go build -o GoFileserver

# 运行
./GoFileserver -port 8080 -path ./www
```

### Docker 部署

```yaml
# docker-compose.yml
version: '3'
services:
  gofileserver:
    build: .
    container_name: GoFileserver
    restart: unless-stopped
    ports:
      - '8080:8080'
    volumes:
      - ./www:/web/www
```

## 使用方法

### 1. 下载文件

```bash
# 基本下载
curl http://localhost:8080/path/to/file -o local-file

# 带参数下载
curl "http://localhost:8080/path/to/file?token=123&version=1.0" -o local-file
```

### 2. 上传文件

```bash
# 基本上传
curl -X POST http://localhost:8080 -F "file=@local-file"

# 带参数上传（指定目录和token）
curl -X POST http://localhost:8080 -F "file=@local-file" -F "dir=upload" -F "token=123"
```

## 多架构编译

```bash
# 创建输出目录
mkdir -p output

# 编译 amd64 版本
GOOS=linux GOARCH=amd64 go build -o output/GoFileserver_linux_amd64

# 编译 arm64 版本
GOOS=linux GOARCH=arm64 go build -o output/GoFileserver_linux_arm64
```

## 配置选项

| 参数 | 默认值 | 说明 |
|------|--------|------|
| -port | 8080 | 服务器监听端口 |
| -path | ./www | 文件存储路径 |
| -c | default | 配置方式（default/env） |

## 环境变量

| 环境变量 | 对应参数 | 说明 |
|---------|---------|------|
| WEB_PORT | -port | 服务器监听端口 |
| WEB_PATH | -path | 文件存储路径 |

## 示例

### 上传文件到指定目录

```bash
curl -X POST http://localhost:8080 -F "file=@test.txt" -F "dir=docs" -F "token=secret123"
```

### 带版本参数下载文件

```bash
curl "http://localhost:8080/docs/test.txt?token=secret123&version=1.0.0" -o test.txt
```

## Kubernetes 部署

### 快速部署（临时存储）

```bash
# 应用部署配置
kubectl apply -f k8s-deployment.yaml

# 查看部署状态
kubectl get deployments
kubectl get pods
kubectl get services
```

### 快速部署（持久化存储）

```bash
# 应用部署配置（包含 PV 和 PVC）
kubectl apply -f k8s-deployment-pvc.yaml

# 查看部署状态
kubectl get pv
kubectl get pvc
kubectl get deployments
kubectl get pods
kubectl get services
```

### 访问服务

服务将通过 NodePort 暴露在 `30080` 端口：

```bash
# 访问地址
http://<节点IP>:30080
```

### 配置说明

#### 临时存储配置（k8s-deployment.yaml）
- **副本数**：1
- **容器端口**：8080
- **NodePort**：30080
- **资源限制**：CPU 1核，内存 512Mi
- **健康检查**： readiness 和 liveness 探针
- **存储**：使用 emptyDir 临时存储

#### 持久化存储配置（k8s-deployment-pvc.yaml）
- **副本数**：1
- **容器端口**：8080
- **NodePort**：30080
- **资源限制**：CPU 1核，内存 512Mi
- **健康检查**： readiness 和 liveness 探针
- **存储**：使用 PersistentVolume 和 PersistentVolumeClaim
  - PV：10Gi 存储空间，hostPath 类型
  - PVC：5Gi 存储请求
  - 存储路径：`/data/gofileserver`

## License

GPL-3.0 License
