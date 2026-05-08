# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## 项目概述

Vicomova 是一个基于 Go 的后端项目，使用字节跳动的 **Hertz** (HTTP 框架) 和 **Kitex** (RPC 框架) 构建微服务架构，采用 DDD（领域驱动设计）架构。

### 架构图

```
Client → HTTP Gateway (Hertz:8080) → User RPC (Kitex:8888)
                                          → Video RPC (Kitex:8889)
                                          → Interaction RPC (Kitex:8890)
                                              ↓
                                      MySQL / Redis / Kafka / Etcd
```

### 配置中心

- **Etcd**: 配置中心，配置热更新
- 配置存储在 `/vicomova/config` key
- 配置文件位于 `docker/config/base.yaml`
- 配置文件变更通过 etcd watch 实时推送到各服务

**支持热更新的配置**：database、redis
**需要重启生效**：jwt.secret、service.addr、email

### Go 版本
```
Go 1.25.0
```

---

## 常用命令

### tmux 开发模式
```bash
make dev          # 启动 tmux 会话（main/user/video/interaction/gateway 五个窗口）
make dev-stop     # 停止 tmux 会话
make dev-logs     # 查看各服务日志
```

### Docker 开发环境
```bash
make docker-up        # 启动 MySQL + Redis + Etcd
make docker-up-prod   # 生产环境 Docker 启动
make docker-down      # 停止服务
make docker-clean     # 清理数据卷
```

### 编译运行
```bash
make build           # 编译所有服务到 bin/
ETCD_ADDR=127.0.0.1:2379 ./bin/user        # 运行 User RPC
ETCD_ADDR=127.0.0.1:2379 ./bin/video        # 运行 Video RPC
ETCD_ADDR=127.0.0.1:2379 ./bin/interaction  # 运行 Interaction RPC
ETCD_ADDR=127.0.0.1:2379 ./bin/gateway       # 运行 Gateway
```

### 代码检查与测试
```bash
make check   # 运行 lint + test
make lint    # 运行 golangci-lint
make test    # 运行测试
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
  ├── user/            # User RPC 服务 (Kitex)
  ├── video/           # Video RPC 服务 (Kitex)
  └── interaction/      # Interaction RPC 服务 (Kitex)

api/rpc/               # Protobuf IDL 定义
  └── rpc/            # Proto 文件目录
    ├── user/user.proto
    ├── video/video.proto
    └── interaction/interaction.proto

internal/              # 内部业务代码
  ├── user/           # 用户服务（完整 DDD）
  │   ├── domain/
  │   │   ├── entity/
  │   │   ├── repository/
  │   │   ├── service/
  │   │   └── valueobject/
  │   ├── application/
  │   │   ├── command/
  │   │   └── query/
  │   ├── infrastructure/
  │   │   ├── persistence/mysql/
  │   │   ├── persistence/redis/
  │   │   └── external/email/
  │   ├── interfaces/
  │   │   ├── http/
  │   │   └── grpc/
  │   └── wire/        # 依赖注入
  ├── video/           # 视频服务（完整 DDD）
  │   ├── domain/
  │   │   ├── entity/   # Video, Category
  │   │   ├── repository/
  │   │   └── service/  # HotAlgorithm
  │   ├── application/
  │   │   ├── command/  # PublishVideo
  │   │   └── query/    # GetVideo, ListByCategory, ListHot
  │   ├── infrastructure/
  │   │   ├── persistence/mysql/
  │   │   ├── storage/  # VideoStorage (策略模式: local/qiniu/aliyun/aws)
  │   │   └── service/  # WilsonHotAlgorithm
  │   ├── interfaces/
  │   └── wire/
  ├── interaction/      # 互动服务（完整 DDD）
  │   ├── domain/
  │   │   ├── entity/   # Like, Comment, Favorite
  │   │   └── repository/
  │   ├── application/
  │   │   ├── command/  # Like, Comment, Favorite
  │   │   └── query/
  │   ├── infrastructure/
  │   │   └── persistence/mysql/
  │   ├── interfaces/
  │   └── wire/
  └── gateway/         # HTTP 网关
      ├── bootstrap.go
      └── router.go

pkg/                   # 公共工具
  ├── config/          # 配置加载 (Etcd) + 热更新
  ├── infrastructure/   # 基础设施 (MySQL/Redis/Kafka/Hertz/Kitex)
  ├── utils/           # 工具函数 (hasher/code)
  └── log/             # 日志封装

third_party/           # 第三方代码
  └── kitex_gen/       # 生成的 RPC 代码

script/                # 脚本
  ├── bootstrap.sh
  └── proto/           # Proto 生成脚本

docs/                  # 文档（Swagger）

docker/                # Docker 配置
  ├── config/          # 配置文件 (base.yaml)
  ├── env/             # 环境变量文件
  ├── data/            # 数据卷目录
  └── script/          # 容器启动脚本
```

