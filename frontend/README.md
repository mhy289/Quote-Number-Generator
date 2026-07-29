# QuoteGenerator 前端应用

基于 **Next.js (App Router)** 构建的随机名言/数字生成器前端应用。

## 技术栈

| 技术 | 说明 |
|---|---|
| Next.js | React 全栈框架 (App Router) |
| TypeScript | 类型安全 |
| ConnectRPC-Web | 与后端 RPC 通信 |

## 快速启动

```bash
# 安装依赖
npm install

# 启动开发服务器
npm run dev
# 或
yarn dev
# 或
pnpm dev
# 或
bun dev
```

应用运行在 [http://localhost:3000](http://localhost:3000)。

## 功能

- **随机数字生成** — 指定范围生成随机整数
- **随机名言展示** — 从后端预设名言列表中随机获取一条

## 后端依赖

前端需要后端服务同时运行（默认连接 `http://localhost:8080`），请先启动后端：

```bash
cd ../backend
go run main.go
```

## 项目结构

```
frontend/
├── app/               # Next.js App Router 页面
│   ├── page.tsx       # 主页面
│   ├── layout.tsx     # 布局组件
│   └── globals.css    # 全局样式
├── public/            # 静态资源
├── package.json       # 依赖配置
└── README.md          # 本文件
```