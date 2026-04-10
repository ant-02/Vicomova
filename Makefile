.PHONY: run docker-build docker-up docker-up-prod docker-down docker-clean proto swagger

# 开发阶段运行
run-user:
	go run ./cmd/user -config config/base.yaml -port 8888

run-gateway:
	go run ./cmd/gateway -config config/base.yaml -port 8080 -user_addr 127.0.0.1:8888

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

# Docker 启动（仅基础设施，用于本地开发）
docker-up:
	cd docker && docker-compose up -d

# Docker 启动（完整部署，需要先 docker-build）
docker-up-prod:
	cd docker && docker-compose -f docker-compose.prod.yaml up -d

docker-down:
	cd docker && docker-compose down

docker-clean:
	cd docker && docker-compose down -v
	rm -rf docker/mysql/* docker/redis/*
