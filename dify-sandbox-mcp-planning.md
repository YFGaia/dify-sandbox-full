# Dify-Sandbox MCP Server 开发规划文档

## 📋 需求分析与目标拆解

### 项目背景
将 dify-sandbox 的核心代码执行功能直接集成到 MCP Server 中，通过 MCP Tools 对外提供服务。**不需要启动独立的 HTTP API 服务**，而是将现有的 `internal/service/` 层代码直接作为 MCP tools 的实现。

### 架构对比

**原始架构**：
```
HTTP Client → HTTP API Server → internal/service/ → core execution
```

**新 MCP 架构**：
```
MCP Client → MCP Server → internal/service/ → core execution
```

### 目标拆解
1. **直接集成核心模块**：将 dify-sandbox 的 `internal/service/`、`internal/core/`、`internal/static/` 等模块集成到 MCP server
2. **MCP Tools 实现**：使用 `github.com/mark3labs/mcp-go` 将服务层函数包装为 MCP tools
3. **环境初始化**：直接在 MCP server 中处理 Python/Node.js 环境的初始化和依赖管理
4. **配置管理**：复用 dify-sandbox 的配置系统
5. **单一进程部署**：整个服务打包为单一可执行文件

## 🏗️ 核心功能分析

基于对 dify-sandbox 代码的分析，需要集成以下核心功能：

### 1. 代码执行服务
- **Python3 执行**：`service.RunPython3Code()`
- **Node.js 执行**：`service.RunNodeJsCode()`
- **沙箱安全**：复用现有的 seccomp 安全机制
- **资源限制**：复用现有的超时和资源限制

### 2. 依赖管理服务
- **Python 依赖列表**：`service.ListPython3Dependencies()`
- **Python 依赖刷新**：`service.RefreshPython3Dependencies()`
- **Python 依赖更新**：`service.UpdateDependencies()`

### 3. 环境初始化
- **配置加载**：`static.InitConfig()`
- **运行环境设置**：`static.SetupRunnerDependencies()`
- **Python 环境准备**：`python.PreparePythonDependenciesEnv()`

## 🚀 技术选型建议

### MCP SDK
- **选择**：`github.com/mark3labs/mcp-go v0.30.0+`
- **传输方式**：StreamableHTTP（默认推荐）
- **备选传输**：SSE、STDIO
- **优势**：高性能、易集成、社区活跃、支持无状态部署

### 代码集成策略
- **导入方式**：直接导入 dify-sandbox 的内部包
- **初始化**：在 MCP server 启动时完成所有环境初始化
- **错误处理**：保持与原有服务层相同的错误处理逻辑

## 📐 核心架构设计

### 架构图
```
┌─────────────────┐    ┌──────────────────────┐
│   MCP Client    │    │    MCP Server        │
│   (AI助手/IDE)   │◄──►│  ┌─────────────────┐ │
└─────────────────┘    │  │   MCP Tools     │ │
                       │  │ ┌─────────────┐ │ │
                       │  │ │run_python   │ │ │
                       │  │ │run_nodejs   │ │ │
                       │  │ │list_deps    │ │ │
                       │  │ │update_deps  │ │ │
                       │  │ └─────────────┘ │ │
                       │  └─────────────────┘ │
                       │  ┌─────────────────┐ │
                       │  │ Dify-Sandbox    │ │
                       │  │ Core Modules    │ │
                       │  │ ┌─────────────┐ │ │
                       │  │ │service/     │ │ │
                       │  │ │core/        │ │ │
                       │  │ │static/      │ │ │
                       │  │ │utils/       │ │ │
                       │  │ └─────────────┘ │ │
                       │  └─────────────────┘ │
                       └──────────────────────┘
```

### 核心组件
1. **MCP Server**: 使用 mark3labs/mcp-go 创建的服务器实例
2. **Tool Handlers**: 将 service 层函数包装为 MCP tool handlers
3. **Environment Manager**: 管理 Python/Node.js 执行环境
4. **Configuration**: 复用 dify-sandbox 的配置系统
5. **Security Layer**: 保持原有的安全沙箱机制

## 🛠️ 模块/功能划分

### 项目结构
```
dify-sandbox-mcp-server/
├── cmd/
│   └── mcp-server/
│       └── main.go              # 程序入口
├── internal/
│   ├── mcp/
│   │   ├── server.go           # MCP 服务器初始化
│   │   ├── tools.go            # MCP tools 注册
│   │   └── handlers/           # MCP tool handlers
│   │       ├── code_execution.go  # 代码执行相关 handlers
│   │       ├── dependencies.go    # 依赖管理相关 handlers
│   │       └── health.go          # 健康检查 handler
│   ├── service/                # 复用 dify-sandbox
│   ├── core/                   # 复用 dify-sandbox  
│   ├── static/                 # 复用 dify-sandbox
│   ├── utils/                  # 复用 dify-sandbox
│   ├── types/                  # 复用 dify-sandbox
│   └── middleware/             # 复用 dify-sandbox
├── conf/
│   └── mcp-config.yaml         # MCP 服务器配置文件
├── dependencies/               # 复用 dify-sandbox
├── go.mod
├── go.sum
└── README.md
```

