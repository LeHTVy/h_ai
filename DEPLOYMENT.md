# Hướng dẫn Deploy và Test H-AI

## 📦 Bước 1: Chuẩn bị và Push lên GitHub

### 1.1. Kiểm tra Git status
```bash
cd h-ai
git status
```

### 1.2. Add tất cả files
```bash
git add .
```

### 1.3. Commit
```bash
git commit -m "Initial commit: H-AI - HexStrike AI Clone in Go

- REST API Server với Gin framework
- MCP Server implementation
- AI Decision Engine và Agents
- 12+ Security Tools integration
- Process Management và Caching
- Error Recovery System"
```

### 1.4. Push lên GitHub
```bash
git push -u origin main
```

## 🚀 Bước 2: Build và Chạy Locally

### 2.1. Cài đặt dependencies
```bash
cd h-ai
go mod download
go mod tidy
```

### 2.2. Build binaries
```bash
# Build server
go build -o bin/h-ai-server ./main.go

# Build MCP client
go build -o bin/h-ai-mcp ./cmd/mcp
```

Hoặc sử dụng Makefile:
```bash
make build
```

### 2.3. Chạy Server
```bash
# Chạy với default settings (port 8888)
./bin/h-ai-server

# Hoặc với options
./bin/h-ai-server --port 8888 --host 0.0.0.0 --debug
```

## 🧪 Bước 3: Test API Server

### 3.1. Test Health Check
```bash
curl http://localhost:8888/health
```

Expected response:
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "tools_status": {...}
}
```

### 3.2. Test Target Analysis
```bash
curl -X POST http://localhost:8888/api/intelligence/analyze-target \
  -H "Content-Type: application/json" \
  -d '{"target": "example.com", "analysis_type": "comprehensive"}'
```

### 3.3. Test Tool Selection
```bash
curl -X POST http://localhost:8888/api/intelligence/select-tools \
  -H "Content-Type: application/json" \
  -d '{"target": "example.com", "target_type": "web_application"}'
```

### 3.4. Test Nmap Scan
```bash
curl -X POST http://localhost:8888/api/tools/nmap \
  -H "Content-Type: application/json" \
  -d '{
    "target": "scanme.nmap.org",
    "scan_type": "-sV",
    "ports": "80,443",
    "additional_args": "-T4"
  }'
```

### 3.5. Test Gobuster Scan
```bash
curl -X POST http://localhost:8888/api/tools/gobuster \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "mode": "dir",
    "wordlist": "/usr/share/wordlists/dirb/common.txt"
  }'
```

## 🤖 Bước 4: Test MCP Client

### 4.1. Chạy MCP Client
```bash
./bin/h-ai-mcp --server http://127.0.0.1:8888
```

### 4.2. Cấu hình Claude Desktop

Chỉnh sửa `~/.config/Claude/claude_desktop_config.json` (Linux/Mac) hoặc `%APPDATA%\Claude\claude_desktop_config.json` (Windows):

```json
{
  "mcpServers": {
    "h-ai": {
      "command": "/path/to/bin/h-ai-mcp",
      "args": ["--server", "http://localhost:8888"],
      "description": "H-AI - Advanced Cybersecurity Automation Platform",
      "timeout": 300
    }
  }
}
```

### 4.3. Test trong Claude Desktop
- Mở Claude Desktop
- Hỏi Claude: "Use h-ai to run an nmap scan on scanme.nmap.org"
- Claude sẽ tự động gọi MCP tools

## 📋 Bước 5: Test với Example Client

### 5.1. Chạy example client
```bash
cd examples
go run example_client.go
```

## 🔍 Bước 6: Kiểm tra Logs

### 6.1. Xem server logs
Server sẽ log ra console khi chạy với `--debug`

### 6.2. Check processes
```bash
curl http://localhost:8888/api/processes/list
```

### 6.3. Check cache stats
```bash
curl http://localhost:8888/api/cache/stats
```

## 🐳 Bước 7: Deploy với Docker (Optional)

### 7.1. Build Docker image
```bash
docker build -t h-ai:latest .
```

### 7.2. Run Docker container
```bash
docker run -d -p 8888:8888 --name h-ai h-ai:latest
```

### 7.3. Check logs
```bash
docker logs h-ai
```

## ✅ Checklist trước khi Push

- [ ] `go mod tidy` đã chạy thành công
- [ ] Không có lỗi compile
- [ ] `.gitignore` đã được cấu hình đúng
- [ ] `README.md` đã được cập nhật
- [ ] Module path trong `go.mod` đúng với GitHub repo
- [ ] Tất cả imports đã được cập nhật

## 🔧 Troubleshooting

### Lỗi: "module path mismatch"
```bash
# Kiểm tra module path trong go.mod
# Đảm bảo nó khớp với GitHub URL
# Sau đó chạy:
go mod tidy
```

### Lỗi: "tool not found"
- Đảm bảo security tools đã được cài đặt (nmap, gobuster, etc.)
- Kiểm tra PATH environment variable

### Lỗi: "port already in use"
```bash
# Thay đổi port
./bin/h-ai-server --port 9999
```

### Lỗi: "connection refused" (MCP client)
- Đảm bảo API server đang chạy
- Kiểm tra URL trong MCP config

## 📚 Tài liệu tham khảo

- README.md - Hướng dẫn sử dụng chi tiết
- PROJECT_STRUCTURE.md - Cấu trúc dự án
- API endpoints được document trong code comments
