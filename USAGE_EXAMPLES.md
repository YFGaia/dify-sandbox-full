# Dify-Sandbox MCP Server 使用示例

本文档提供了 Dify-Sandbox MCP Server 的详细使用示例，重点介绍新的 StreamableHTTP 传输模式。

## 🚀 快速开始

### 1. 构建 Docker 镜像

```bash
# 使用构建脚本
./build/build_docker_mcp.sh

# 或直接使用 Docker 命令
docker build -f docker/amd64/dockerfile.mcp -t dify-sandbox-mcp:latest .
```

### 2. 启动 MCP Server

#### StreamableHTTP 模式（默认推荐）

```bash
# 本地测试
docker run -p 3000:3000 dify-sandbox-mcp:latest

# 远程访问
docker run -p 3000:3000 \
  -e MCP_BASE_URL=http://your-server-ip:3000 \
  dify-sandbox-mcp:latest

# 生产环境
docker run -d \
  --name dify-sandbox-mcp \
  -p 3000:3000 \
  -e MCP_TRANSPORT=streamable-http \
  -e MCP_BASE_URL=http://your-domain.com:3000 \
  -e MCP_SHOW_LOG=true \
  --restart unless-stopped \
  dify-sandbox-mcp:latest
```

#### SSE 模式（备选）

```bash
docker run -p 3000:3000 \
  -e MCP_TRANSPORT=sse \
  -e MCP_BASE_URL=http://your-server-ip:3000 \
  dify-sandbox-mcp:latest
```

#### STDIO 模式（AI 助手集成）

```bash
# 直接运行
docker run --rm -i dify-sandbox-mcp:latest

# Claude Desktop 配置
{
  "mcpServers": {
    "dify-sandbox": {
      "command": "docker",
      "args": ["run", "--rm", "-i", "dify-sandbox-mcp:latest"]
    }
  }
}
```

## 🔧 MCP 客户端测试

### 使用 MCP Inspector

1. **访问 MCP Inspector**：https://modelcontextprotocol.io/docs/tools/inspector

2. **StreamableHTTP 连接**：
   - 连接类型：选择 "StreamableHTTP"
   - URL：`http://your-server:3000/mcp`
   - 点击 "Connect"

3. **SSE 连接**：
   - 连接类型：选择 "Server-Sent Events (SSE)"
   - SSE URL：`http://your-server:3000/sse`
   - 点击 "Connect"

### 使用 curl 测试 StreamableHTTP

```bash
# 获取服务器信息
curl -X POST http://localhost:3000/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 1,
    "method": "initialize",
    "params": {
      "protocolVersion": "2024-11-05",
      "capabilities": {},
      "clientInfo": {
        "name": "test-client",
        "version": "1.0.0"
      }
    }
  }'

# 列出可用工具
curl -X POST http://localhost:3000/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 2,
    "method": "tools/list"
  }'

# 执行 Python 代码
curl -X POST http://localhost:3000/mcp \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "id": 3,
    "method": "tools/call",
    "params": {
      "name": "run_python_code",
      "arguments": {
        "code": "print(\"Hello from Dify-Sandbox!\")\nprint(f\"2 + 2 = {2 + 2}\")"
      }
    }
  }'
```

## 🛠️ MCP Tools 使用示例

### 1. Python 代码执行

```json
{
  "name": "run_python_code",
  "arguments": {
    "code": "import math\nimport datetime\n\nprint(f'当前时间: {datetime.datetime.now()}')\nprint(f'圆周率: {math.pi}')\nprint(f'2的平方根: {math.sqrt(2)}')\n\n# 简单计算\nnumbers = [1, 2, 3, 4, 5]\nprint(f'数字列表: {numbers}')\nprint(f'总和: {sum(numbers)}')\nprint(f'平均值: {sum(numbers) / len(numbers)}')",
    "enable_network": false
  }
}
```

### 2. Node.js 代码执行

```json
{
  "name": "run_nodejs_code",
  "arguments": {
    "code": "console.log('Hello from Node.js!');\nconsole.log('当前时间:', new Date().toISOString());\n\n// 简单计算\nconst numbers = [1, 2, 3, 4, 5];\nconsole.log('数字数组:', numbers);\nconsole.log('总和:', numbers.reduce((a, b) => a + b, 0));\nconsole.log('平均值:', numbers.reduce((a, b) => a + b, 0) / numbers.length);",
    "enable_network": false
  }
}
```

### 3. 网络访问示例

```json
{
  "name": "run_python_code",
  "arguments": {
    "code": "import requests\nimport json\n\ntry:\n    response = requests.get('https://httpbin.org/json', timeout=10)\n    data = response.json()\n    print('网络请求成功:')\n    print(json.dumps(data, indent=2, ensure_ascii=False))\nexcept Exception as e:\n    print(f'网络请求失败: {e}')",
    "enable_network": true
  }
}
```

### 4. 依赖管理

```json
// 列出 Python 依赖
{
  "name": "list_python_dependencies"
}

// 刷新依赖列表
{
  "name": "refresh_python_dependencies"
}

// 更新依赖环境
{
  "name": "update_python_dependencies"
}
```

### 5. 健康检查

```json
{
  "name": "health_check"
}
```

## 🔧 配置选项

### 环境变量

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `MCP_TRANSPORT` | `streamable-http` | 传输模式 |
| `MCP_HTTP_PORT` | `3000` | HTTP 服务端口 |
| `MCP_BASE_URL` | `http://localhost:3000` | 客户端访问的基础 URL |
| `MCP_SHOW_LOG` | `true` | 是否显示日志 |
| `MCP_LOG_LEVEL` | `info` | 日志级别 |

### 传输模式对比

| 模式 | 适用场景 | 优势 | 端点 |
|------|----------|------|------|
| `streamable-http` | Web 应用、API 集成 | 高性能、无状态、易扩展 | `/mcp` |
| `sse` | 实时应用、浏览器集成 | 实时推送、浏览器友好 | `/sse`, `/message` |
| `stdio` | AI 助手、命令行工具 | 简单集成、无网络依赖 | 标准输入输出 |

## 🚨 故障排除

### 常见问题

1. **连接被拒绝**
   ```bash
   # 检查容器是否运行
   docker ps
   
   # 检查端口是否开放
   netstat -tlnp | grep 3000
   ```

2. **CORS 错误**
   - MCP Server 已内置 CORS 支持，允许所有来源访问

3. **工具调用失败**
   ```bash
   # 查看容器日志
   docker logs dify-sandbox-mcp
   ```

4. **网络访问被阻止**
   - 确保设置 `"enable_network": true`
   - 检查防火墙设置

### 调试模式

```bash
# 启用详细日志
docker run -p 3000:3000 \
  -e MCP_LOG_LEVEL=debug \
  -e MCP_SHOW_LOG=true \
  dify-sandbox-mcp:latest
```

## 📚 更多资源

- [Model Context Protocol 官方文档](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/mark3labs/mcp-go)
- [Dify-Sandbox MCP 设计文档](dify-sandbox-mcp-planning.md)
- [原始 Dify-Sandbox](https://github.com/langgenius/dify-sandbox) 