### 模块职责

#### 1. MCP 层 (`internal/mcp/`)
- **server.go**: MCP 服务器创建、配置、启动
- **tools.go**: 所有 MCP tools 的定义和注册
- **handlers/**: 各个 tool 的具体实现逻辑

#### 2. 服务层（复用 dify-sandbox）
- **service/**: 代码执行和依赖管理的核心业务逻辑
- **core/**: 底层执行引擎（Python、Node.js runners）
- **static/**: 配置管理和环境设置
- **utils/**: 工具函数
- **types/**: 数据类型定义

## 🔌 详细的 MCP Tools 设计

### 1. `run_python_code` Tool
```go
mcp.NewTool("run_python_code",
    mcp.WithDescription("在安全沙箱环境中执行 Python 代码"),
    mcp.WithString("code",
        mcp.Required(),
        mcp.Description("要执行的 Python 代码"),
    ),
    mcp.WithString("preload",
        mcp.Description("预加载的代码（可选）"),
    ),
    mcp.WithBoolean("enable_network",
        mcp.Description("是否启用网络访问"),
        mcp.DefaultBool(false),
    ),
)
```

**Handler 实现**：
```go
func (h *CodeExecutionHandler) RunPythonCode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    code, _ := request.RequireString("code")
    preload, _ := request.OptionalString("preload")
    enableNetwork, _ := request.OptionalBool("enable_network")
    
    // 直接调用 dify-sandbox 的服务层
    result := service.RunPython3Code(code, preload, &runner_types.RunnerOptions{
        EnableNetwork: enableNetwork,
    })
    
    return convertToMCPResult(result), nil
}
```

### 2. `run_nodejs_code` Tool
```go
mcp.NewTool("run_nodejs_code",
    mcp.WithDescription("在安全沙箱环境中执行 Node.js 代码"),
    mcp.WithString("code",
        mcp.Required(),
        mcp.Description("要执行的 Node.js 代码"),
    ),
    mcp.WithString("preload",
        mcp.Description("预加载的代码（可选）"),
    ),
    mcp.WithBoolean("enable_network",
        mcp.Description("是否启用网络访问"),
        mcp.DefaultBool(false),
    ),
)
```

### 3. `list_python_dependencies` Tool
```go
mcp.NewTool("list_python_dependencies",
    mcp.WithDescription("获取 Python 环境中已安装的依赖包列表"),
)
```

### 4. `refresh_python_dependencies` Tool
```go
mcp.NewTool("refresh_python_dependencies",
    mcp.WithDescription("刷新并获取最新的 Python 依赖包列表"),
)
```

### 5. `update_python_dependencies` Tool
```go
mcp.NewTool("update_python_dependencies",
    mcp.WithDescription("更新 Python 依赖环境"),
)
```

### 6. `health_check` Tool
```go
mcp.NewTool("health_check",
    mcp.WithDescription("检查代码执行环境的健康状态"),
)
```

## 📊 数据结构设计

### MCP Tool Response 格式
```go
// 代码执行结果
type CodeExecutionResult struct {
    Success bool   `json:"success"`
    Stdout  string `json:"stdout"`
    Stderr  string `json:"stderr"`
    Error   string `json:"error,omitempty"`
}

// 依赖列表结果
type DependencyListResult struct {
    Success      bool                          `json:"success"`
    Dependencies []runner_types.Dependency    `json:"dependencies"`
}

// 健康检查结果
type HealthCheckResult struct {
    Status  string            `json:"status"`
    Details map[string]string `json:"details"`
}
```

## 📝 详细的待办事项列表 (Todolist)

### Phase 1: 项目设置和模块集成 🏗️
- [x] **1.1** 创建新项目目录结构
- [x] **1.2** 初始化 Go module
- [x] **1.3** 从 dify-sandbox 复制核心模块（service、core、static、utils、types）
- [x] **1.4** 集成 `github.com/mark3labs/mcp-go` 依赖
- [x] **1.5** 复制并调整配置文件和依赖文件
- [x] **1.6** 解决模块间的导入路径问题

### Phase 2: MCP Server 基础架构 🔧
- [x] **2.1** 实现 MCP server 初始化 (`internal/mcp/server.go`)
- [x] **2.2** 创建环境初始化逻辑（复用 dify-sandbox 的初始化代码）
- [x] **2.3** 实现配置加载和管理
- [x] **2.4** 设置日志系统
- [x] **2.5** 创建错误处理和响应转换工具

### Phase 3: MCP Tools 实现 🛠️
- [x] **3.1** 实现 `run_python_code` tool 和 handler
- [x] **3.2** 实现 `run_nodejs_code` tool 和 handler  
- [x] **3.3** 实现 `list_python_dependencies` tool 和 handler
- [x] **3.4** 实现 `refresh_python_dependencies` tool 和 handler
- [x] **3.5** 实现 `update_python_dependencies` tool 和 handler
- [x] **3.6** 实现 `health_check` tool 和 handler
- [x] **3.7** 实现响应格式转换函数

### Phase 4: 集成和测试 🧪
- [x] **4.1** 完成主程序入口 (`cmd/mcp-server/main.go`)
- [x] **4.2** 测试环境初始化流程
- [x] **4.3** 测试各个 MCP tools 的功能（简化版本）
- [x] **4.4** 实现 StreamableHTTP 传输模式
- [x] **4.5** 将 StreamableHTTP 设置为默认传输模式
- [ ] **4.6** 测试错误处理和边缘情况
- [ ] **4.7** 性能测试和优化

### Phase 5: 文档和部署 📚
- [x] **5.1** 编写 README 文档
- [x] **5.2** 更新文档以反映 StreamableHTTP 默认模式
- [ ] **5.3** 编写使用示例和集成指南
- [ ] **5.4** 创建构建和部署脚本
- [ ] **5.5** 编写故障排除指南

### Phase 6: Linux 环境完整版本 🐧
- [ ] **6.1** 在 Linux 环境中构建底层库
- [ ] **6.2** 测试完整版本的代码执行功能
- [ ] **6.3** 验证所有 MCP tools 在 Linux 环境下的工作
- [ ] **6.4** 性能优化和稳定性测试

### Phase 7: 高级功能 🚀
- [ ] **7.1** 实现会话管理（如果需要）
- [ ] **7.2** 添加缓存机制
- [ ] **7.3** 实现工具执行监控和统计
- [ ] **7.4** 添加配置热重载
- [ ] **7.5** 支持多个 dify-sandbox 实例的负载均衡

## 🎉 当前进展总结

### ✅ 已完成的重要里程碑

1. **MCP Server 架构完成** - 成功集成 mark3labs/mcp-go SDK v0.30.0
2. **所有 MCP Tools 实现** - 6个核心工具全部实现并注册
3. **Handler 层完成** - 代码执行、依赖管理、健康检查处理器
4. **StreamableHTTP 传输模式** - 实现并设置为默认传输模式
5. **多传输模式支持** - 支持 StreamableHTTP、SSE、STDIO 三种模式
6. **简化版本可用** - 在 macOS 上成功构建和测试
7. **MCP 协议验证** - 工具列表和调用功能正常工作
8. **文档更新完成** - README 和规划文档已更新

### 🚧 当前状态

- **StreamableHTTP 模式**：已实现并设为默认，支持无状态部署
- **简化版本**：完全可用，适合 MCP 协议测试和开发
- **完整版本**：需要 Linux 环境进行最终测试
- **文档**：基础文档和配置说明已完成

### 📋 下一步重点

1. **Linux 环境测试** - 验证完整的代码执行功能
2. **错误处理优化** - 完善异常情况的处理
3. **使用示例** - 创建更多实际使用场景的示例

## 🎯 实现策略

### 模块复用策略
1. **完整复用**：`internal/service/`, `internal/core/`, `internal/static/`, `internal/utils/`, `internal/types/`
2. **选择性复用**：`conf/`, `dependencies/`（根据需要调整）
3. **跳过复用**：`internal/controller/`, `internal/server/`, `internal/middleware/`（HTTP 相关）

### 初始化流程
```go
// 在 main.go 中的初始化顺序
func main() {
    // 1. 加载配置
    err := static.InitConfig("conf/config.yaml")
    
    // 2. 设置运行环境
    err = static.SetupRunnerDependencies()
    
    // 3. 初始化 Python 环境（异步）
    go initializePythonEnvironment()
    
    // 4. 创建和启动 MCP 服务器
    server := mcp.NewMCPServer(...)
    mcp.RegisterAllTools(server)
    server.ServeStdio()
}
```

### 错误处理策略
- **保持一致性**：与 dify-sandbox 的错误处理保持一致
- **MCP 适配**：将 `types.DifySandboxResponse` 转换为 MCP 响应格式
- **详细日志**：保持详细的错误日志记录

## 🔒 安全考虑

1. **沙箱隔离**：完全依赖 dify-sandbox 现有的安全机制
2. **资源限制**：保持原有的超时和资源限制配置
3. **网络控制**：默认禁用网络访问，用户可选择启用
4. **依赖管理**：只允许管理预定义的 Python 依赖

## 📈 成功指标

1. **功能完整性**：所有 dify-sandbox 核心功能都能通过 MCP tools 使用
2. **性能指标**：执行性能与原 dify-sandbox 相当或更好
3. **稳定性**：长时间运行无内存泄漏或崩溃
4. **易用性**：AI 助手能够轻松调用和理解各个工具
5. **兼容性**：与主流 MCP 客户端完全兼容