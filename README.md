# Dify-Sandbox 代码执行沙箱

提供安全的 Python 和 Node.js 代码执行能力，支持两种服务模式：

1. **HTTP API 模式** - 传统的 REST API 服务
2. **MCP Server 模式** - 基于 Model Context Protocol 的新一代服务

## 🚀 功能特性

### 核心能力
- **Python 代码执行** - 安全沙箱中执行 Python 3.10 代码
- **Node.js 代码执行** - 安全沙箱中执行 Node.js 20.x 代码
- **依赖管理** - Python 包的安装、列表、更新管理
- **安全隔离** - 基于 seccomp 的系统调用限制
- **资源控制** - 内存、CPU 时间、文件访问限制

### 安全特性
- **沙箱隔离**：基于 libseccomp 的系统调用过滤
- **资源限制**：内存、CPU 时间、文件访问限制
- **网络控制**：可选的网络访问控制
- **文件系统隔离**：受限的文件系统访问权限

## 📦 系统要求

- **操作系统**：Linux (推荐)，macOS (开发/测试)
- **依赖**：Docker 或 libseccomp, pkg-config, gcc, Go 1.23+
- **架构**：AMD64, ARM64

---

## 🌐 HTTP API 模式

传统的 REST API 服务，提供 HTTP 接口进行代码执行。

### Docker 快速启动

#### 构建镜像
```bash
# 使用构建脚本
./build/build_docker.sh

# 或直接使用 Docker 命令
docker build -f docker/amd64/dockerfile -t dify-sandbox:latest .
```

#### 运行容器
```bash
# 基本运行
docker run -p 8194:8194 dify-sandbox:latest

# 自定义端口
docker run -p 9000:8194 dify-sandbox:latest

# 启用网络访问
docker run -p 8194:8194 -e ENABLE_NETWORK=true dify-sandbox:latest

# 完整配置示例
docker run -d \
  --name dify-sandbox \
  -p 8194:8194 \
  -e ENABLE_NETWORK=false \
  -e WORKER_TIMEOUT=60 \
  -e MAX_WORKERS=4 \
  --restart unless-stopped \
  dify-sandbox:latest
```

### API 端点

#### 代码执行
- **POST** `/v1/sandbox/run` - 执行代码

**请求参数：**
```json
{
  "language": "python3",  // 或 "nodejs"
  "code": "print('Hello World')",
  "preload": "",  // 可选，预加载代码
  "enable_network": false  // 可选，是否启用网络
}
```

#### 依赖管理
- **GET** `/v1/sandbox/dependencies` - 获取依赖列表
- **POST** `/v1/sandbox/dependencies/refresh` - 刷新依赖
- **POST** `/v1/sandbox/dependencies/update` - 更新依赖

### 环境变量配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `PORT` | `8194` | HTTP 服务端口 |
| `ENABLE_NETWORK` | `false` | 是否启用网络访问 |
| `MAX_WORKERS` | `4` | 最大工作进程数 |
| `WORKER_TIMEOUT` | `60` | 工作进程超时时间(秒) |

---

## 🔌 MCP Server 模式

基于 Model Context Protocol (MCP) 的新一代服务模式，专为 AI 助手和智能工具设计。

> 📖 **详细设计文档**：请参考 [dify-sandbox-mcp-planning.md](dify-sandbox-mcp-planning.md) 了解完整的架构设计和实现细节。

### MCP Tools 概览

#### 代码执行工具
- **`run_python_code`** - 在安全沙箱中执行 Python 代码
- **`run_nodejs_code`** - 在安全沙箱中执行 Node.js 代码

#### 依赖管理工具
- **`list_python_dependencies`** - 列出已安装的 Python 依赖包
- **`refresh_python_dependencies`** - 刷新 Python 依赖包列表
- **`update_python_dependencies`** - 更新 Python 依赖环境

#### 系统工具
- **`health_check`** - 检查代码执行环境的健康状态

### Docker 构建与运行

#### 构建 MCP 镜像
```bash
# 使用构建脚本
./build/build_docker_mcp.sh

# 或直接使用 Docker 命令
docker build -f docker/amd64/dockerfile.mcp -t dify-sandbox-mcp:latest .
```

#### 运行 MCP 容器

**本地访问模式：**
```bash
# 基本运行
docker run -p 3000:3000 dify-sandbox-mcp:latest

# 自定义端口
docker run -p 8080:3000 \
  -e MCP_HTTP_PORT=3000 \
  dify-sandbox-mcp:latest
```

**远程访问模式：**
```bash
# 使用服务器 IP
docker run -p 3000:3000 \
  -e MCP_BASE_URL=http://your-server-ip:3000 \
  dify-sandbox-mcp:latest

# 使用域名
docker run -p 3000:3000 \
  -e MCP_BASE_URL=http://your-domain.com:3000 \
  dify-sandbox-mcp:latest

# 完整配置示例
docker run -d \
  --name dify-sandbox-mcp \
  -p 3000:3000 \
  -e MCP_TRANSPORT=sse \
  -e MCP_HTTP_PORT=3000 \
  -e MCP_BASE_URL=http://your-server.com:3000 \
  -e MCP_SHOW_LOG=true \
  -e MCP_LOG_LEVEL=info \
  --restart unless-stopped \
  dify-sandbox-mcp:latest
```

