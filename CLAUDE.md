# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Vicomova 是一个基于 Go 的后端项目，使用字节跳动的 **Hertz** (HTTP 框架) 和 **Kitex** (RPC 框架) 构建微服务架构。

### 架构图

```
Client → HTTP Gateway (Hertz:8080) → User RPC (Kitex:8888)
                                              ↓
                                      MySQL / Redis / Kafka
```

---

## 常用命令

### 开发阶段
```bash
# 启动 MySQL 和 Redis (Docker)
make docker-up

# 运行 User RPC 服务
make run-user

# 新终端运行 Gateway
make run-gateway

# 重新生成 RPC 代码 (修改 proto 后)
make proto
```

### Docker 部署
```bash
# 构建镜像
make docker-build

# 启动所有服务
make docker-up

# 停止服务
make docker-down

# 清理（包含数据卷）
make docker-clean
```

### Go 命令
```bash
go mod tidy          # 整理依赖
go run ./cmd/user    # 直接运行 user 服务
go run ./cmd/gateway # 直接运行 gateway
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
  ├── domain/user/     # 领域实体
  ├── repository/      # 仓储层 (接口 + 实现)
  ├── service/         # 业务逻辑层
  ├── data/            # 数据访问层 (MySQL/Redis/Kafka)
  ├── rpc/user/        # RPC Handler + Client
  └── gateway/         # HTTP Handler + Router

pkg/                   # 公共导出包
  ├── config/          # 配置加载 (Viper)
  ├── hertz/           # Hertz 封装
  └── kitex/           # Kitex 封装

third_party/           # Kitex 生成的代码
  └── kitex_gen/user/  # Protobuf 生成的桩代码
```

---

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

---

## 配置文件

- 开发环境: `config/base.yaml`
- Docker 环境: `docker/config/base.yaml`

配置通过 **Viper** 加载，支持命令行 flag 覆盖和环境变量。

### 启动参数
```bash
./user -config config/base.yaml -port 8888
./gateway -config config/base.yaml -port 8080 -user_addr 127.0.0.1:8888
```

---

## Proto 代码生成

修改 `api/rpc/user/user.proto` 后运行:
```bash
sh script/proto/gen.sh
# 或直接
kitex -module vicomova -service user -gen-path third_party/kitex_gen ./api/rpc/user/user.proto
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
