# 后端 Docker 镜像构建文件

# --- 阶段1：编译 ---
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o server main.go

# --- 阶段2：运行 ---
FROM alpine:latest

RUN apk --no-cache add ca-certificates tzdata

ENV TZ=Asia/Shanghai

WORKDIR /app

COPY --from=builder /app/server .
COPY --from=builder /app/config.yaml .

RUN mkdir -p uploads

EXPOSE 8080

CMD ["./server"]