---

## 核心设计

### 配置热更新

配置存储在 Etcd 中，通过 `pkg/config` 包管理：

- `config.Init(serviceName)` - 初始化配置，连接 etcd 加载配置
- `config.Get()` - 获取全局配置
- `config.RegisterCallback()` - 注册配置变更回调
- `config.Close()` - 关闭 etcd 连接

**支持热更新的配置**：
| 配置 | 回调动作 |
|------|----------|
| database | 重建 MySQL 连接池 |
| redis | 重建 Redis 连接池 |

**需要重启生效的配置**：
| 配置 | 说明 |
|------|------|
| jwt.secret | 变更后所有用户 Token 失效 |
| service.addr | 服务监听地址，无法热切换 |
| email | 邮件服务配置 |

### 存储策略模式

视频存储支持多种后端切换：

| 类型 | 说明 |
|------|------|
| `local` | 本地磁盘 |
| `qiniu` | 七牛云存储 |
| `aliyun` | 阿里云 OSS |
| `aws` | AWS S3 |

```go
type VideoStorage interface {
    Upload(ctx context.Context, key string, r io.Reader) (string, error)
    Delete(ctx context.Context, key string) error
    GetURL(ctx context.Context, key string) (string, error)
}
```

### 热度算法（Wilson 区间）

视频热度采用 Wilson 置信区间算法，避免头部效应：

```
score = (p + z²/2n - z*√((p(1-p)+z²/4n)/n)) / (1+z²/n)
```

- `p`: 点赞率（likes / views）
- `n`: 总浏览量
- `z`: 置信度参数（默认 1.96）

### 双 Token 认证

采用 Access Token + Refresh Token 旋转方案：

| Token | 过期时间 | 存储 |
|-------|---------|------|
| Access Token | 30 分钟 | 客户端内存 |
| Refresh Token | 7 天 | Redis 缓存 |

**认证流程**:
1. **注册**: 发送邮箱验证码 → 验证并创建用户
2. **登录**: 用户名+密码 → 返回 Access Token + Refresh Token
3. **访问**: 带 Access Token → 验证 → 返回数据
4. **刷新**: Refresh Token 有效 → 旋转更新 → 返回新 Token
5. **登出**: 使 Redis 中的 Refresh Token 失效

---

## 配置文件

- 配置模板: `docker/config/base.yaml.example`
- 实际配置: `docker/config/base.yaml`（已 gitignore）

### 配置结构

```yaml
database:
  host: "localhost"
  port: 3306
  user: "vicomova"
  password: "password"
  dbname: "vicomova"
  max_open_conns: 100
  max_idle_conns: 10

redis:
  host: "localhost"
  port: 6379
  password: ""
  db: 0

services:
  user:
    name: "user"
    addr: "127.0.0.1:8888"
  video:
    name: "video"
    addr: "127.0.0.1:8889"
  interaction:
    name: "interaction"
    addr: "127.0.0.1:8890"
  gateway:
    name: "gateway"
    addr: "127.0.0.1:8080"
```

### 环境变量
```bash
ETCD_ADDR=127.0.0.1:2379  # etcd 地址（必需）
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

修改 `api/rpc/*.proto` 后运行:
```bash
make proto
```

---

## 依赖框架

| 类别 | 框架 | 版本 |
|------|------|------|
| HTTP | cloudwego/hertz | v0.10.4 |
| RPC | cloudwego/kitex | v0.16.1 |
| ORM | gorm.io/gorm + gorm.io/driver/mysql | v1.31.1 / v1.6.0 |
| Redis | redis/go-redis/v9 | v9.18.0 |
| Kafka | IBM/sarama | v1.47.0 |
| 配置 | go.etcd.io/etcd/client/v3 | v3.6.10 |
| JWT | golang-jwt/jwt/v5 | v5.3.1 |
| 邮件 | github.com/alibabacloud-go/dm-20151123/v2 | v2.9.0 |
| 凭据 | github.com/aliyun/credentials-go | v1.4.12 |
| Swagger | github.com/swaggo/swag | v1.16.6 |

---

## CI/CD

GitHub Actions 自动构建：

- **dev 分支 push**: 运行 lint + build
- **PR to main**: 运行 lint + build + 测试 + Docker 构建

---

## 新增服务模板

新增服务时，在 `internal/` 下创建独立的服务目录：

```
internal/{service}/
├── domain/
│   ├── entity/
│   ├── repository/
│   ├── service/
│   └── valueobject/
├── application/
│   ├── command/
│   └── query/
├── infrastructure/
│   ├── persistence/
│   └── external/
├── interfaces/
│   ├── http/
│   └── grpc/
└── wire/           # 依赖注入
```
