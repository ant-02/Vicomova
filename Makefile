.PHONY: dev dev-stop dev-logs run-user run-gateway docker-build docker-up docker-up-prod docker-down docker-clean proto swagger

CONFIG_FILE = config/base.yaml
SESSION_NAME = vicomova

# tmux 开发模式：2 面板，左边最近访问窗口，右边主终端
dev:
	@if tmux has-session -t $(SESSION_NAME) 2>/dev/null; then \
		tmux kill-session -t $(SESSION_NAME); \
	fi
	@echo "Starting Vicomova services..."
	@tmux new-session -d -s $(SESSION_NAME) -n main
	@tmux split-window -h -t $(SESSION_NAME):main
	@tmux new-window -t $(SESSION_NAME) -n user
	@tmux new-window -t $(SESSION_NAME) -n gateway
	@tmux send-keys -t $(SESSION_NAME):main.0 "echo 'User window: Ctrl+B , then 1'" C-m
	@tmux send-keys -t $(SESSION_NAME):main.1 "echo '=== TERMINAL ===' && echo 'Main operational window'" C-m
	@tmux send-keys -t $(SESSION_NAME):user "make run-user" C-m
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
	@echo "=== GATEWAY LOGS ===" && tmux capture-pane -t $(SESSION_NAME):gateway -p | tail -30

# 直接运行
run-user:
	go run ./cmd/user -config $(CONFIG_FILE) -port 8888

run-gateway:
	go run ./cmd/gateway -config $(CONFIG_FILE) -port 8080

# Proto 代码生成
proto:
	sh script/proto/gen.sh

# Swagger 文档生成
swagger:
	$(HOME)/go/1.24.10/bin/swag init -g cmd/gateway/main.go -o docs

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
