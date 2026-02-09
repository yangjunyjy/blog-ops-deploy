# 多阶段构建
FROM golang:1.21-alpine AS builder

WORKDIR /app

# 安装依赖
COPY go.mod go.sum ./
RUN go mod download

# 复制源代码
COPY . .

# 构建应用
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# 运行阶段
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata
ENV TZ=Asia/Shanghai

WORKDIR /root/

COPY --from=builder /app/server .
COPY --from=builder /app/config ./config

EXPOSE 8080

CMD ["./server"]
