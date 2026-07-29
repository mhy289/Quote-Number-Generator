# QuoteGenerator - 随机名言/数字生成器

一个基于 **Go + ConnectRPC** 后端和 **Next.js** 前端的全栈项目，提供随机名言和随机数字生成服务。

## 项目架构

```
QuoteGenerator/
├── backend/          # Go 后端服务 (ConnectRPC)
│   ├── main.go                    # 业务逻辑 + 服务入口
│   ├── generator.proto            # Protobuf 接口定义
│   ├── generatorpb/               # 自动生成的 Protobuf 消息类型
│   ├── generatorpb/generatorpbconnect/  # 自动生成的 ConnectRPC 骨架
│   ├── generator_service_test.go  # 单元测试
│   ├── go.mod / go.sum            # Go 依赖管理
│   └── README.md                  # 后端文档
├── frontend/         # Next.js 前端应用
│   ├── app/                       # 页面组件
│   ├── README.md                  # 前端文档
│   └── ...
└── README.md         # 本文件（项目总览）
```

## 技术栈

| 层级 | 技术 | 说明 |
|---|---|---|
| **后端** | Go 1.24 + ConnectRPC + Protocol Buffers | RPC 服务，支持 HTTP/1.1 和 JSON |
| **前端** | Next.js (App Router) | React 全栈框架 |
| **通信** | ConnectRPC 协议 | 支持 Connect/gRPC/gRPC-Web |

## 快速启动

### 1. 启动后端服务

```bash
cd backend
go run main.go
# 服务运行在 http://localhost:8080
```

### 2. 启动前端应用

```bash
cd frontend
npm install
npm run dev
# 应用运行在 http://localhost:3000
```

## API 接口

后端提供两个 RPC 接口（通过 HTTP POST 调用）：

| 接口路径 | 请求体 | 响应体 | 说明 |
|---|---|---|---|
| `/generator.GeneratorService/GetRandomNumber` | `{"min": 1, "max": 100}` | `{"number": 42}` | 生成指定范围随机整数 |
| `/generator.GeneratorService/GetRandomQuote` | `{}` | `{"quote": "..."}` | 返回随机名言 |

## 开发

- **修改接口定义**：编辑 `backend/generator.proto`，然后运行 `protoc` 重新生成代码
- **添加名言**：编辑 `backend/main.go` 中的 `quotes` 切片
- **运行测试**：`cd backend && go test -v`