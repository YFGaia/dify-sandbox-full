#!/bin/bash

# Dify-Sandbox MCP Docker Build Script
# 构建 MCP Server 模式的 Docker 镜像

set -e

# 颜色定义
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# 打印带颜色的消息
print_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

print_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

print_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

print_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

# 检查 Docker 是否安装
check_docker() {
    if ! command -v docker &> /dev/null; then
        print_error "Docker is not installed. Please install Docker first."
        exit 1
    fi
    print_info "Docker is available"
}

# 检查 Docker 是否运行
check_docker_running() {
    if ! docker info &> /dev/null; then
        print_error "Docker is not running. Please start Docker first."
        exit 1
    fi
    print_info "Docker is running"
}

# 获取项目根目录
get_project_root() {
    SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
    PROJECT_ROOT="$(dirname "$SCRIPT_DIR")"
    echo "$PROJECT_ROOT"
}

# 构建 Docker 镜像
build_image() {
    local project_root="$1"
    local image_name="${2:-dify-sandbox-mcp}"
    local tag="${3:-latest}"
    local dockerfile="docker/amd64/dockerfile.mcp"
    
    print_info "Building Docker image: ${image_name}:${tag}"
    print_info "Project root: ${project_root}"
    print_info "Dockerfile: ${dockerfile}"
    
    # 检查 Dockerfile 是否存在
    if [[ ! -f "${project_root}/${dockerfile}" ]]; then
        print_error "Dockerfile not found: ${project_root}/${dockerfile}"
        exit 1
    fi
    
    # 构建镜像
    cd "$project_root"
    docker build \
        -f "$dockerfile" \
        -t "${image_name}:${tag}" \
        . || {
        print_error "Failed to build Docker image"
        exit 1
    }
    
    print_success "Docker image built successfully: ${image_name}:${tag}"
}

# 显示镜像信息
show_image_info() {
    local image_name="$1"
    local tag="$2"
    
    print_info "Image information:"
    docker images "${image_name}:${tag}" --format "table {{.Repository}}\t{{.Tag}}\t{{.ID}}\t{{.CreatedAt}}\t{{.Size}}"
}

# 显示使用说明
show_usage() {
    echo "Usage: $0 [OPTIONS]"
    echo ""
    echo "Options:"
    echo "  -n, --name NAME     Docker image name (default: dify-sandbox-mcp)"
    echo "  -t, --tag TAG       Docker image tag (default: latest)"
    echo "  -h, --help          Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0                                    # Build with default settings"
    echo "  $0 -n my-sandbox -t v1.0.0          # Build with custom name and tag"
    echo "  $0 --name dify-mcp --tag dev         # Build development version"
}

# 显示运行示例
show_run_examples() {
    local image_name="$1"
    local tag="$2"
    
    print_info "Docker run examples:"
    echo ""
    echo "# 本地访问模式："
    echo "docker run -p 3000:3000 ${image_name}:${tag}"
    echo ""
    echo "# 远程访问模式（使用服务器 IP）："
    echo "docker run -p 3000:3000 \\"
    echo "  -e MCP_BASE_URL=http://your-server-ip:3000 \\"
    echo "  ${image_name}:${tag}"
    echo ""
    echo "# 完整配置示例："
    echo "docker run -d \\"
    echo "  --name dify-sandbox-mcp \\"
    echo "  -p 3000:3000 \\"
    echo "  -e MCP_TRANSPORT=streamable-http \\"
    echo "  -e MCP_HTTP_PORT=3000 \\"
    echo "  -e MCP_BASE_URL=http://your-server.com:3000 \\"
    echo "  -e MCP_SHOW_LOG=true \\"
    echo "  -e MCP_LOG_LEVEL=info \\"
    echo "  --restart unless-stopped \\"
    echo "  ${image_name}:${tag}"
    echo ""
    echo "# MCP Inspector 测试："
    echo "# 1. 启动容器后，打开 https://modelcontextprotocol.io/docs/tools/inspector"
    echo "# 2. StreamableHTTP 模式：选择 'HTTP' 连接类型，输入 URL: http://your-server:3000/mcp"
    echo "# 3. SSE 模式：选择 'Server-Sent Events (SSE)' 连接类型，输入 SSE URL: http://your-server:3000/sse"
    echo "# 4. 点击连接测试"
}

# 主函数
main() {
    local image_name="dify-sandbox-mcp"
    local tag="latest"
    
    # 解析命令行参数
    while [[ $# -gt 0 ]]; do
        case $1 in
            -n|--name)
                image_name="$2"
                shift 2
                ;;
            -t|--tag)
                tag="$2"
                shift 2
                ;;
            -h|--help)
                show_usage
                exit 0
                ;;
            *)
                print_error "Unknown option: $1"
                show_usage
                exit 1
                ;;
        esac
    done
    
    print_info "Starting Dify-Sandbox MCP Docker build process..."
    
    # 检查环境
    check_docker
    check_docker_running
    
    # 获取项目根目录
    project_root=$(get_project_root)
    
    # 构建镜像
    build_image "$project_root" "$image_name" "$tag"
    
    # 显示镜像信息
    show_image_info "$image_name" "$tag"
    
    # 显示运行示例
    show_run_examples "$image_name" "$tag"
    
    print_success "Build completed successfully!"
}

# 运行主函数
main "$@" 