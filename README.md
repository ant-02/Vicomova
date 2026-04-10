# Vicomova

基于 Go + Hertz + Kitex 的微服务后端项目。

## 架构

```
Client → HTTP Gateway (Hertz:8080) → User RPC (Kitex:8888)
                                              ↓
                                      MySQL / Redis / Kafka
```

## 项目结构

```
├── cmd/                      # 服务入口
│   ├── gateway/             # HTTP 网关 (Hertz)
│   └── user/                # User RPC 服务 (Kitex)
│
├── api/rpc/user/            # Protobuf IDL 定义
│   └── user.proto           # RPC 接口定义
│
├── internal/                # 业务代码
│   ├── domain/user/         # 领域实体
│   ├── repository/          # 仓储层
│   ├── service/             # 业务逻辑
│   ├── data/                # 数据访问 (MySQL/Redis/Kafka)
│   ├── rpc/user/            # RPC Handler + Client
│   └── gateway/             # HTTP Handler + Router
│
├── pkg/                     # 公共工具
│   ├── config/              # 配置加载 (Viper)
│   ├── hertz/               # Hertz 封装
│   └── kitex/               # Kitex 封装
│
├── docker/                  # Docker 配置
│   ├── docker-compose.yaml       # 开发环境 (仅 MySQL/Redis)
│   ├── docker-compose.prod.yaml  # 生产环境 (完整服务)
│   └── Dockerfile.*         # 多阶段构建
│
├── config/                  # 配置文件
├── docs/                    # Swagger 文档
└── Makefile                # 构建命令
```

## 快速开始

### 开发环境

```bash
# 启动 MySQL 和 Redis
make docker-up

# 运行 User RPC 服务
make run-user

# 新终端运行 Gateway
make run-gateway
```

### Docker 部署

```bash
# 构建镜像
make docker-build

# 启动所有服务
make docker-up-prod
```

## 常用命令

| 命令 | 用途 |
|------|------|
| `make run-user` | 运行 User RPC |
| `make run-gateway` | 运行 Gateway |
| `make docker-up` | 启动开发环境 (MySQL/Redis) |
| `make docker-up-prod` | 启动生产环境 |
| `make docker-down` | 停止服务 |
| `make docker-clean` | 清理数据卷 |
| `make proto` | 重新生成 RPC 代码 |
| `make swagger` | 生成 API 文档 |

## 微服务分层

请求链路: **Gateway Handler → RPC Client → RPC Handler → Service → Repository → DAO**

| 层级 | 文件位置 | 职责 |
|------|----------|------|
| Domain | `internal/domain/user/entity.go` | User 实体，GORM tag 映射 |
| Repository 接口 | `internal/repository/user_interface.go` | 定义数据操作接口 |
| Repository 实现 | `internal/repository/user.go` | 调用 DAO |
| DAO | `internal/data/mysql/dao/user.go` | GORM 数据库操作 |
| Service | `internal/service/user.go` | 业务逻辑、密码哈希、JWT |
| RPC Handler | `internal/rpc/user/handler.go` | RPC 请求处理 |
| Gateway Handler | `internal/gateway/handler/user.go` | HTTP 请求处理 |

## API 文档

启动服务后访问: http://localhost:8080/swagger/

## 依赖框架

- **HTTP**: cloudwego/hertz
- **RPC**: cloudwego/kitex
- **ORM**: gorm.io/gorm + gorm.io/driver/mysql
- **Redis**: redis/go-redis/v9
- **Kafka**: IBM/sarama
- **配置**: spf13/viper
- **JWT**: golang-jwt/jwt/v5
