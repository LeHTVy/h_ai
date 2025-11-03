# Hướng dẫn Sử Dụng H-AI

## 🤔 H-AI là gì?

H-AI là **cả hai**:
- ✅ **HTTP API Server độc lập** - Dùng trực tiếp qua REST API
- ✅ **MCP Server** - Bridge để AI agents (Claude, GPT, Cursor) gọi tools

## 📊 So Sánh 2 Cách Dùng

| Tính năng | HTTP API (Độc lập) | MCP (Với AI) |
|-----------|-------------------|--------------|
| **Cần AI?** | ❌ Không | ✅ Có (Claude/GPT/Cursor) |
| **Cách gọi** | HTTP requests | Qua AI agent |
| **Tự động?** | Manual/script | AI tự quyết định |
| **Dễ dùng?** | Cần biết API | Chỉ cần chat với AI |

## 🎯 Cách 1: Dùng Trực Tiếp (Không Cần AI)

### Setup

```bash
# Chạy HTTP API Server
./bin/h-ai-server --port 8888
```

### Sử dụng

```bash
# Gọi trực tiếp qua curl
curl -X POST http://localhost:8888/api/tools/nmap \
  -H "Content-Type: application/json" \
  -d '{"target": "scanme.nmap.org", "scan_type": "-sV"}'

# Hoặc dùng trong script Python
import requests
response = requests.post('http://localhost:8888/api/tools/nmap', 
  json={"target": "example.com"})
print(response.json())
```

**Ưu điểm:**
- ✅ Độc lập, không cần AI
- ✅ Tích hợp dễ vào automation/scripts
- ✅ Control hoàn toàn

**Nhược điểm:**
- ❌ Cần biết API endpoints
- ❌ Phải tự quyết định tool nào dùng

## 🤖 Cách 2: Dùng Qua AI Agents (Với AI)

### Setup

```bash
# Terminal 1: Chạy HTTP API Server
./bin/h-ai-server --port 8888

# Terminal 2: Chạy MCP Server (để AI agents connect)
./bin/h-ai-mcp --server http://127.0.0.1:8888
```

### Cấu hình Claude Desktop

Edit `~/.config/Claude/claude_desktop_config.json`:
```json
{
  "mcpServers": {
    "h-ai": {
      "command": "/path/to/bin/h-ai-mcp",
      "args": ["--server", "http://localhost:8888"]
    }
  }
}
```

### Sử dụng

Mở Claude Desktop và chat:
```
"I'm a security researcher. Can you scan scanme.nmap.org using h-ai tools?"
```

Claude sẽ:
1. Hiểu bạn muốn scan
2. Tự động gọi `nmap_scan` tool từ H-AI
3. Hiển thị kết quả cho bạn

**Ưu điểm:**
- ✅ Dễ dùng - chỉ cần chat
- ✅ AI tự quyết định tool nào phù hợp
- ✅ Natural language interface

**Nhược điểm:**
- ❌ Cần AI agent (Claude Desktop, Cursor, etc.)
- ❌ Phụ thuộc vào AI để hiểu intent

## 🔄 Workflow Chi Tiết Với AI

### Scenario: User muốn scan một website

```
1. User: "I need to scan example.com for vulnerabilities"

2. Claude (AI Agent):
   - Phân tích: User muốn vulnerability scan
   - Quyết định: Cần dùng nhiều tools
   - Gọi tools qua MCP:
     * nmap_scan (để tìm open ports)
     * nuclei_scan (để tìm vulnerabilities)
     * gobuster_scan (để tìm hidden directories)

3. MCP Server (h-ai-mcp):
   - Nhận requests từ Claude
   - Convert thành HTTP API calls
   - Gửi đến h-ai-server

4. HTTP API Server (h-ai-server):
   - Execute nmap, nuclei, gobuster
   - Trả kết quả về MCP Server

5. MCP Server:
   - Format results
   - Trả về cho Claude qua JSON-RPC

6. Claude:
   - Phân tích results
   - Tạo report dễ đọc
   - Hiển thị cho User
```

## 🛠️ Khi Nào Dùng Cách Nào?

### Dùng HTTP API trực tiếp khi:
- ✅ Tích hợp vào automation/CI/CD
- ✅ Viết script tự động
- ✅ Cần control chính xác từng bước
- ✅ Không muốn phụ thuộc AI

### Dùng MCP với AI khi:
- ✅ Muốn giao tiếp tự nhiên (chat)
- ✅ Cần AI tự quyết định tool chain
- ✅ Muốn AI phân tích và tổng hợp kết quả
- ✅ Làm bug bounty hoặc CTF với AI assistance

## 📝 Examples

### Example 1: Script tự động scan

```bash
#!/bin/bash
# scan.sh - Automated scan script

SERVER="http://localhost:8888"
TARGET="$1"

# Step 1: Analyze target
curl -X POST "$SERVER/api/intelligence/analyze-target" \
  -H "Content-Type: application/json" \
  -d "{\"target\": \"$TARGET\"}" > analysis.json

# Step 2: Get recommended tools
TOOLS=$(curl -s -X POST "$SERVER/api/intelligence/select-tools" \
  -H "Content-Type: application/json" \
  -d "{\"target\": \"$TARGET\"}" | jq -r '.selected_tools[]')

# Step 3: Run each tool
for tool in $TOOLS; do
  echo "Running $tool..."
  curl -X POST "$SERVER/api/tools/$tool" \
    -H "Content-Type: application/json" \
    -d "{\"target\": \"$TARGET\"}" > "${tool}_results.json"
done
```

### Example 2: Chat với Claude

```
User: "I'm testing my website example.com. Can you help me run a comprehensive security scan?"

Claude: "I'll help you scan example.com. Let me start with:
1. Analyzing the target structure
2. Running Nmap to find open ports
3. Scanning for vulnerabilities with Nuclei
4. Checking for hidden directories with Gobuster

[Claude tự động gọi các h-ai tools và hiển thị kết quả]"

User: "What vulnerabilities did you find?"

Claude: [Phân tích kết quả và tóm tắt vulnerabilities]
```

## ❓ FAQ

### Q: Tôi có bắt buộc phải dùng AI không?

**A:** Không! HTTP API Server hoạt động độc lập. Bạn có thể dùng trực tiếp qua curl, Postman, hoặc script.

### Q: AI Agent là gì và ở đâu?

**A:** AI Agent là Claude Desktop, Cursor, hoặc các MCP-compatible clients. Chúng đã có sẵn AI (Claude, GPT). H-AI chỉ cung cấp tools cho các AI này sử dụng.

### Q: MCP Server có bắt buộc không?

**A:** Chỉ cần khi muốn dùng với AI agents. Nếu chỉ dùng HTTP API trực tiếp thì không cần MCP Server.

### Q: Làm sao AI biết khi nào dùng tool nào?

**A:** AI agent (Claude/GPT) sẽ phân tích prompt của bạn và tự quyết định:
- "scan website" → gọi nmap_scan, nuclei_scan
- "find directories" → gọi gobuster_scan
- "test SQL injection" → gọi sqlmap_scan

AI sử dụng intelligence có sẵn của nó + tools từ H-AI.

### Q: Có thể dùng cả 2 cách cùng lúc không?

**A:** Có! Bạn có thể:
- Chạy HTTP API Server
- Một số requests đi trực tiếp qua HTTP API
- Một số requests đi qua AI + MCP
