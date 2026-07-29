#!/bin/bash

ROOT_DIR="$(cd "$(dirname "$0")" && pwd)"
BACKEND_DIR="$ROOT_DIR/backend"
FRONTEND_DIR="$ROOT_DIR/frontend"

echo "========================================"
echo "  🚀 QuoteGenerator 启动中..."
echo "========================================"
echo ""

# 启动后端服务
echo "📦 启动后端服务..."
cd "$BACKEND_DIR" || exit
go run . --mode=http &
BACKEND_PID=$!

sleep 3

# 启动前端服务
echo "🎨 启动前端服务..."
cd "$FRONTEND_DIR" || exit
npm run dev &
FRONTEND_PID=$!

echo ""
echo "========================================"
echo "  ✅ 服务启动完成！"
echo "  🌐 前端: http://localhost:3000"
echo "  🔧 后端: http://localhost:8080"
echo "========================================"
echo ""
echo "按 Ctrl+C 停止所有服务"

# 捕获退出信号，清理子进程
cleanup() {
    echo ""
    echo "🛑 正在停止所有服务..."
    kill $BACKEND_PID 2>/dev/null
    kill $FRONTEND_PID 2>/dev/null
    echo "👋 已退出"
    exit 0
}

trap cleanup SIGINT SIGTERM

# 等待任意子进程结束
wait