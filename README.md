# GoIP - IP查询服务

GoIP是一个基于Golang的高性能IP查询服务，同时支持gRPC和HTTP REST API，提供IP地理位置查询功能。

## 功能特性

- 🚀 **高性能**: 基于Golang构建，支持并发查询
- 🔌 **多协议**: 同时支持gRPC和HTTP REST API
- 💾 **缓存机制**: 内置内存缓存，支持Redis扩展
- 🐳 **容器化**: 完整的Docker支持
- 📊 **监控**: 集成Prometheus监控
- 🔧 **配置灵活**: 支持YAML配置文件和环境变量
- 📝 **日志完善**: 结构化日志输出

## 项目结构

```
goip/
├── api/proto/          # gRPC协议定义
├── cmd/server/         # 服务端主程序
├── internal/           # 内部实现
│   ├── config/         # 配置管理
│   ├── handler/        # HTTP/gRPC处理器
│   ├── service/        # 核心业务逻辑
│   └── ipquery/        # IP查询实现
├── pkg/                # 公共库
├── configs/            # 配置文件
├── deployments/        # 部署配置
├── scripts/            # 脚本文件
├── Dockerfile          # Docker镜像构建
├── Makefile           # 构建脚本
└── README.md          # 项目文档
```

## 快速开始

### 环境要求

- Go 1.23+
- Docker (可选)

### 本地运行

1. 克隆项目
```bash
git clone https://github.com/ushell/goip.git
cd goip
```

2. 安装依赖
```bash
make deps
```

3. 运行服务
```bash
make run
```

### Docker运行

1. 构建镜像
```bash
make docker-build
```

2. 运行容器
```bash
make docker-run
```

3. 使用Docker Compose
```bash
make docker-compose
```

## API文档

### HTTP REST API

#### 查询单个IP
```bash
GET /api/v1/ip/{ip}
```

**示例请求:**
```bash
curl http://localhost:8080/api/v1/ip/8.8.8.8
```

**响应示例:**
```json
{
  "code": 0,
  "data": {
    "ip": "8.8.8.8",
    "country": "美国",
    "country_code": "US",
    "region": "加利福尼亚州",
    "city": "山景城",
    "district": "山景城区",
    "isp": "",
    "latitude": "",
    "longitude": "",
    "timezone": "",
    "postal_code": "94043",
    "is_valid": true
  }
}
```

#### 批量查询IP
```bash
POST /api/v1/ip/batch
```

**示例请求:**
```bash
curl -X POST http://localhost:8080/api/v1/ip/batch \
  -H "Content-Type: application/json" \
  -d '{"ips": ["8.8.8.8", "1.1.1.1"]}'
```

#### 获取客户端IP
```bash
GET /api/v1/ip/client
```

#### 健康检查
```bash
GET /api/v1/health
```

#### 服务状态
```bash
GET /api/v1/status
```

### gRPC API

#### 生成客户端代码
```bash
make proto
```

#### 服务定义
- `QueryIP` - 查询单个IP
- `BatchQueryIP` - 批量查询IP
- `GetServiceStatus` - 获取服务状态

## 配置说明

配置文件位于 `configs/config.yaml`，支持以下配置：

```yaml
server:
  http:
    host: "0.0.0.0"
    port: 8080
  grpc:
    host: "0.0.0.0"
    port: 8081

logging:
  level: "info"
  format: "json"

cache:
  enabled: true
  ttl: "1h"
```

## 开发指南

### 项目设置

1. 安装开发依赖
```bash
make deps
```

2. 运行测试
```bash
make test
```

3. 代码格式化
```bash
make fmt
```

4. 静态检查
```bash
make lint
```

### 添加新的IP查询源

1. 实现 `ipquery.QueryProvider` 接口
2. 在 `service.NewIPService` 中注册新的提供者

## 部署

### 单机部署

```bash
# 构建
make build

# 运行
./goip
```

## IP数据源

本项目使用 [**ip2region**](https://github.com/lionsoul2014/ip2region) 作为IP地理位置查询的数据源。

#### 数据源详情
- **数据库**: `ip2region.xdb` (11MB二进制格式)
- **版本**: v2.11.2
- **更新频率**: 支持自动重载，每24小时检查更新
- **数据格式**: 国家|区域|省份|城市|ISP

#### 数据覆盖范围
| 字段 | 说明 | 数据来源 |
|------|------|----------|
| 国家 | 国家名称 | ip2region |
| 省份 | 省份/州信息 | ip2region |
| 城市 | 城市信息 | ip2region |
| ISP | 网络服务商 | ip2region |
| 国家代码 | ISO国家代码 | 内部映射表 |
| 经纬度 | 地理坐标 | 暂不支持 |
| 时区 | 时区信息 | 暂不支持 |
| 邮编 | 邮政编码 | 暂不支持 |

#### 数据库配置
```yaml
ip_database:
  type: "local"
  path: "./data/ip2region.xdb"
  cache_size: 512  # MB
  auto_reload: true
  reload_interval: "24h"
```

#### 扩展支持
项目设计了 `QueryProvider` 接口，支持未来集成其他IP数据源：

```go
type QueryProvider interface {
    Query(ip string) (*IPInfo, error)
    BatchQuery(ips []string) ([]*IPInfo, error)
    Close() error
}
```