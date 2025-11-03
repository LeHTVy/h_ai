# H-AI - HexStrike AI Clone in Go

🚀 **Advanced AI-Powered Penetration Testing Framework** - Golang Implementation

H-AI được viết bằng Golang, cung cấp nền tảng tự động hóa bảo mật mạnh mẽ với hơn 150 công cụ bảo mật và khả năng tích hợp AI agents.

## ✨ Tính năng

- **150+ Security Tools** - Tích hợp các công cụ bảo mật hàng đầu
- **REST API Server** - HTTP API server với Gin framework
- **MCP Support** - Hỗ trợ Model Context Protocol cho AI agents
- **Process Management** - Quản lý và theo dõi processes
- **Intelligent Caching** - Hệ thống cache thông minh
- **Error Recovery** - Tự động xử lý và phục hồi lỗi
- **High Performance** - Được tối ưu hóa cho hiệu suất với Golang

## 🏗️ Kiến trúc

```
H-AI/
├── cmd/
│   ├── server/          # HTTP API Server
│   └── mcp/             # MCP Server Client
├── internal/
│   ├── server/          # HTTP server handlers
│   ├── executor/        # Command execution engine
│   ├── cache/           # Caching system
│   ├── tools/            # Security tools manager
│   ├── client/           # API client
│   └── models/           # Data models
└── main.go              # Server entry point
```

## 📦 Cài đặt

### Yêu cầu

- Go 1.21 hoặc cao hơn
- Các security tools cần thiết (nmap, metasploit, v.v.)

### Build

```bash
# Clone repository
git clone https://github.com/LeHTVy/h_ai.git
cd h_ai

# Build server
go build -o bin/h-ai-server ./main.go

# Build MCP client
go build -o bin/h-ai-mcp ./cmd/mcp
```

### Chạy Server

```bash
# Chạy với mặc định (port 8888)
./bin/h-ai-server

# Chạy với tùy chọn
./bin/h-ai-server --port 9999 --host 0.0.0.0 --debug
```

### Chạy MCP Client

```bash
# Kết nối đến server
./bin/h-ai-mcp --server http://127.0.0.1:8888
```

## 🔧 Cấu hình

### Environment Variables

```bash
export H_AI_PORT=8888
export H_AI_HOST=0.0.0.0
export H_AI_DEBUG=true
```

## 📡 API Endpoints

### Health Check

```bash
GET /health
```

### Command Execution

```bash
POST /api/command
{
  "command": "nmap -sV target.com",
  "use_cache": true
}
```

### Security Tools

```bash
# Nmap scan
POST /api/tools/nmap
{
  "target": "target.com",
  "scan_type": "-sV",
  "ports": "80,443",
  "additional_args": "-T4"
}

# Metasploit
POST /api/tools/metasploit
{
  "module": "exploit/windows/smb/ms17_010_eternalblue",
  "options": {
    "RHOSTS": "192.168.1.1"
  }
}

# Gobuster
POST /api/tools/gobuster
{
  "url": "https://target.com",
  "mode": "dir",
  "wordlist": "/usr/share/wordlists/dirb/common.txt"
}
```

### Process Management

```bash
# List processes
GET /api/processes/list

# Process status
GET /api/processes/status/:pid

# Terminate process
POST /api/processes/terminate/:pid

# Dashboard
GET /api/processes/dashboard
```

### Intelligence

```bash
# Analyze target
POST /api/intelligence/analyze-target
{
  "target": "target.com",
  "analysis_type": "comprehensive"
}

# Select tools
POST /api/intelligence/select-tools
{
  "target": "target.com",
  "target_type": "web_application"
}

# Optimize parameters
POST /api/intelligence/optimize-parameters
{
  "tool": "nmap",
  "parameters": {
    "target": "target.com"
  }
}
```

## 🛠️ Security Tools

H-AI hỗ trợ các công cụ bảo mật sau:

### Network Scanning
- Nmap, Masscan, Rustscan
- AutoRecon, Amass, Subfinder
- NetExec, Enum4linux-ng

### Web Application
- Gobuster, Feroxbuster, FFuf
- Nuclei, Nikto, SQLMap
- WPScan, Arjun, ParamSpider

### Password Cracking
- Hydra, John the Ripper
- Hashcat, Medusa

### Exploitation
- Metasploit Framework
- MSFVenom

### Cloud Security
- Prowler, Scout Suite
- Trivy, Kube-Hunter

## 🤖 AI Agents Integration

### Claude Desktop

Chỉnh sửa `~/.config/Claude/claude_desktop_config.json`:

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

### Cursor / VS Code

Cấu hình trong settings:

```json
{
  "servers": {
    "h-ai": {
      "type": "stdio",
      "command": "/path/to/bin/h-ai-mcp",
      "args": ["--server", "http://localhost:8888"]
    }
  }
}
```

## 🔒 Bảo mật

⚠️ **Lưu ý quan trọng**:
- Công cụ này cung cấp quyền truy cập hệ thống mạnh mẽ cho AI agents
- Chạy trong môi trường cô lập hoặc VM chuyên dụng
- AI agents có thể thực thi các công cụ bảo mật tùy ý
- Giám sát hoạt động AI agents qua dashboard
- Xem xét triển khai xác thực cho môi trường production

## ⚖️ Sử dụng hợp pháp

✅ **Được phép**:
- Penetration testing có giấy phép
- Bug bounty programs
- CTF competitions
- Security research trên hệ thống sở hữu
- Red team exercises có phê duyệt

❌ **KHÔNG được phép**:
- Testing không có giấy phép
- Hoạt động độc hại
- Đánh cắp dữ liệu

## 📝 Giấy phép

MIT License

## 👨‍💻 Tác giả

Dựa trên HexStrike AI v6.0 - Ported to Golang

## 🙏 Đóng góp

Mọi đóng góp đều được chào đón! Vui lòng tạo issue hoặc pull request.

---

**Made with ❤️ - Golang Implementation of HexStrike AI**
