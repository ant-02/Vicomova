# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Vicomova 是一个基于 Go 的后端项目，使用字节跳动的 **Hertz** (HTTP 框架) 和 **Kitex** (RPC 框架) 构建微服务架构，采用 DDD（领域驱动设计）架构。

### 架构图

```
Client → HTTP Gateway (Hertz:8080) → User RPC (Kitex:8888)
                                              ↓
                                      MySQL / Redis / Kafka / Etcd
```

### 服务发现

- **Etcd**: 服务注册与发现，配置热更新
- 服务启动时注册到 etcd，Gateway 从 etcd 发现服务地址
- 配置文件变更通过 etcd watch 实时推送到各服务

---

## 常用命令

### tmux 开发模式
```bash
make dev      # 启动 tmux 会话（main/user/gateway 三个窗口）
make dev-stop # 停止 tmux 会话
```

### Docker 开发环境
```bash
make docker-up      # 启动 MySQL + Redis + Etcd
make docker-down    # 停止服务
make docker-clean   # 清理数据卷
```

### 直接运行
```bash
make run-user     # 运行 User RPC 服务
make run-gateway  # 运行 Gateway
```

### 代码生成
```bash
make proto   # 重新生成 RPC 代码
make swagger # 生成 Swagger 文档
```

---

## 项目结构

```
cmd/                    # 服务入口
  ├── gateway/         # HTTP 网关 (Hertz)
  └── user/            # User RPC 服务 (Kitex)

api/rpc/user/          # Protobuf IDL 定义
  └── user.proto       # RPC 接口定义

internal/              # 内部业务代码
  ├── user/            # 用户服务（完整 DDD）
  │   ├── domain/      # 领域层
  │   │   ├── entity/ # 实体
  │   │   ├── repository/ # 仓储接口
  │   │   ├── service/  # 领域服务
  │   │   └── valueobject/ # 值对象
  │   ├── application/ # 应用层
  │   │   ├── command/ # 写操作
  │   │   └── query/   # 读操作
  │   ├── infrastructure/ # 基础设施层
  │   │   ├── persistence/mysql/  # MySQL 实现
  │   │   ├── persistence/redis/   # Redis 实现
  │   │   └── external/email/      # 邮件服务（阿里云 DirectMail）
  │   └── interfaces/  # 接口层
  │       ├── http/   # HTTP Handler + Router
  │       └── grpc/   # RPC Handler + Client
  ├── shared/          # 跨服务共享
  │   ├── infrastructure/data/ # MySQL/Redis/Kafka 连接
  │   └── pkg/        # errors/log/constants
  └── gateway/         # HTTP 网关
      ├── bootstrap.go # 启动初始化
      └── dynamic_router.go # 动态路由

pkg/                   # 公共工具
  ├── config/         # 配置加载 (Viper) + 热更新
  ├── etcd/           # Etcd 客户端封装
  │   ├── client.go   # Etcd 连接
  │   ├── registry.go # 服务注册（带租约）
  │   ├── discovery.go # 服务发现
  │   └── config.go   # 配置
  └── ...

docker/                # Docker 配置
config/                # 配置文件（已 gitignore）
  └── base.yaml.example # 配置模板
third_party/kitex_gen/ # 生成的 RPC 代码
```

---

## 配置文件

**注意**: `config/base.yaml` 包含敏感信息，已从 git 跟踪中移除。

- 配置模板: `config/base.yaml.example`
- Docker 配置: `docker/config/base.yaml`（已 gitignore）

### 配置项说明

| 配置 | 说明 |
|------|------|
| `email.access_key` | 阿里云 AccessKey ID（建议使用环境变量） |
| `email.access_secret` | 阿里云 AccessKey Secret（建议使用环境变量） |
| `email.account_name` | 发件地址，如 `noreply@mail.xhhx.xyz` |
| `jwt.secret` | JWT 密钥（生产环境必须修改） |

配置通过 **Viper** 加载，支持命令行 flag 覆盖。

### 启动参数
```bash
./user -config config/base.yaml -port 8888
./gateway -config config/base.yaml -port 8080
```

---

## 微服务分层

请求链路: **Gateway Handler → RPC Client → RPC Handler → Application Service → Domain Service → Repository**

| 层级 | 文件位置 | 职责 |
|------|----------|------|
| Interface | `internal/*/interfaces/` | 适配器层，处理请求/响应 |
| Application | `internal/*/application/` | 应用服务，编排业务用例 |
| Domain | `internal/*/domain/` | 领域实体、业务规则、领域服务 |
| Infrastructure | `internal/*/infrastructure/` | 持久化、缓存、外部服务 |

---

## Proto 代码生成

修改 `api/rpc/user/user.proto` 后运行:
```bash
make proto
```

---

## 依赖框架

- **HTTP**: cloudwego/hertz v0.10.4
- **RPC**: cloudwego/kitex v0.16.1
- **ORM**: gorm.io/gorm + gorm.io/driver/mysql
- **Redis**: redis/go-redis/v9
- **Kafka**: IBM/sarama
- **配置**: spf13/viper
- **JWT**: golang-jwt/jwt/v5
- **Etcd**: go.etcd.io/etcd/client/v3
- **邮件**: github.com/alibabacloud-go/dm-20151123/v2
- **凭据**: github.com/aliyun/credentials-go