### MCP 环境变量配置

| 变量名 | 默认值 | 说明 |
|--------|--------|------|
| `MCP_TRANSPORT` | `sse` | 传输模式：`stdio`, `sse` |
| `MCP_HTTP_PORT` | `3000` | HTTP 服务端口 (SSE 模式) |
| `MCP_BASE_URL` | `http://localhost:3000` | 客户端访问的基础 URL |
| `MCP_SHOW_LOG` | `true` | 是否显示日志 |
| `MCP_LOG_LEVEL` | `info` | 日志级别：`debug`, `info`, `warn`, `error` |

### MCP 客户端集成

#### 与 AI 助手集成 (STDIO)

**Claude Desktop 配置：**
```json
{
  "mcpServers": {
    "dify-sandbox": {
      "command": "docker",
      "args": ["run", "--rm", "-i", "dify-sandbox-mcp:latest"]
    }
  }
}
```

#### HTTP 客户端集成 (SSE)

**端点访问：**
- **SSE 连接**: `http://your-server:3000/sse`
- **消息发送**: `http://your-server:3000/message`

**MCP Inspector 测试：**
1. 打开 [MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector)
2. 选择 "Server-Sent Events (SSE)" 连接类型
3. 输入 SSE URL：`http://your-server:3000/sse`
4. 点击连接测试

## 🛠️ MCP Tools 详细说明

### 代码执行工具

#### `run_python_code`
在安全沙箱环境中执行 Python 代码。

**参数：**
- `code` (string, 必需) - 要执行的 Python 代码
- `preload` (string, 可选) - 预加载的代码，在主代码执行前运行
- `enable_network` (boolean, 可选, 默认: false) - 是否启用网络访问

**示例：**
```json
{
  "name": "run_python_code",
  "arguments": {
    "code": "import math\nprint(f'圆周率: {math.pi}')\nprint(f'2的平方根: {math.sqrt(2)}')",
    "enable_network": false
  }
}
```

#### `run_nodejs_code`
在安全沙箱环境中执行 Node.js 代码。

**参数：**
- `code` (string, 必需) - 要执行的 Node.js 代码
- `preload` (string, 可选) - 预加载的代码
- `enable_network` (boolean, 可选, 默认: false) - 是否启用网络访问

**示例：**
```json
{
  "name": "run_nodejs_code",
  "arguments": {
    "code": "console.log('Hello from Node.js!'); console.log('当前时间:', new Date().toISOString());",
    "enable_network": false
  }
}
```

### 依赖管理工具

#### `list_python_dependencies`
获取当前 Python 环境中已安装的依赖包列表。

**参数：** 无

**返回：** 包含包名、版本等信息的依赖列表

#### `refresh_python_dependencies`
刷新并重新扫描 Python 依赖包列表。

**参数：** 无

#### `update_python_dependencies`
更新 Python 依赖环境，重新安装和配置依赖包。

**参数：** 无

### 系统工具

#### `health_check`
检查代码执行环境的健康状态，包括 Python、Node.js 环境和系统资源。

**参数：** 无

**返回：** 系统健康状态报告

## 🏗️ 架构对比

### HTTP API 架构
```
HTTP Client → HTTP API Server → internal/service/ → core execution
```

### MCP 架构
```
MCP Client → MCP Server → internal/service/ → core execution
```

### MCP 优势
1. **协议标准化** - 遵循 MCP 规范，与 AI 助手无缝集成
2. **传输灵活** - 支持 STDIO、SSE 等多种传输方式
3. **类型安全** - 强类型的工具定义和参数验证
4. **易于集成** - AI 助手可直接理解和调用工具

## 🔒 安全考虑

1. **沙箱隔离**：基于 libseccomp 的系统调用过滤
2. **资源限制**：CPU、内存、执行时间限制
3. **网络控制**：默认禁用网络访问，可选择性启用
4. **文件系统**：受限的文件系统访问权限
5. **进程隔离**：每次执行都在独立的进程中运行

## 📄 许可证

本项目采用 Apache-2.0 许可证。详见 [LICENSE](LICENSE) 文件。

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📚 相关文档

- [Model Context Protocol (MCP)](https://modelcontextprotocol.io/)
- [MCP Go SDK](https://github.com/mark3labs/mcp-go)
- [Sandbox MCP 详细设计文档](dify-sandbox-mcp-planning.md)
- [原始 Dify-Sandbox](https://github.com/langgenius/dify-sandbox)

---

**注意**：在生产环境中建议使用 Linux 系统以获得完整的沙箱安全特性。