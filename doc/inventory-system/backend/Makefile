.PHONY: run build test test-cover fmt lint db-init db-seed clean

# 运行后端（开发模式）
run:
	go run main.go

# 编译为可执行文件
build:
	CGO_ENABLED=0 go build -o bin/server main.go

# 运行所有测试
test:
	go test ./internal/... -v

# 运行测试并生成覆盖率报告
test-cover:
	go test ./internal/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "覆盖率报告已生成: coverage.html"

# 代码格式化
fmt:
	gofmt -w .

# 初始化数据库
db-init:
	mysql -u root -p < scripts/init.sql
	@echo "数据库初始化完成"

# 插入测试数据
db-seed:
	mysql -u root -p < scripts/seed.sql
	@echo "测试数据插入完成"

# 清理编译产物
clean:
	rm -rf bin/ coverage.out coverage.html
