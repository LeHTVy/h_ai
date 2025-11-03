# Kiến trúc H-AI

## 📐 Tổng quan

H-AI có **2 components** có thể hoạt động độc lập hoặc kết hợp:

1. **HTTP API Server** (`h-ai-server`) - Chương trình độc lập
2. **MCP Server** (`h-ai-mcp`) - Kết nối AI agents với HTTP API Server

## 🏗️ Kiến trúc 2-Tier

```
┌─────────────────────────────────────────────────────────────┐
│                    USER / AI AGENT                           │
└──────────────────────────┬──────────────────────────────────┘
                           │
           ┌───────────────┴───────────────┐
           │                               │
           ▼                               ▼
    ┌─────────────┐              ┌─────────────────┐
    │  Option 1  │              │    Option 2     │
    │ Direct API │              │  Via AI Agent   │
    └──────┬──────┘              └────────┬────────┘
           │                               │
           │                               │
           ▼                               ▼
    ┌──────────────────────────────────────────────┐
    │      HTTP API Server (h-ai-server)          │
    │      Port: 8888                              │
    │      - REST API Endpoints                    │
    │      - Security Tools Execution              │
    │      - AI Decision Engine                    │
    └───────────────────┬──────────────────────────┘
                        │
                        │
        ┌───────────────┴───────────────┐
        │                               │
        ▼                               ▼
┌──────────────┐              ┌──────────────────┐
│ Security     │              │ Process Manager  │
│ Tools        │              │ & Cache          │
│ (nmap, etc.) │              └──────────────────┘
└──────────────┘
```

## 🔄 Hai Cách Sử Dụng

### Option 1: Sử Dụng Trực Tiếp (Không Cần AI)

**HTTP API Server** là chương trình **độc lập**, có thể dùng trực tiếp:

```bash
# 1. Chạy server
./bin/h-ai-server --port 8888

# 2. Gọi API trực tiếp bằng curl/Postman/browser
curl -X POST http://localhost:8888/api/tools/nmap \
  -H "Content-Type: application/json" \
  -d '{"target": "scanme.nmap.org", "scan_type": "-sV"}'
```

**Workflow:**
```
User → HTTP Request → H-AI Server → Security Tools → Response → User
```

### Option 2: Sử Dụng Qua AI Agents (Với AI)

**MCP Server** là bridge giữa AI agents và HTTP API Server:

```bash
# Terminal 1: Chạy HTTP API Server
./bin/h-ai-server --port 8888

# Terminal 2: Chạy MCP Server (optional - chỉ khi dùng với AI)
./bin/h-ai-mcp --server http://127.0.0.1:8888
```

**Workflow với AI:**
```
User Prompt → AI Agent (Claude/GPT) → MCP Protocol → MCP Server → HTTP API → Security Tools → Response → MCP → AI Agent → User
```

## 🤖 Luồng Hoạt Động với AI Agents

### Bước 1: User nhập prompt vào AI

Ví dụ trong Claude Desktop hoặc Cursor:
```
"I'm a security researcher. Can you help me run an nmap scan on scanme.nmap.org using h-ai tools?"
```

### Bước 2: AI Agent xử lý prompt

AI agent (Claude, GPT, Cursor) phân tích prompt và quyết định:
- Cần dùng tool nào (nmap_scan)
- Cần tham số gì (target: scanme.nmap.org)

### Bước 3: AI Agent gọi MCP Tools

AI agent gửi JSON-RPC request qua MCP protocol:
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "nmap_scan",
    "arguments": {
      "target": "scanme.nmap.org",
      "scan_type": "-sV"
    }
  }
}
```

### Bước 4: MCP Server nhận request

MCP Server (`h-ai-mcp`) nhận request từ stdin, parse và:
- Xác định tool cần gọi
- Convert MCP arguments thành HTTP API request
- Gửi request đến HTTP API Server

### Bước 5: HTTP API Server xử lý

HTTP API Server (`h-ai-server`):
- Nhận HTTP request từ MCP Server
- Gọi security tool (nmap)
- Trả về kết quả

### Bước 6: MCP Server trả kết quả cho AI

MCP Server format response và trả về cho AI agent qua stdout (JSON-RPC)

### Bước 7: AI Agent hiển thị cho User

AI agent nhận kết quả và format thành response dễ đọc cho user

## 📋 Setup với AI Agents

### Claude Desktop

1. **Chạy HTTP API Server:**
```bash
./bin/h-ai-server --port 8888
```

2. **Cấu hình Claude Desktop:**
Edit `~/.config/Claude/claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "h-ai": {
      "command": "/path/to/bin/h-ai-mcp",
      "args": ["--server", "http://localhost:8888"],
      "description": "H-AI Cybersecurity Tools",
      "timeout": 300
    }
  }
}
```

3. **Claude Desktop tự động:**
- Khi user chat với Claude, Claude sẽ tự động gọi `h-ai-mcp`
- MCP Server sẽ connect đến HTTP API Server
- Tools được execute và trả kết quả về Claude

### Cursor / VS Code

Tương tự, cấu hình trong settings và Cursor sẽ tự động gọi tools khi cần.

## 🔑 Điểm Quan Trọng

1. **HTTP API Server là độc lập**: 
   - Có thể dùng trực tiếp qua REST API
   - Không bắt buộc phải có AI agent
   - Có thể tích hợp vào script, web app, automation

2. **MCP Server là optional**:
   - Chỉ cần khi muốn dùng với AI agents
   - Là bridge giữa AI và HTTP API
   - Không cần AI, bạn vẫn dùng được HTTP API

3. **AI Agent (Claude/GPT) là phần của client**:
   - Không phải part của H-AI
   - Claude Desktop, Cursor đã có sẵn AI
   - H-AI chỉ cung cấp tools cho AI sử dụng

## 💡 Tóm Tắt

- **H-AI Server**: Chương trình độc lập, có thể dùng trực tiếp
- **H-AI MCP**: Optional component để AI agents có thể gọi tools
- **AI Agent**: Có sẵn trong Claude Desktop/Cursor, không phải part của H-AI
- **User**: Có thể dùng trực tiếp (HTTP API) hoặc qua AI (MCP)

## 🎯 Use Cases

### Use Case 1: Trực tiếp qua API
```bash
# Script tự động scan
curl -X POST http://localhost:8888/api/tools/nmap \
  -d '{"target": "target.com"}' > results.json
```

### Use Case 2: Qua AI Agent
```
User: "Scan example.com for open ports"
Claude: [Tự động gọi h-ai nmap_scan tool và hiển thị kết quả]
```

### Use Case 3: Hybrid
```bash
# AI agent xử lý planning
# HTTP API thực thi actual scans
# AI agent phân tích results
```
