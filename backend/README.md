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
├── generatorpb/
│   ├── generator.pb.go                  # 自动生成的 Protobuf 消息类型
│   └── generatorpbconnect/
│       └── generator.connect.go         # 自动生成的 ConnectRPC 客户端/服务端骨架
├── generator_service_test.go            # 单元测试
├── go.mod                               # Go 模块定义
├── go.sum                               # 依赖校验文件
└── README.md                            # 本文件
```

## 快速启动

### 启动方式

后端提供 **两种启动模式**，同时运行互不干扰：

#### 模式一：标准 HTTP 服务

```bash
# 启动服务（默认监听 :8080）
go run main.go
```

#### 模式二：控制台交互模式（内置）

启动后控制台自动显示交互菜单，无需前端即可直接测试功能：

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

**预设名言列表**：
1. "Stay hungry, stay foolish."
2. "Life is what happens when you're busy making other plans."
3. "The only limit to our realization of tomorrow is our doubts of today."
4. "Do not take life too seriously. You will never get out of it alive."

## 运行测试

```bash
go test -v
```

包含 3 个测试用例：
- `TestGetRandomNumber` — 参数校验 + 范围验证
- `TestGetRandomQuote` — 返回值有效性验证
- `TestGetRandomNumberMultipleTimes` — 100 次连续调用稳定性验证

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

编辑 `main.go` 中 `generatorService` 初始化时的 `quotes` 切片即可。

## CORS 配置

当前允许 `http://localhost:3000` 跨域访问。生产环境部署时请修改 `main.go` 中的 `AllowedOrigins`。