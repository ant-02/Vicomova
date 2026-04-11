# Vicomova

基于 Go + Hertz + Kitex 的微服务后端项目，采用 DDD（Domain-Driven Design）架构设计。

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
├── internal/                 # 业务代码 (DDD)
│   ├── user/               # 用户服务（完整 DDD）
│   │   ├── domain/         # 领域层
│   │   │   ├── entity/     # 实体
│   │   │   ├── repository/ # 仓储接口
│   │   │   ├── service/    # 领域服务
│   │   │   ├── valueobject/# 值对象
│   │   │   ├── event/      # 领域事件
│   │   │   └── errors.go   # 领域错误
│   │   ├── application/    # 应用层
│   │   │   ├── command/    # 写操作（注册/登录等）
│   │   │   └── query/      # 读操作
│   │   ├── infrastructure/  # 基础设施层
│   │   │   ├── persistence/mysql/  # MySQL 实现
│   │   │   ├── persistence/redis/  # Redis 实现
│   │   │   └── external/email/     # 邮件服务
│   │   └── interfaces/     # 接口层
│   │       ├── http/       # HTTP Handler + Router
│   │       └── grpc/       # RPC Handler + Client
│   ├── shared/              # 跨服务共享
│   │   ├── infrastructure/data/  # MySQL/Redis/Kafka 连接
│   │   └── pkg/            # errors/log/middleware
│   └── wire/wire.go        # 依赖注入
│
├── pkg/                     # 公共工具
│   ├── config/             # 配置加载 (Viper)
│   ├── hertz/              # Hertz 封装
│   └── kitex/              # Kitex 封装
│
├── docker/                  # Docker 配置
├── config/                  # 配置文件
├── third_party/kitex_gen/   # 生成的 RPC 代码
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

## 新增服务

新增服务（如 order）时，在 `internal/` 下创建独立的服务目录：

```
internal/order/              # 新服务
├── domain/
│   ├── entity/
│   ├── repository/
│   ├── service/
│   └── errors.go
├── application/
│   ├── command/
│   └── query/
├── infrastructure/
│   ├── persistence/
│   └── external/
└── interfaces/
    ├── http/
    └── grpc/
```

## DDD 分层架构

请求链路: **Gateway Handler → RPC Client → RPC Handler → Application Service → Domain Service → Repository**

| 层级 | 目录 | 职责 |
|------|------|------|
| Interface | `internal/*/interfaces/` | 适配器层，处理请求/响应 |
| Application | `internal/*/application/` | 应用服务，编排业务用例 |
| Domain | `internal/*/domain/` | 领域实体、业务规则、领域服务 |
| Infrastructure | `internal/*/infrastructure/` | 持久化、缓存、外部服务 |

## 认证机制

采用双 Token 认证策略：

| Token | 过期时间 | 存储 |
|-------|---------|------|
| Access Token | 30 分钟 | 客户端内存 |
| Refresh Token | 7 天 | Redis 缓存 |

### 认证流程

1. **注册**: 发送邮箱验证码 → 验证并创建用户
2. **登录**: 用户名+密码 → 返回 Access Token + Refresh Token
3. **访问**: 带 Access Token → 验证 → 返回数据
4. **刷新**: Refresh Token 有效 → 旋转更新 → 返回新 Access Token + 新 Refresh Token
5. **登出**: 使 Redis 中的 Refresh Token 失效

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
