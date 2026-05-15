.PHONY: dev dev-stop dev-logs run-user run-gateway run-video run-interaction run-chat run-commerce docker-build docker-up docker-up-prod docker-down docker-clean proto swagger goimports vet test check

SESSION_NAME = vicomova

# tmux 开发模式：6 个窗口
dev:
	@if tmux has-session -t $(SESSION_NAME) 2>/dev/null; then \
		tmux kill-session -t $(SESSION_NAME); \
	fi
	@echo "Starting Vicomova services..."
	@tmux new-session -d -s $(SESSION_NAME) -n main
	@tmux new-window -t $(SESSION_NAME) -n user
	@tmux new-window -t $(SESSION_NAME) -n video
	@tmux new-window -t $(SESSION_NAME) -n interaction
	@tmux new-window -t $(SESSION_NAME) -n chat
	@tmux new-window -t $(SESSION_NAME) -n commerce
	@tmux new-window -t $(SESSION_NAME) -n gateway
	@tmux send-keys -t $(SESSION_NAME):main "echo 'Vicomova dev session'" C-m
	@tmux send-keys -t $(SESSION_NAME):user "make run-user" C-m
	@tmux send-keys -t $(SESSION_NAME):video "make run-video" C-m
	@tmux send-keys -t $(SESSION_NAME):interaction "make run-interaction" C-m
	@tmux send-keys -t $(SESSION_NAME):chat "make run-chat" C-m
	@tmux send-keys -t $(SESSION_NAME):commerce "make run-commerce" C-m
	@tmux send-keys -t $(SESSION_NAME):gateway "make run-gateway" C-m
	@tmux attach-session -t $(SESSION_NAME)

# 停止
dev-stop:
	@if tmux has-session -t $(SESSION_NAME) 2>/dev/null; then \
		tmux kill-session -t $(SESSION_NAME); \
		echo "Session killed"; \
	else \
		echo "Session does not exist"; \
	fi

# 查看服务日志
dev-logs:
	@echo "=== USER LOGS ===" && tmux capture-pane -t $(SESSION_NAME):user -p | tail -30
	@echo "=== VIDEO LOGS ===" && tmux capture-pane -t $(SESSION_NAME):video -p | tail -30
	@echo "=== INTERACTION LOGS ===" && tmux capture-pane -t $(SESSION_NAME):interaction -p | tail -30
	@echo "=== CHAT LOGS ===" && tmux capture-pane -t $(SESSION_NAME):chat -p | tail -30
	@echo "=== COMMERCE LOGS ===" && tmux capture-pane -t $(SESSION_NAME):commerce -p | tail -30
	@echo "=== GATEWAY LOGS ===" && tmux capture-pane -t $(SESSION_NAME):gateway -p | tail -30

# 编译输出到 bin/
build:
	mkdir -p bin
	go build -o bin/user ./cmd/user
	go build -o bin/video ./cmd/video
	go build -o bin/interaction ./cmd/interaction
	go build -o bin/chat ./cmd/chat
	go build -o bin/commerce ./cmd/commerce
	go build -o bin/gateway ./cmd/gateway

# 直接运行（编译后执行）
run-user:
	mkdir -p bin && go build -o bin/user ./cmd/user && set -a && . ./docker/env/base.env && set +a && ./bin/user

run-video:
	mkdir -p bin && go build -o bin/video ./cmd/video && set -a && . ./docker/env/base.env && set +a && ./bin/video

run-interaction:
	mkdir -p bin && go build -o bin/interaction ./cmd/interaction && set -a && . ./docker/env/base.env && set +a && ./bin/interaction

run-chat:
	mkdir -p bin && go build -o bin/chat ./cmd/chat && set -a && . ./docker/env/base.env && set +a && ./bin/chat

run-commerce:
	mkdir -p bin && go build -o bin/commerce ./cmd/commerce && set -a && . ./docker/env/base.env && set +a && ./bin/commerce

run-gateway:
	mkdir -p bin && go build -o bin/gateway ./cmd/gateway && set -a && . ./docker/env/base.env && set +a && ./bin/gateway

# Proto 代码生成
proto:
	sh script/proto/gen.sh

# Swagger 文档生成
swagger:
	$(HOME)/go/1.25.0/bin/swag init -g cmd/gateway/main.go -o docs

# 格式化 import
goimports:
	goimports -w .

# 代码检查
vet:
	go vet ./...

# 运行测试
test:
	go test -race -cover $$(go list ./... | grep -v -E 'cmd|docs|pkg|third_party|wire|interfaces|infrastructure|gateway|mock|repository|entity|valueobject') -short

# 检查项目（goimports + vet + test）
check: goimports vet test

# Docker 构建
docker-build:
	docker build -t vicomova/user -f docker/Dockerfile.user .
	docker build -t vicomova/video -f docker/Dockerfile.video .
	docker build -t vicomova/interaction -f docker/Dockerfile.interaction .
	docker build -t vicomova/chat -f docker/Dockerfile.chat .
	docker build -t vicomova/commerce -f docker/Dockerfile.commerce .
	docker build -t vicomova/gateway -f docker/Dockerfile.gateway .

# Docker 启动
docker-up:
	cd docker && docker-compose up -d

docker-up-prod:
	cd docker && docker-compose -f docker-compose.prod.yaml up -d

docker-down:
	cd docker && docker-compose down

docker-clean:
	cd docker && docker-compose down -v
	rm -rf docker/data/*
