# Quick Start Guide - H-AI

## 🚀 Push lên GitHub và Test

### Bước 1: Khởi tạo Git Repository

```bash
cd h-ai

# Khởi tạo git (nếu chưa có)
git init

# Add remote (đã có từ trước)
git remote add origin https://github.com/LeHTVy/h_ai.git

# Hoặc nếu đã có, set lại:
git remote set-url origin https://github.com/LeHTVy/h_ai.git
```

### Bước 2: Kiểm tra và Cleanup

```bash
# Đảm bảo go.mod đúng
go mod tidy

# Kiểm tra không có lỗi compile
go build ./...
```

### Bước 3: Add và Commit

```bash
# Xem files sẽ được add
git status

# Add tất cả files
git add .

# Commit
git commit -m "feat: Initial release of H-AI

- Complete Go implementation of HexStrike AI
- REST API Server with Gin framework  
- MCP Server for AI agent integration
- AI Decision Engine and specialized agents
- 12+ Security tools integration (nmap, metasploit, gobuster, etc.)
- Process management and intelligent caching
- Error recovery system
- Full documentation and examples"
```

### Bước 4: Push lên GitHub

```bash
# Set branch
git branch -M main

# Push
git push -u origin main
```

## 🧪 Test Local

### Test 1: Build và Chạy Server

```bash
# Build
go build -o bin/h-ai-server ./main.go
go build -o bin/h-ai-mcp ./cmd/mcp

# Chạy server
./bin/h-ai-server --port 8888 --debug

# Trong terminal khác, test health
curl http://localhost:8888/health
```

### Test 2: Test API Endpoints

```bash
# Test analyze target
curl -X POST http://localhost:8888/api/intelligence/analyze-target \
  -H "Content-Type: application/json" \
  -d '{"target": "example.com"}'

# Test select tools
curl -X POST http://localhost:8888/api/intelligence/select-tools \
  -H "Content-Type: application/json" \
  -d '{"target": "example.com", "target_type": "web_application"}'

# Test nmap (cần có nmap installed)
curl -X POST http://localhost:8888/api/tools/nmap \
  -H "Content-Type: application/json" \
  -d '{
    "target": "scanme.nmap.org",
    "scan_type": "-sV",
    "ports": "80,443"
  }'
```

### Test 3: Test MCP Client

```bash
# Terminal 1: Chạy server
./bin/h-ai-server

# Terminal 2: Chạy MCP client
./bin/h-ai-mcp --server http://127.0.0.1:8888
```

## ✅ Verification Checklist

Sau khi push, kiểm tra:

1. **GitHub Repository**
   - [ ] Tất cả files đã được push
   - [ ] README.md hiển thị đúng
   - [ ] .gitignore hoạt động đúng

2. **Code Quality**
   - [ ] `go mod tidy` không có lỗi
   - [ ] Build thành công: `go build ./...`
   - [ ] Không có linter errors

3. **Functionality**
   - [ ] Server start được: `./bin/h-ai-server`
   - [ ] Health endpoint trả về đúng: `curl http://localhost:8888/health`
   - [ ] API endpoints hoạt động

## 🔧 Nếu có lỗi

### Lỗi: "module path mismatch"
```bash
# Kiểm tra go.mod
cat go.mod

# Đảm bảo module path là: github.com/LeHTVy/h_ai
# Nếu sai, sửa và chạy:
go mod tidy
```

### Lỗi: "cannot find package"
```bash
# Download dependencies
go mod download
go mod tidy
```

### Lỗi khi build
```bash
# Clean và build lại
go clean -cache
go mod tidy
go build ./...
```
