# Hướng dẫn Test H-AI trên Ubuntu

## 🐧 Bước 1: Clone Repository trên Ubuntu

```bash
# Clone repository
git clone https://github.com/LeHTVy/h_ai.git
cd h_ai

# Kiểm tra Go version (cần Go 1.21+)
go version
```

## 📦 Bước 2: Cài đặt Dependencies

```bash
# Cài đặt Go nếu chưa có
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Kiểm tra lại
go version

# Download Go dependencies
cd h_ai
go mod download
go mod tidy
```

## 🔧 Bước 3: Cài đặt Security Tools (Kali Linux hoặc Ubuntu)

### Option A: Trên Kali Linux (Recommended - có sẵn nhiều tools)

```bash
# Kali Linux đã có sẵn nhiều tools
# Kiểm tra tools
which nmap gobuster nuclei sqlmap metasploit

# Nếu thiếu, cài thêm:
sudo apt update
sudo apt install -y nmap gobuster nuclei sqlmap metasploit-framework
```

### Option B: Trên Ubuntu thông thường

```bash
# Cài đặt basic tools
sudo apt update
sudo apt install -y golang-go

# Cài đặt security tools từ repositories
sudo apt install -y nmap gobuster

# Cài đặt Nuclei
go install -v github.com/projectdiscovery/nuclei/v3/cmd/nuclei@latest

# Cài đặt SQLMap
sudo apt install -y sqlmap

# Cài đặt Metasploit (optional)
curl https://raw.githubusercontent.com/rapid7/metasploit-omnibus/master/config/templates/metasploit-framework-wrappers/msfupdate.erb | sudo bash

# Kiểm tra tools
nmap --version
gobuster version
nuclei -version
sqlmap --version
```

## 🏗️ Bước 4: Build trên Ubuntu

```bash
cd h_ai

# Build server
go build -o bin/h-ai-server ./main.go

# Build MCP client
go build -o bin/h-ai-mcp ./cmd/mcp

# Kiểm tra binaries đã được tạo
ls -lh bin/
```

## 🚀 Bước 5: Chạy Server

```bash
# Chạy server
./bin/h-ai-server --port 8888 --debug

# Server sẽ chạy trên http://0.0.0.0:8888
# Kiểm tra logs để xem tools nào đã được detect
```

## 🧪 Bước 6: Test các Chức năng

### Test 1: Health Check

```bash
curl http://localhost:8888/health | jq
```

Expected output:
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "tools_status": {
    "nmap": true,
    "gobuster": true,
    ...
  }
}
```

### Test 2: Analyze Target

```bash
curl -X POST http://localhost:8888/api/intelligence/analyze-target \
  -H "Content-Type: application/json" \
  -d '{"target": "example.com", "analysis_type": "comprehensive"}' | jq
```

### Test 3: Select Tools

```bash
curl -X POST http://localhost:8888/api/intelligence/select-tools \
  -H "Content-Type: application/json" \
  -d '{"target": "example.com", "target_type": "web_application"}' | jq
```

### Test 4: Nmap Scan (cần có nmap)

```bash
# Scan một target công khai (scanme.nmap.org là target được phép scan)
curl -X POST http://localhost:8888/api/tools/nmap \
  -H "Content-Type: application/json" \
  -d '{
    "target": "scanme.nmap.org",
    "scan_type": "-sV",
    "ports": "22,80,443",
    "additional_args": "-T4"
  }' | jq
```

### Test 5: Gobuster Scan

```bash
curl -X POST http://localhost:8888/api/tools/gobuster \
  -H "Content-Type: application/json" \
  -d '{
    "url": "https://example.com",
    "mode": "dir",
    "wordlist": "/usr/share/wordlists/dirb/common.txt"
  }' | jq
```

### Test 6: Nuclei Scan

```bash
curl -X POST http://localhost:8888/api/tools/nuclei \
  -H "Content-Type: application/json" \
  -d '{
    "target": "https://example.com",
    "severity": "critical,high"
  }' | jq
```

### Test 7: Create Attack Chain

```bash
curl -X POST http://localhost:8888/api/intelligence/create-attack-chain \
  -H "Content-Type: application/json" \
  -d '{"target": "example.com", "analysis_type": "comprehensive"}' | jq
```

### Test 8: Smart Scan

```bash
curl -X POST http://localhost:8888/api/intelligence/smart-scan \
  -H "Content-Type: application/json" \
  -d '{
    "target": "example.com",
    "analysis_type": "comprehensive",
    "max_tools": 5
  }' | jq
```

## 📊 Bước 7: Monitor Processes

```bash
# List processes
curl http://localhost:8888/api/processes/list | jq

# Get dashboard
curl http://localhost:8888/api/processes/dashboard | jq

# Cache stats
curl http://localhost:8888/api/cache/stats | jq
```

## 🤖 Bước 8: Test MCP Client

### Terminal 1: Chạy API Server
```bash
./bin/h-ai-server --port 8888
```

### Terminal 2: Chạy MCP Client
```bash
./bin/h-ai-mcp --server http://127.0.0.1:8888
```

MCP client sẽ chạy và listen trên stdio, sẵn sàng nhận requests từ AI agents.

## 🔍 Test Script Tự động

Tạo file `test_api.sh`:

```bash
#!/bin/bash

BASE_URL="http://localhost:8888"

echo "🧪 Testing H-AI API..."
echo ""

echo "1. Health Check..."
curl -s "$BASE_URL/health" | jq .
echo ""

echo "2. Analyze Target..."
curl -s -X POST "$BASE_URL/api/intelligence/analyze-target" \
  -H "Content-Type: application/json" \
  -d '{"target": "example.com"}' | jq .
echo ""

echo "3. Select Tools..."
curl -s -X POST "$BASE_URL/api/intelligence/select-tools" \
  -H "Content-Type: application/json" \
  -d '{"target": "example.com"}' | jq .
echo ""

echo "✅ Tests completed!"
```

Chạy:
```bash
chmod +x test_api.sh
./test_api.sh
```

## 📝 Expected Results

Sau khi chạy trên Ubuntu với security tools đã cài:

1. **Health check**: Tất cả tools sẽ show `true` nếu đã cài đặt
2. **Nmap scan**: Sẽ có kết quả scan thực tế từ scanme.nmap.org
3. **Gobuster/Nuclei**: Sẽ chạy và trả về kết quả
4. **Process management**: Có thể list và terminate processes
5. **Intelligence endpoints**: Sẽ trả về analysis và tool recommendations

## ⚠️ Lưu ý

- Đảm bảo có quyền chạy các security tools (một số cần sudo cho network operations)
- Kiểm tra firewall nếu test từ máy khác
- Sử dụng targets hợp pháp để test (ví dụ: scanme.nmap.org, example.com)
- Một số tools như metasploit cần setup database riêng

## 🐛 Troubleshooting

### Lỗi: "permission denied" khi chạy nmap
```bash
# Một số scan types cần quyền root
sudo ./bin/h-ai-server --port 8888
```

### Lỗi: "tool not found"
```bash
# Kiểm tra PATH
echo $PATH
which nmap

# Thêm Go bin vào PATH nếu cần
export PATH=$PATH:$(go env GOPATH)/bin
```

### Lỗi: "port already in use"
```bash
# Tìm process đang dùng port 8888
lsof -i :8888

# Kill process
kill -9 <PID>

# Hoặc dùng port khác
./bin/h-ai-server --port 9999
```
