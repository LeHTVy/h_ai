# 🔗 Tích hợp Ollama AI với H-AI

H-AI hỗ trợ tích hợp Ollama để sử dụng AI local cho phân tích thông minh và quyết định tự động.

## 📋 Yêu cầu

1. **Cài đặt Ollama**
   ```bash
   # Ubuntu/Debian
   curl -fsSL https://ollama.ai/install.sh | sh
   
   # hoặc download từ https://ollama.ai/download
   ```

2. **Tải mô hình LLM**
   ```bash
   # Ví dụ với llama2 (mô hình mặc định)
   ollama pull llama2
   
   # Hoặc các mô hình khác:
   ollama pull mistral      # Mistral 7B - tốt cho cybersecurity
   ollama pull codellama    # CodeLlama - tốt cho phân tích code/exploit
   ollama pull deepseek-chat # DeepSeek - tốt cho tiếng Việt
   ```

3. **Khởi động Ollama**
   ```bash
   # Ollama tự động chạy như một service
   # Kiểm tra xem đã chạy chưa:
   curl http://localhost:11434/api/tags
   ```

## 🚀 Sử dụng

### 1. Khởi động H-AI với Ollama

```bash
# Sử dụng mô hình mặc định (llama2) và URL mặc định (localhost:11434)
./bin/h-ai-server --ollama-model llama2

# Chỉ định URL Ollama tùy chỉnh
./bin/h-ai-server --ollama-url http://localhost:11434 --ollama-model mistral

# Hoặc tắt AI (chỉ dùng rule-based logic)
./bin/h-ai-server
```

### 2. Kiểm tra tích hợp

```bash
# Health check sẽ hiển thị trạng thái Ollama
curl http://localhost:8888/health

# Response sẽ có:
# {
#   "ollama_enabled": true,
#   "ollama_model": "llama2"
# }
```

## 🎯 Chức năng AI

### 1. Phân tích Target thông minh

```bash
curl -X POST http://localhost:8888/api/intelligence/analyze-target \
  -H "Content-Type: application/json" \
  -d '{"target": "example.com"}'
```

AI sẽ phân tích target và đưa ra:
- Đánh giá rủi ro
- Mức độ tin cậy
- Các khuyến nghị

### 2. Đề xuất Tools tối ưu

```bash
curl -X POST http://localhost:8888/api/intelligence/select-tools \
  -H "Content-Type: application/json" \
  -d '{
    "target": "https://example.com",
    "target_type": "comprehensive"
  }'
```

AI sẽ đề xuất các tools phù hợp dựa trên:
- Loại target (web app, network host, API)
- Technologies được phát hiện
- Objective (quick, stealth, comprehensive)

### 3. Tối ưu hóa Parameters

```bash
curl -X POST http://localhost:8888/api/intelligence/optimize-parameters \
  -H "Content-Type: application/json" \
  -d '{
    "tool": "nmap",
    "target": "example.com",
    "context": {"previous_scan": "ports 80,443 open"}
  }'
```

AI sẽ tối ưu parameters dựa trên:
- Context từ các scans trước
- Target profile
- Best practices cho tool đó

### 4. Phân tích Kết quả Scan

```bash
curl -X POST http://localhost:8888/api/intelligence/analyze-results \
  -H "Content-Type: application/json" \
  -d '{
    "tool": "nmap",
    "target": "example.com",
    "results": "PORT   STATE SERVICE\n80/tcp open  http\n443/tcp open https"
  }'
```

AI sẽ phân tích và đưa ra:
- Key findings và vulnerabilities tiềm năng
- Recommended next steps
- Risk assessment
- Suggested follow-up tools

## 📊 Workflow với AI

```
┌─────────────────┐
│ User Request    │
│ "Scan example"  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ AI Analyzes     │ ← Ollama
│ Target          │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ AI Suggests     │ ← Ollama
│ Tools & Params  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Execute Tools   │
│ (nmap, nuclei)  │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ AI Analyzes     │ ← Ollama
│ Results         │
└────────┬────────┘
         │
         ▼
┌─────────────────┐
│ Next Steps      │
│ Recommendations │
└─────────────────┘
```

## ⚙️ Cấu hình nâng cao

### Environment Variables

```bash
export OLLAMA_URL=http://localhost:11434
export OLLAMA_MODEL=mistral
./bin/h-ai-server --ollama-url $OLLAMA_URL --ollama-model $OLLAMA_MODEL
```

### Chọn Model phù hợp

| Model | Tốt cho | Kích thước | RAM cần |
|-------|---------|------------|---------|
| llama2 | General purpose | 7B | 8GB |
| mistral | Cybersecurity | 7B | 8GB |
| codellama | Code analysis | 7B-13B | 8-16GB |
| deepseek-chat | Tiếng Việt | 7B | 8GB |
| llava | Image analysis | 7B | 8GB |

### Tối ưu Performance

1. **Sử dụng GPU** (nếu có):
   ```bash
   # Ollama tự động dùng GPU nếu có CUDA
   # Kiểm tra:
   ollama run llama2 "test"
   ```

2. **Giảm Temperature** cho kết quả nhất quán hơn:
   - Code đã set `Temperature: 0.3` cho analysis tasks
   - Có thể điều chỉnh trong `internal/ai/ollama.go`

3. **Cache responses**: H-AI tự động cache kết quả để tránh gọi lại

## 🔍 Troubleshooting

### Ollama không kết nối được

```bash
# Kiểm tra Ollama đang chạy
curl http://localhost:11434/api/tags

# Nếu không có response, khởi động lại:
ollama serve
```

### Model chưa được tải

```bash
# Kiểm tra models đã tải:
ollama list

# Tải model:
ollama pull llama2
```

### Out of Memory

Nếu model quá lớn, dùng model nhỏ hơn:
```bash
ollama pull llama2:7b  # Thay vì 13b
```

## 📚 Tài liệu tham khảo

- [Ollama Documentation](https://github.com/ollama/ollama)
- [Available Models](https://ollama.ai/library)
- [API Reference](https://github.com/ollama/ollama/blob/main/docs/api.md)

## 💡 Tips

1. **Test model trước**: Chạy `ollama run <model>` để đảm bảo model hoạt động
2. **Monitor RAM**: LLMs tốn nhiều RAM, đảm bảo có đủ tài nguyên
3. **Fallback tự động**: Nếu Ollama không available, H-AI tự động dùng rule-based logic
4. **Kết hợp với MCP**: Có thể dùng cả Ollama (local) và Claude Desktop (cloud) cùng lúc

