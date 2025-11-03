# 💬 Hướng dẫn Chat với Ollama LLM qua H-AI

Bạn có thể chat trực tiếp với LLM (Ollama) qua H-AI API để hỏi các câu hỏi phức tạp về cybersecurity, exploit development, tool usage, v.v.

## 🚀 Cách sử dụng

### 1. Khởi động server với Ollama

```bash
./bin/h-ai-server --ollama-model llama2 --port 8888
```

### 2. Chat đơn giản (một lần)

```bash
curl -X POST http://localhost:8888/api/intelligence/chat \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "user",
        "content": "Giải thích cách SQL injection hoạt động và cách phòng chống"
      }
    ]
  }'
```

### 3. Chat với nhiều tin nhắn (conversation)

```bash
curl -X POST http://localhost:8888/api/intelligence/chat \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "system",
        "content": "Bạn là một chuyên gia cybersecurity với 20 năm kinh nghiệm."
      },
      {
        "role": "user",
        "content": "Tôi cần exploit CVE-2023-1234, bạn có thể hướng dẫn không?"
      },
      {
        "role": "assistant",
        "content": "Để exploit CVE-2023-1234, bạn cần..."
      },
      {
        "role": "user",
        "content": "Vậy làm sao để tạo payload?"
      }
    ]
  }'
```

### 4. Điều chỉnh Temperature và TopP

```bash
curl -X POST http://localhost:8888/api/intelligence/chat \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "user",
        "content": "Viết một exploit Python cho buffer overflow"
      }
    ],
    "temperature": 0.3,
    "top_p": 0.9
  }'
```

**Temperature:**
- `0.0-0.3`: Có cấu trúc, chính xác (tốt cho code, exploit)
- `0.4-0.7`: Cân bằng (mặc định)
- `0.8-1.0`: Sáng tạo, đa dạng (tốt cho brainstorming)

**TopP:**
- `0.1-0.5`: Tập trung vào top tokens
- `0.6-0.9`: Đa dạng hơn (mặc định 0.9)
- `1.0`: Tất cả tokens

## 📝 Ví dụ sử dụng thực tế

### 1. Hỏi về Security Tools

```bash
curl -X POST http://localhost:8888/api/intelligence/chat \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "user",
        "content": "Làm sao để dùng Nmap để scan stealthy nhất có thể? Cho tôi command cụ thể."
      }
    ],
    "temperature": 0.2
  }'
```

### 2. Phân tích Vulnerability

```bash
curl -X POST http://localhost:8888/api/intelligence/chat \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "user",
        "content": "Tôi đã scan và thấy port 443 mở với Apache 2.4.41. Có những lỗ hổng nào tôi nên kiểm tra?"
      }
    ],
    "temperature": 0.3
  }'
```

### 3. Viết Exploit Code

```bash
curl -X POST http://localhost:8888/api/intelligence/chat \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "system",
        "content": "Bạn là một exploit developer chuyên nghiệp. Viết code Python rõ ràng, có comment."
      },
      {
        "role": "user",
        "content": "Viết exploit cho buffer overflow trên port 9999, có thể overflow 200 bytes."
      }
    ],
    "temperature": 0.2,
    "top_p": 0.8
  }'
```

### 4. Hỏi về Attack Chain

```bash
curl -X POST http://localhost:8888/api/intelligence/chat \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "user",
        "content": "Tôi đã có shell trên một web server Linux. Bây giờ tôi cần privilege escalation. Hãy đưa ra một attack chain chi tiết từng bước."
      }
    ]
  }'
```

### 5. Phân tích Kết quả Scan

```bash
curl -X POST http://localhost:8888/api/intelligence/chat \
  -H "Content-Type: application/json" \
  -d '{
    "messages": [
      {
        "role": "user",
        "content": "Tôi scan một target và thấy:\n- Port 22: OpenSSH 7.4\n- Port 80: Apache 2.4.6\n- Port 3306: MySQL 5.7\n- Port 8080: Tomcat 8.5\n\nHãy phân tích và đưa ra kế hoạch tấn công chi tiết."
      }
    ],
    "temperature": 0.4
  }'
```

## 🔧 Sử dụng với Python

```python
import requests
import json

def chat_with_hai(prompt, conversation_history=None, temperature=0.7):
    url = "http://localhost:8888/api/intelligence/chat"
    
    messages = conversation_history or []
    messages.append({
        "role": "user",
        "content": prompt
    })
    
    payload = {
        "messages": messages,
        "temperature": temperature
    }
    
    response = requests.post(url, json=payload)
    return response.json()

# Ví dụ sử dụng
result = chat_with_hai("Giải thích cách XSS hoạt động")
print(result["response"])

# Conversation
history = [
    {"role": "user", "content": "CVE-2023-1234 là gì?"},
    {"role": "assistant", "content": "CVE-2023-1234 là một lỗ hổng..."}
]
result = chat_with_hai("Làm sao để exploit?", history)
print(result["response"])
```

## 🔧 Sử dụng với JavaScript/Node.js

```javascript
async function chatWithHAI(prompt, conversationHistory = null, temperature = 0.7) {
  const url = 'http://localhost:8888/api/intelligence/chat';
  
  const messages = conversationHistory || [];
  messages.push({
    role: 'user',
    content: prompt
  });
  
  const response = await fetch(url, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({
      messages: messages,
      temperature: temperature
    })
  });
  
  return await response.json();
}

// Sử dụng
const result = await chatWithHAI('Giải thích cách SQL injection hoạt động');
console.log(result.response);
```

## 🎯 Best Practices

1. **System Prompt**: Sử dụng `role: "system"` để set context cho LLM
   ```json
   {
     "role": "system",
     "content": "Bạn là một penetration tester chuyên nghiệp..."
   }
   ```

2. **Temperature thấp cho code**: Khi cần code chính xác, dùng `temperature: 0.2-0.3`

3. **Conversation History**: Giữ lại lịch sử để LLM hiểu context

4. **Chia nhỏ câu hỏi phức tạp**: Thay vì một câu hỏi dài, chia thành nhiều câu ngắn hơn

5. **Kiểm tra response**: Luôn kiểm tra `success: true` trước khi dùng response

## 📊 Response Format

```json
{
  "success": true,
  "response": "Câu trả lời từ LLM...",
  "message_count": 1
}
```

Hoặc nếu có lỗi:
```json
{
  "success": false,
  "error": "Ollama AI is not available",
  "message": "Please configure --ollama-url and --ollama-model flags..."
}
```

## ⚠️ Lưu ý

- **Timeout**: Mặc định timeout là 120 giây. Câu hỏi phức tạp có thể cần thời gian lâu hơn
- **Model size**: Model lớn (13B+) sẽ chính xác hơn nhưng chậm hơn
- **RAM**: Đảm bảo có đủ RAM cho model (thường cần 8-16GB cho model 7B-13B)
- **Context length**: Một số model có giới hạn độ dài context, chia nhỏ conversation nếu quá dài

## 🔗 Kết hợp với H-AI Tools

Bạn có thể kết hợp chat với các tools của H-AI:

```bash
# 1. Scan target
curl -X POST http://localhost:8888/api/tools/nmap \
  -d '{"target": "example.com"}' > scan_results.json

# 2. Hỏi LLM về kết quả scan
curl -X POST http://localhost:8888/api/intelligence/chat \
  -d '{
    "messages": [{
      "role": "user",
      "content": "Tôi scan được kết quả sau, hãy phân tích:\n'$(cat scan_results.json)'"
    }]
  }'
```

Happy hacking! 🚀

