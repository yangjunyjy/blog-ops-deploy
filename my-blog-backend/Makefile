.PHONY: run build test clean docker

# 变量
APP_NAME=my-blog-backend
CMD_PATH=./cmd

# 运行
run:
	go run $(CMD_PATH)/main.go

# 构建
build:
	go build -o bin/$(APP_NAME) $(CMD_PATH)/main.go

# Wire依赖注入
wire:
	cd $(CMD_PATH) && wire

# 测试
test:
	go test -v ./...

# 清理
clean:
	rm -rf bin/

# Docker
docker-build:
	docker build -f deployments/Dockerfile -t $(APP_NAME):latest .

docker-up:
	docker-compose -f deployments/docker-compose.yaml up -d

docker-down:
	docker-compose -f deployments/docker-compose.yaml down
