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
```

### 2. 上传文件

```bash
# 基本上传
curl -X POST http://localhost:8080 -F "file=@local-file"

# 带参数上传（指定目录）
curl -X POST http://localhost:8080 -F "file=@local-file" -F "dir=upload"
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
| -user | "" | 认证用户名（为空表示不需要认证） |
| -pass | "" | 认证密码（为空表示不需要认证） |

## 环境变量

| 环境变量 | 对应参数 | 说明 |
|---------|---------|------|
| WEB_PORT | -port | 服务器监听端口 |
| WEB_PATH | -path | 文件存储路径 |
| AUTH_USER | -user | 认证用户名 |
| AUTH_PASS | -pass | 认证密码 |

## 认证功能

GoFileserver 支持基本的 HTTP 认证功能，可以通过以下方式启用：

### 命令行方式
```bash
./GoFileserver -port 8080 -path ./www -user admin -pass password123
```

### Docker 方式
```bash
docker run -d -p 8080:8080 -v ./www:/web/www \
  -e AUTH_USER=admin -e AUTH_PASS=password123 \
  registry.cn-hangzhou.aliyuncs.com/public_hjj_images/go-fileserver:latest
```

### 访问带认证的服务

使用 curl 访问需要认证的服务：
```bash
# 下载文件
curl -u admin:password123 "http://localhost:8080/path/to/file" -o local-file

# 上传文件
curl -u admin:password123 -X POST http://localhost:8080 \
  -F "file=@local-file" -F "dir=upload"
```

## 示例

### 上传文件到指定目录

```bash
curl -X POST http://localhost:8080 -F "file=@test.txt" -F "dir=docs"
```

### 下载文件

```bash
curl http://localhost:8080/docs/test.txt -o test.txt
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

### 快速部署（Ingress + 持久化存储）

```bash
# 应用部署配置（包含 PV、PVC 和 Ingress）
kubectl apply -f k8s-deployment-ingress.yaml

# 查看部署状态
kubectl get pv
kubectl get pvc
kubectl get deployments
kubectl get pods
kubectl get services
kubectl get ingress
```

### 访问服务

#### NodePort 访问
服务将通过 NodePort 暴露在 `30080` 端口：

```bash
# 访问地址
http://<节点IP>:30080
```

#### Ingress 访问
通过 Ingress 配置，可以使用域名访问服务：

1. **配置域名解析**：
   - 将 `fileserver.example.com` 解析到您的 Kubernetes 集群入口 IP

2. **访问地址**：
   ```bash
   http://fileserver.example.com
   ```

3. **自定义域名**：
   - 编辑 `k8s-deployment-ingress.yaml` 文件，将 `fileserver.example.com` 改为您自己的域名

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
- **存储**：使用 PersistentVolume 和 PersistentVolumeClaim
  - PV：10Gi 存储空间，hostPath 类型
  - PVC：5Gi 存储请求
  - 存储路径：`/data/gofileserver`

#### Ingress + 持久化存储配置（k8s-deployment-ingress.yaml）
- **副本数**：1
- **容器端口**：8080
- **服务类型**：ClusterIP
- **资源限制**：CPU 1核，内存 512Mi
- **存储**：使用 PersistentVolume 和 PersistentVolumeClaim
  - PV：10Gi 存储空间，hostPath 类型
  - PVC：5Gi 存储请求
  - 存储路径：`/data/gofileserver`
- **Ingress 配置**：
  - 域名：`fileserver.example.com`
  - 路径：`/`
  - 后端服务：gofileserver:8080
  - 注解：
    - 使用 nginx ingress 控制器
    - 文件上传大小限制：100MB
    - 代理读写超时：10分钟（适合大文件传输）

## License

GPL-3.0 License
