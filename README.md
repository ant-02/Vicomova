# Vicomova

基于 Go + Hertz + Kitex 的微服务后端项目，采用 DDD（Domain-Driven Design）架构设计。

## 架构

```
Client → HTTP Gateway (Hertz:8080) → User RPC (Kitex:8888)
                                              ↓
                                      MySQL / Redis / Kafka
```

## DDD 项目结构

```
├── cmd/                      # 服务入口
│   ├── gateway/             # HTTP 网关 (Hertz)
│   └── user/                # User RPC 服务 (Kitex)
│
├── api/rpc/user/            # Protobuf IDL 定义
│   └── user.proto           # RPC 接口定义
│
├── internal/                # 业务代码 (DDD)
│   ├── domain/user/        # 领域层 (实体、值对象、领域服务)
│   ├── application/user/   # 应用层 (用例、命令/查询)
│   ├── infrastructure/      # 基础设施层 (持久化、RPC客户端)
│   │   └── persistence/
│   │       ├── mysql/      # MySQL 仓储实现
│   │       └── redis/     # Redis 缓存实现
│   └── interface/          # 接口层 (RPC Handler、Gateway Handler)
│       ├── rpc/
│       └── gateway/
│
├── pkg/                     # 公共工具
│   ├── config/             # 配置加载 (Viper)
│   ├── hertz/              # Hertz 封装
│   └── kitex/              # Kitex 封装
│
├── docker/                  # Docker 配置
│   ├── docker-compose.yaml       # 开发环境 (MySQL/Redis)
│   ├── docker-compose.prod.yaml  # 生产环境 (完整服务)
│   └── Dockerfile.*         # 多阶段构建
│
├── config/                  # 配置文件
├── .github/workflows/       # GitHub Actions CI/CD
├── Makefile                # 构建命令
└── README.md
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
| `make lint` | 代码检查 |
| `make build` | 编译所有服务 |

## DDD 分层架构

请求链路: **Gateway Handler → RPC Client → RPC Handler → Application Service → Domain Service → Repository**

| 层级 | 目录 | 职责 |
|------|------|------|
| Interface | `internal/interface/` | 适配器层，处理请求/响应 |
| Application | `internal/application/user/` | 应用服务，编排业务用例 |
| Domain | `internal/domain/user/` | 领域实体、业务规则、领域服务 |
| Infrastructure | `internal/infrastructure/` | 持久化、缓存、RPC 客户端等 |

### Domain 层

| 文件 | 类型 | 说明 |
|------|------|------|
| `entity.go` | 聚合根 | User 实体，GORM tag 映射 |
| `refresh_token.go` | 值对象 | RefreshToken 值对象 |
| `repository.go` | 接口 | Repository Port (仓储接口) |
| `service.go` | 领域服务 | Token 生成与验证 |
| `events.go` | 领域事件 | 领域事件定义 |
| `errors.go` | 错误 | 领域错误定义 |

### Application 层

| 文件 | 说明 |
|------|------|
| `command.go` | 写操作 (注册/登录/刷新Token/登出) |
| `query.go` | 读操作 (获取用户信息) |

## 认证机制

采用双 Token 认证策略：

| Token | 过期时间 | 存储 |
|-------|---------|------|
| Access Token | 30 分钟 | 客户端内存 |
| Refresh Token | 7 天 | Redis 缓存 |

### 认证流程

1. **登录**: 用户名+密码 → 返回 Access Token + Refresh Token
2. **访问**: 带 Access Token → 验证 → 返回数据
3. **刷新**: Refresh Token 有效 → 旋转更新 → 返回新 Access Token + 新 Refresh Token
4. **登出**: 使 Redis 中的 Refresh Token 失效

## CI/CD

GitHub Actions 自动构建：

- **dev 分支 push**: 运行 lint + build
- **PR to main**: 运行 lint + build + Docker 构建

## 依赖框架

- **HTTP**: cloudwego/hertz v0.10.4
- **RPC**: cloudwego/kitex v0.16.1
- **ORM**: gorm.io/gorm + gorm.io/driver/mysql
- **Redis**: redis/go-redis/v9
- **Kafka**: IBM/sarama
- **配置**: spf13/viper
- **JWT**: golang-jwt/jwt/v5
