# QuoteGenerator 后端服务

基于 **Go + ConnectRPC + Protocol Buffers** 构建的随机数/名言生成器后端服务。

## 技术栈

| 技术 | 版本 | 用途 |
|---|---|---|
| Go | 1.24.5 | 编程语言 |
| ConnectRPC | v1.18.1 | RPC 框架（轻量级 gRPC 替代） |
| Protocol Buffers | proto3 | 接口定义与序列化 |
| rs/cors | v1.11.1 | CORS 中间件 |

## 项目结构

```
backend/
├── main.go                              # 业务逻辑 + 服务启动入口
├── generator.proto                      # Protobuf 接口定义文件
├── generator_service_test.go            # 单元测试
├── go.mod                               # Go 模块定义
├── go.sum                               # 依赖校验文件
├── README.md                            # 本文件
│
├── data/
│   └── quotes.json                      # 名言数据源（86.8KB，约 300 条）
│
├── generatorpb/                         # Protobuf 自动生成代码
│   ├── generator.pb.go                  # 消息类型
│   └── generatorpbconnect/
│       └── generator.connect.go         # ConnectRPC 客户端/服务端骨架
│
└── tools/                               # 辅助工具（独立运行，不参与编译）
    ├── fetch_quotes.go                  # 从网页抓取名言并保存到 data/quotes.json
    └── fetch_check.go                   # 检查网页结构，辅助调试爬虫
```

### 文件职责说明

| 文件 | 职责 |
|------|------|
| `main.go` | 核心业务逻辑 + HTTP 服务启动 + 控制台交互菜单 |
| `generator.proto` | 接口契约定义（RPC 方法 + 消息结构） |
| `generator_service_test.go` | 单元测试（参数校验、范围验证、稳定性测试） |
| `data/quotes.json` | 名言数据源，服务启动时自动加载 |
| `tools/fetch_quotes.go` | 名言爬虫，从指定网页抓取并解析名言 |
| `tools/fetch_check.go` | 网页结构探测工具，用于调试爬虫正则 |

## 快速启动

后端支持 **三种启动模式**，通过 `--mode` 参数控制：

```bash
# 模式一：仅启动 HTTP 服务（无控制台菜单）
go run . --mode=http

# 模式二：仅启动控制台交互模式（无 HTTP 服务）
go run . --mode=console

# 模式三：同时启动 HTTP 服务 + 控制台菜单（默认）
go run .
# 或显式指定
go run . --mode=all
```

### 命令行参数

| 参数 | 类型 | 默认值 | 说明 |
|------|------|--------|------|
| `--mode` | string | `all` | 启动模式：`http` / `console` / `all` |
| `--port` | int | `8080` | HTTP 服务监听端口 |

> 示例：`go run . --mode=http --port=9090` 仅启动 HTTP 服务并监听 9090 端口。

### 控制台交互模式

启动 `console` 或 `all` 模式后，控制台自动显示交互菜单：

```
========================================
  🎯 QuoteGenerator 控制台交互模式
========================================
  1. 生成随机数
  2. 获取随机名言
  0. 退出
========================================
```

| 数字 | 操作 | 说明 |
|------|------|------|
| `1` | 生成随机数 | 输入范围，返回随机整数 |
| `2` | 获取随机名言 | 从名言列表中随机返回一条 |
| `0` | 退出 | 退出交互菜单，HTTP 服务仍在后台运行 |

> 退出菜单后 HTTP 服务继续在后台运行，可通过 `http://localhost:8080` 正常调用 API。

## API 接口

### GetRandomNumber — 生成随机整数

```
POST /generator.GeneratorService/GetRandomNumber
Content-Type: application/json

请求: {"min": 1, "max": 100}
响应: {"number": 42}
```

**参数说明**：
- `min` (int32) — 随机数最小值（包含）
- `max` (int32) — 随机数最大值（包含）

**错误处理**：当 `min > max` 时返回 `CodeInvalidArgument` 错误。

### GetRandomQuote — 获取随机名言

```
POST /generator.GeneratorService/GetRandomQuote
Content-Type: application/json

请求: {}
响应: {"quote": "Stay hungry, stay foolish."}
```

**名言数据源**：服务启动时自动加载 `data/quotes.json` 文件（约 300 条名言），文件不存在时使用内置默认列表。

## 运行测试

```bash
go test -v
```

包含 3 个测试用例：
- `TestGetRandomNumber` — 参数校验 + 范围验证
- `TestGetRandomQuote` — 返回值有效性验证
- `TestGetRandomNumberMultipleTimes` — 100 次连续调用稳定性验证

## 辅助工具

### 抓取名言

```bash
go run tools/fetch_quotes.go
```

从指定网页抓取名言列表，自动解析并保存到 `data/quotes.json`。

### 检查网页结构

```bash
go run tools/fetch_check.go
```

探测目标网页的 HTML 结构，辅助调试爬虫正则表达式。

## 开发指南

### 修改接口定义

1. 编辑 `generator.proto`
2. 运行 protoc 重新生成代码：
   ```bash
   protoc --go_out=. --go_opt=paths=source_relative \
       --go_opt=Mgenerator.proto=quote_generator/generatorpb \
       --connect-go_out=. --connect-go_opt=paths=source_relative \
       generator.proto
   ```
3. 在 `main.go` 中实现新增的 RPC 方法

### 添加名言

编辑 `data/quotes.json` 文件，或运行 `go run tools/fetch_quotes.go` 从网页重新抓取。

## CORS 配置

当前允许 `http://localhost:3000` 跨域访问。生产环境部署时请修改 `main.go` 中的 `AllowedOrigins`。