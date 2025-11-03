# Cấu trúc dự án H-AI


## Cấu trúc thư mục

```
h-ai/
├── cmd/                    # Command-line applications
│   └── mcp/               # MCP client application
│       └── main.go        # MCP client entry point
│
├── internal/               # Internal packages (not exported)
│   ├── cache/             # Caching system
│   │   └── cache.go       # LRU cache implementation
│   │
│   ├── client/            # HTTP API client
│   │   └── client.go      # Client for communicating with API server
│   │
│   ├── executor/          # Command execution engine
│   │   └── executor.go    # Process management and command execution
│   │
│   ├── models/            # Data models
│   │   └── requests.go    # Request/response models
│   │
│   ├── recovery/          # Error recovery system
│   │   └── recovery.go    # Recovery strategies and handlers
│   │
│   ├── server/            # HTTP server
│   │   ├── server.go      # Server setup and routing
│   │   └── handlers.go    # HTTP request handlers
│   │
│   ├── tools/             # Security tools manager
│   │   └── manager.go     # Tool execution and management
│   │
│   └── utils/             # Utility functions
│       └── strings.go     # String utilities
│
├── examples/              # Example code
│   └── example_client.go  # Example API client usage
│
├── main.go                # Server entry point
├── go.mod                 # Go module definition
├── go.sum                 # Go module checksums
├── Makefile              # Build automation
├── Dockerfile            # Docker image definition
├── README.md             # Main documentation
├── CONTRIBUTING.md       # Contribution guidelines
├── CHANGELOG.md          # Version history
└── h-ai-mcp.json         # MCP configuration example
```

## Các thành phần chính

### 1. Server (`main.go`, `internal/server/`)

HTTP REST API server sử dụng Gin framework:
- Health check endpoints
- Security tools endpoints
- Process management endpoints
- Intelligence endpoints
- Cache management

### 2. Executor (`internal/executor/`)

Command execution engine:
- Execute shell commands với timeout
- Process tracking và management
- Output capture (stdout/stderr)
- Process termination

### 3. Cache (`internal/cache/`)

Caching system:
- In-memory cache với TTL
- Automatic cleanup
- Statistics tracking

### 4. Tools Manager (`internal/tools/`)

Security tools integration:
- Tool availability checking
- Command building
- Tool execution wrappers
- Result formatting

### 5. Recovery (`internal/recovery/`)

Error recovery system:
- Recovery strategies
- Backoff calculation
- Error handling

### 6. Client (`internal/client/`)

HTTP API client:
- GET/POST requests
- Health checking
- Error handling

### 7. MCP Client (`cmd/mcp/`)

MCP protocol client:
- Connect to API server
- Expose tools to AI agents
- Protocol handling

## Luồng hoạt động

1. **Startup**: Server khởi động và thiết lập routes
2. **Request**: Client gửi HTTP request đến API endpoint
3. **Handler**: Handler xử lý request và gọi tool manager
4. **Tool Execution**: Tool manager build command và gọi executor
5. **Process**: Executor chạy command và track process
6. **Result**: Kết quả được cache và trả về client
7. **Recovery**: Nếu có lỗi, recovery system xử lý

## Dependencies

- **Gin**: HTTP web framework
- **Zap**: Structured logging
- **go-cache**: Caching utilities (not actively used, using custom implementation)

## Build và chạy

```bash
# Build
make build

# Run server
make run

# Test
make test
```

## Mở rộng

Để thêm tool mới:
1. Thêm model trong `internal/models/requests.go`
2. Thêm handler trong `internal/server/handlers.go`
3. Thêm execution method trong `internal/tools/manager.go`
4. Update routes trong `internal/server/server.go`
