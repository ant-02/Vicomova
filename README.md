# Vicomova

基于 Go + Hertz + Kitex 的微服务后端项目，采用 DDD（领域驱动设计）架构，支持服务注册发现和配置热更新。

## 架构

```
Client → HTTP Gateway (Hertz:8080) → User RPC (Kitex:8888)
                                              ↓
                                      MySQL / Redis / Kafka / Etcd
```

### 核心特性

- **配置中心**: 基于 Etcd 的配置管理，配置热更新
- **服务发现**: 基于 Etcd 的服务注册与发现
- **双 Token 认证**: Access Token + Refresh Token 旋转方案
- **邮件服务**: 阿里云 DirectMail SDK
- **视频服务**: 投稿、播放、分类、热度算法（Wilson 区间）
- **互动服务**: 点赞、评论、收藏（可复用）
- **存储策略**: local / 七牛 / 阿里云 OSS / AWS S3

## 项目结构

```
├── cmd/                      # 服务入口
│   ├── gateway/             # HTTP 网关 (Hertz)
│   ├── user/                # User RPC 服务 (Kitex)
│   ├── video/               # Video RPC 服务 (Kitex)
│   └── interaction/         # Interaction RPC 服务 (Kitex)
│
├── api/rpc/                 # Protobuf IDL 定义
│   ├── user/user.proto
│   ├── video/video.proto
│   └── interaction/interaction.proto
│
├── internal/                 # 业务代码 (DDD)
│   ├── user/               # 用户服务
│   ├── video/              # 视频服务
│   ├── interaction/        # 互动服务
│   └── gateway/            # HTTP 网关
│
├── pkg/                     # 公共工具
│   ├── config/             # 配置加载 (Etcd) + 热更新
│   ├── infrastructure/    # 基础设施 (MySQL/Redis/Kafka/Hertz)
│   ├── utils/              # 工具函数
│   └── log/                # 日志封装
│
├── docker/                  # Docker 配置
│   ├── config/             # 配置文件 (base.yaml)
│   └── script/              # 容器启动脚本
│
├── third_party/kitex_gen/   # 生成的 RPC 代码
└── Makefile
```

## 快速开始

### 1. 启动开发环境

```bash
# 启动 MySQL + Redis + Etcd
make docker-up

# tmux 开发模式（五个窗口）
make dev
```

### 2. 直接运行

```bash
# 各服务独立运行（需先启动 etcd）
ETCD_ADDR=127.0.0.1:2379 ./bin/user
ETCD_ADDR=127.0.0.1:2379 ./bin/video
ETCD_ADDR=127.0.0.1:2379 ./bin/interaction
ETCD_ADDR=127.0.0.1:2379 ./bin/gateway
```

### 3. Docker 部署

```bash
# 构建镜像
make docker-build

# 启动生产环境
make docker-up-prod
```

## 常用命令

| 命令 | 用途 |
|------|------|
| `make dev` | tmux 开发模式 |
| `make dev-stop` | 停止 tmux |
| `make build` | 编译所有服务到 bin/ |
| `make run-user` | 运行 User RPC |
| `make run-video` | 运行 Video RPC |
| `make run-interaction` | 运行 Interaction RPC |
| `make run-gateway` | 运行 Gateway |
| `make docker-up` | 启动开发环境 |
| `make docker-up-prod` | 启动生产环境 |
| `make docker-down` | 停止服务 |
| `make docker-clean` | 清理数据卷 |
| `make proto` | 重新生成 RPC 代码 |
| `make swagger` | 生成 Swagger 文档 |

## 配置热更新

配置存储在 Etcd 中，通过 `docker/config/base.yaml` 上传到 `/vicomova/config` key。

**支持热更新的配置**：
- `database` - MySQL 连接配置
- `redis` - Redis 连接配置

**需要重启生效的配置**：
- `jwt.secret` - JWT 密钥（变更后所有用户 Token 失效）
- `service.addr` - 服务监听地址
- `email` - 邮件服务配置

修改 `docker/config/base.yaml` 后，配置会自动同步到 etcd 并热更新到各服务。

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
4. **刷新**: Refresh Token 有效 → 旋转更新 → 返回新 Token
5. **登出**: 使 Redis 中的 Refresh Token 失效

## 核心设计

### 存储策略模式

视频存储支持多种后端切换：

| 类型 | 说明 |
|------|------|
| `local` | 本地磁盘 |
| `qiniu` | 七牛云存储 |
| `aliyun` | 阿里云 OSS |
| `aws` | AWS S3 |

### 热度算法（Wilson 区间）

视频热度采用 Wilson 置信区间算法，避免头部效应：

```
score = (p + z²/2n - z*√((p(1-p)+z²/4n)/n)) / (1+z²/n)
```

## 新增服务

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

## CI/CD

GitHub Actions 自动构建：

- **dev 分支 push**: 运行 lint + build
- **PR to main**: 运行 lint + build + 测试 + Docker 构建

## 依赖框架

- **HTTP**: cloudwego/hertz v0.10.4
- **RPC**: cloudwego/kitex v0.16.1
- **ORM**: gorm.io/gorm + gorm.io/driver/mysql
- **Redis**: redis/go-redis/v9
- **Kafka**: IBM/sarama
- **配置**: go.etcd.io/etcd/client/v3
- **JWT**: golang-jwt/jwt/v5
- **邮件**: github.com/alibabacloud-go/dm-20151123/v2
- **凭据**: github.com/aliyun/credentials-go
