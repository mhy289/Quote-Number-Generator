# QuoteGenerator 优化待办清单

> 生成时间：2026-07-29

---

## 🔴 P0 - 严重问题（优先修复）

### 1. 后端随机数种子初始化 Bug
- **文件**: `backend/main.go:60`
- **问题**: `rand.New(rand.NewSource(...))` 创建新实例但未赋值给全局 `rand`，实际仍使用默认固定种子
- **方案**: 改为 `rand.Seed(time.Now().UnixNano())` 或创建局部实例使用

### 2. 后端名言文件路径硬编码
- **文件**: `backend/main.go:68`
- **问题**: `loadQuotes("data/quotes.json")` 使用相对路径，从不同目录启动会找不到文件
- **方案**: 使用 `os.Executable()` 获取可执行文件路径，或通过命令行参数传入

---

## 🟡 P1 - 中等优化

### 3. 前端后端地址硬编码
- **文件**: `frontend/src/app/page.tsx`
- **问题**: `fetch("http://localhost:8080/...")` 硬编码，部署到不同环境需手动修改
- **方案**: 使用 Next.js 环境变量 `NEXT_PUBLIC_API_URL`，或通过 API Route 代理转发

### 4. 前端跨域依赖
- **问题**: 前端 `localhost:3000` 直接请求 `localhost:8080`，依赖后端 CORS 中间件
- **方案**: 利用 Next.js API Routes 做代理转发，消除跨域依赖，同时隐藏后端地址

### 5. 后端代码拆分
- **文件**: `backend/main.go`（207 行）
- **问题**: 服务定义、HTTP 启动、控制台交互全部在一个文件中，职责不清晰
- **方案**: 拆分为 `service.go`（业务逻辑）、`server.go`（HTTP 启动）、`console.go`（控制台交互）

### 6. 后端缺少优雅关闭
- **文件**: `backend/main.go:78`
- **问题**: `select{}` 永久阻塞，无法响应 `SIGINT/SIGTERM`
- **方案**: 使用 `signal.Notify` 监听系统信号实现优雅关闭

### 7. 启动脚本健康检查不准确
- **文件**: `start.bat`
- **问题**: 用 `GET /` 检查后端，但 ConnectRPC 没有根路由，始终返回 404
- **方案**: 后端添加 `/health` 健康检查端点

### 8. 前端缺少错误边界和骨架屏
- **文件**: `frontend/src/app/page.tsx`
- **问题**: 仅用文字提示加载状态，用户体验单一
- **方案**: 添加 `Suspense` 边界、骨架屏组件、重试按钮

### 9. 前端请求状态管理分散
- **文件**: `frontend/src/app/page.tsx`
- **问题**: 每个请求独立管理 `loading`、`error`、`data` 三个状态，重复代码多
- **方案**: 封装自定义 Hook（如 `useRequest`）统一管理请求状态

### 10. 前端缺少 API 类型定义
- **文件**: `frontend/src/app/page.tsx`
- **问题**: API 请求/响应没有定义接口类型，`data.number` 和 `data.quote` 无类型提示
- **方案**: 定义 `ApiResponse` 类型，或从 proto 生成 TypeScript 类型

---

## 🟢 P2 - 轻微优化

### 11. 前端元数据未更新
- **文件**: `frontend/src/app/layout.tsx`
- **问题**: `title` 和 `description` 仍是 `create next app` 默认值
- **方案**: 改为 `"QuoteGenerator"` 和项目描述

### 12. 前端语言标签未更新
- **文件**: `frontend/src/app/layout.tsx`
- **问题**: `<html lang="en">` 应改为中文
- **方案**: 改为 `lang="zh-CN"`

### 13. 名言抓取工具缺少超时控制
- **文件**: `backend/tools/sources/source.go`
- **问题**: `http.Get(s.URL)` 无超时，可能永久阻塞
- **方案**: 使用 `http.Client` 并设置 `Timeout` 字段

### 14. 测试覆盖率不足
- **文件**: `backend/generator_service_test.go`
- **问题**: 仅有后端单元测试，缺少前端测试和集成测试
- **方案**: 添加前端组件测试（Vitest + Testing Library）和后端 HTTP 集成测试

### 15. `.gitignore` 缺少 `.next` 目录
- **文件**: `.gitignore`
- **问题**: `frontend/.next/` 构建产物未在 `.gitignore` 中排除
- **方案**: 在根 `.gitignore` 中添加 `frontend/.next/`

---

## 📋 优先级总览

| 优先级 | 数量 | 说明 |
|--------|------|------|
| 🔴 P0 | 2 | 功能正确性 & 部署稳定性 |
| 🟡 P1 | 8 | 架构合理性 & 可维护性 |
| 🟢 P2 | 5 | 细节体验 & 质量保障 |
| **合计** | **15** | |

> 建议按 P0 → P1 → P2 的顺序逐步推进优化。