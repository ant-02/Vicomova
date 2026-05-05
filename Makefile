.PHONY: dev dev-stop dev-logs run-user run-gateway run-video run-interaction docker-build docker-up docker-up-prod docker-down docker-clean proto swagger lint test

CONFIG_FILE = config/base.yaml
SESSION_NAME = vicomova

# tmux 开发模式：5 个窗口
dev:
	@if tmux has-session -t $(SESSION_NAME) 2>/dev/null; then \
		tmux kill-session -t $(SESSION_NAME); \
	fi
	@echo "Starting Vicomova services..."
	@tmux new-session -d -s $(SESSION_NAME) -n main
	@tmux new-window -t $(SESSION_NAME) -n user
	@tmux new-window -t $(SESSION_NAME) -n video
	@tmux new-window -t $(SESSION_NAME) -n interaction
	@tmux new-window -t $(SESSION_NAME) -n gateway
	@tmux send-keys -t $(SESSION_NAME):main "echo 'Vicomova dev session'" C-m
	@tmux send-keys -t $(SESSION_NAME):user "make run-user" C-m
	@tmux send-keys -t $(SESSION_NAME):video "make run-video" C-m
	@tmux send-keys -t $(SESSION_NAME):interaction "make run-interaction" C-m
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
	@echo "=== GATEWAY LOGS ===" && tmux capture-pane -t $(SESSION_NAME):gateway -p | tail -30

# 编译输出到 bin/
build:
	mkdir -p bin
	go build -o bin/user ./cmd/user
	go build -o bin/video ./cmd/video
	go build -o bin/interaction ./cmd/interaction
	go build -o bin/gateway ./cmd/gateway

# 直接运行（编译后执行）
run-user: build
	./bin/user -config $(CONFIG_FILE) -port 8888

run-video: build
	./bin/video -config $(CONFIG_FILE) -port 8889

run-interaction: build
	./bin/interaction -config $(CONFIG_FILE) -port 8890

run-gateway: build
	./bin/gateway -config $(CONFIG_FILE) -port 8080

# Proto 代码生成
proto:
	sh script/proto/gen.sh

# Swagger 文档生成
swagger:
	$(HOME)/go/1.25.0/bin/swag init -g cmd/gateway/main.go -o docs

# Lint 代码
lint:
	golangci-lint run ./...

# 运行测试
test:
	go test ./... -short

# 检查项目（lint + test）
check: lint test

# Docker 构建
docker-build:
	docker build -t vicomova-user -f docker/Dockerfile.user .
	docker build -t vicomova-gateway -f docker/Dockerfile.gateway .

# Docker 启动
docker-up:
	cd docker && docker-compose up -d

docker-up-prod:
	cd docker && docker-compose -f docker-compose.prod.yaml up -d

docker-down:
	cd docker && docker-compose down

docker-clean:
	cd docker && docker-compose down -v
	rm -rf docker/mysql/* docker/redis/*
